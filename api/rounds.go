package api

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/wordgame"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initRoundRoutes initializes round routes under groups
func (r *Router) initRoundRoutes() {
	g := r.apiGroup("groups")

	// Rounds
	g.Get("/:id/rounds/current", r.requireAuth, r.requireGroupMember, r.getGroupRound)
	g.Post("/:id/rounds/current/word", r.requireAuth, r.requireGroupMember, r.createGroupRoundWord)
	g.Post("/:id/rounds/current/guesses", r.requireAuth, r.requireGroupMember, r.createGroupRoundGuess)
	g.Get("/:id/rounds/current/reveal", r.requireAuth, r.requireGroupMember, r.getGroupRoundReveal)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupRound returns the current round state for the caller
func (r *Router) getGroupRound(c *fiber.Ctx) error {
	p, ctx, _ := principalCtx(c)
	groupID := c.Params("id")
	roundDate := time.Now().Format("2006-01-02")

	round, err := r.appDao.GetCurrentRound(ctx, groupID, roundDate)
	if err != nil || round == nil {
		return c.JSON(fiber.Map{"status": "none"})
	}

	guess, _ := r.appDao.GetGuess(ctx, round.ID, p.UserID)
	yourRole := "guesser"
	if round.PickerUserID == p.UserID {
		yourRole = "picker"
	}

	resp := fiber.Map{
		"roundId": round.ID, "status": round.Status, "yourRole": yourRole,
		"attemptsUsed": 0, "finished": false, "rows": []interface{}{},
	}
	if guess != nil {
		var rows []interface{}
		_ = json.Unmarshal([]byte(guess.RowsJSON), &rows)
		resp["attemptsUsed"] = guess.AttemptsUsed
		resp["finished"] = guess.Finished
		resp["rows"] = rows
		resp["solved"] = guess.Solved
	}

	return c.JSON(resp)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroupRoundWord submits the picker's word for the current round
func (r *Router) createGroupRoundWord(c *fiber.Ctx) error {
	p, ctx, _ := principalCtx(c)
	groupID := c.Params("id")
	roundDate := time.Now().Format("2006-01-02")

	round, _ := r.appDao.GetCurrentRound(ctx, groupID, roundDate)
	if round == nil || round.Status != models.RoundAwaitingWord {
		return errorResponse(c, fiber.StatusBadRequest, "No round awaiting word", nil)
	}

	if round.PickerUserID != p.UserID {
		return errorResponse(c, fiber.StatusForbidden, "Not the picker", nil)
	}

	req := &submitRoundWordRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid body", nil)
	}

	word := strings.ToUpper(strings.TrimSpace(req.Word))
	if len(word) != 5 {
		return errorResponse(c, fiber.StatusBadRequest, "Word must be 5 letters", nil)
	}

	if !r.app.Dictionary.IsValidAnswer(word) {
		return errorResponse(c, fiber.StatusBadRequest, "Not a valid answer word", nil)
	}

	hash := wordgame.HashWord(r.app.Config.DataDir, word)
	round.WordHash = &hash
	round.WordPlain = &word
	round.Status = models.RoundActive
	if err := r.appDao.UpdateRound(ctx, round); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to save word", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroupRoundGuess submits a guess for the current round
func (r *Router) createGroupRoundGuess(c *fiber.Ctx) error {
	p, ctx, _ := principalCtx(c)
	groupID := c.Params("id")
	roundDate := time.Now().Format("2006-01-02")

	round, _ := r.appDao.GetCurrentRound(ctx, groupID, roundDate)
	if round == nil || round.Status != models.RoundActive {
		return errorResponse(c, fiber.StatusBadRequest, "Round not active", nil)
	}

	if round.PickerUserID == p.UserID {
		return errorResponse(c, fiber.StatusForbidden, "Picker cannot guess", nil)
	}

	if round.WordPlain == nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Word not set", nil)
	}

	req := &submitRoundGuessRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid body", nil)
	}

	word := strings.ToUpper(strings.TrimSpace(req.Word))
	if len(word) != 5 {
		return errorResponse(c, fiber.StatusBadRequest, "Word must be 5 letters", nil)
	}

	if !r.app.Dictionary.IsValidGuess(word) {
		return errorResponse(c, fiber.StatusBadRequest, "Not in word list", nil)
	}

	guess, _ := r.appDao.GetGuess(ctx, round.ID, p.UserID)
	if guess == nil {
		guess = &models.Guess{RoundID: round.ID, UserID: p.UserID, RowsJSON: "[]"}
		now := time.Now().UTC().Format(time.RFC3339)
		guess.FirstGuessAt = &now
		_ = r.appDao.CreateGuess(ctx, guess)
	}

	if guess.Finished {
		return errorResponse(c, fiber.StatusBadRequest, "Already finished", nil)
	}

	if guess.AttemptsUsed >= 6 {
		return errorResponse(c, fiber.StatusBadRequest, "No attempts left", nil)
	}

	answer := *round.WordPlain
	result := wordgame.Grade(word, answer)
	won := wordgame.IsWin(result)
	var rows []map[string]interface{}
	_ = json.Unmarshal([]byte(guess.RowsJSON), &rows)
	rows = append(rows, map[string]interface{}{
		"word": word, "result": result,
	})
	guess.AttemptsUsed++
	guess.Solved = won
	b, _ := json.Marshal(rows)
	guess.RowsJSON = string(b)
	if won || guess.AttemptsUsed >= 6 {
		guess.Finished = true
		guess.Score = wordgame.ScoreForAttempt(guess.AttemptsUsed, won)
		now := time.Now().UTC().Format(time.RFC3339)
		guess.CompletedAt = &now
	}
	_ = r.appDao.UpdateGuess(ctx, guess)

	if r.allGuessersDone(ctx, round) {
		round.Status = models.RoundCompleted
		_ = r.appDao.UpdateRound(ctx, round)
	}

	return c.JSON(fiber.Map{
		"result": result, "attempt": guess.AttemptsUsed, "won": won, "finished": guess.Finished,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupRoundReveal returns the completed round word and guess summaries
func (r *Router) getGroupRoundReveal(c *fiber.Ctx) error {
	_, ctx, _ := principalCtx(c)
	groupID := c.Params("id")
	roundDate := time.Now().Format("2006-01-02")

	round, _ := r.appDao.GetCurrentRound(ctx, groupID, roundDate)
	if round == nil {
		return errorResponse(c, fiber.StatusNotFound, "No round", nil)
	}

	if round.Status != models.RoundCompleted && round.Status != models.RoundSkipped {
		return errorResponse(c, fiber.StatusBadRequest, "Round not finished", nil)
	}

	picker, _ := r.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: round.PickerUserID}))
	resp := fiber.Map{
		"status":       round.Status,
		"pickerUserId": round.PickerUserID,
	}
	if picker != nil {
		resp["pickerDisplayName"] = picker.DisplayName
	}

	if round.Status == models.RoundCompleted && round.WordPlain != nil {
		resp["word"] = *round.WordPlain
	}

	guesses, _ := r.appDao.ListGuessesForRound(ctx, round.ID)
	var summaries []fiber.Map
	for _, g := range guesses {
		if g.UserID == round.PickerUserID {
			continue
		}

		u, _ := r.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: g.UserID}))
		name := g.UserID
		if u != nil {
			name = u.DisplayName
		}

		summaries = append(summaries, fiber.Map{
			"userId": g.UserID, "displayName": name,
			"attemptsUsed": g.AttemptsUsed, "solved": g.Solved, "score": g.Score,
		})
	}
	resp["guesses"] = summaries

	return c.JSON(resp)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundSummary returns a summary of the current round for a group member
func (r *Router) roundSummary(ctx context.Context, g *models.Group, userID string) fiber.Map {
	roundDate := time.Now().Format("2006-01-02")
	round, _ := r.appDao.GetCurrentRound(ctx, g.ID, roundDate)
	if round == nil {
		return fiber.Map{"status": "none"}
	}

	m, _ := r.appDao.GetGroupMember(ctx, g.ID, userID)
	yourRole := "guesser"
	if round.PickerUserID == userID {
		yourRole = "picker"
	}

	out := fiber.Map{
		"status": round.Status, "yourRole": yourRole,
	}
	if round.Status == models.RoundCompleted || round.Status == models.RoundSkipped {
		out["canReveal"] = true
	}
	if m != nil {
		out["groupRole"] = m.GroupRole
	}

	return out
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// allGuessersDone reports whether every non-picker member has finished guessing
func (r *Router) allGuessersDone(ctx context.Context, round *models.Round) bool {
	members, _ := r.appDao.ListGroupMembers(ctx, round.GroupID)
	guesses, _ := r.appDao.ListGuessesForRound(ctx, round.ID)
	done, need := 0, 0
	for _, m := range members {
		if m.UserID == round.PickerUserID {
			continue
		}

		need++
		for _, g := range guesses {
			if g.UserID == m.UserID && g.Finished {
				done++
			}
		}
	}

	return need > 0 && done >= need
}
