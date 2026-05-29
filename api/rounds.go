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
	g.Get("/:id/rounds/current", r.requireAccess(accessGroupMemberScope), r.getGroupRound)
	g.Post("/:id/rounds/current/word", r.requireAccess(accessGroupMemberScope), r.createGroupRoundWord)
	g.Post("/:id/rounds/current/guesses", r.requireAccess(accessGroupMemberScope), r.createGroupRoundGuess)
	g.Get("/:id/rounds/current/reveal", r.requireAccess(accessGroupMemberScope), r.getGroupRoundReveal)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupRound returns the current round state for the caller
func (r *Router) getGroupRound(c *fiber.Ctx) error {
	p, ctx := principalAndCtx(c)
	groupID := c.Params("id")
	roundDate := time.Now().Format("2006-01-02")

	round, err := r.appDao.GetCurrentRound(ctx, groupID, roundDate)
	if err != nil || round == nil {
		return c.JSON(groupRoundResponseHelper(nil, nil, p.UserID))
	}

	guess, _ := r.appDao.GetGuess(ctx, round.ID, p.UserID)

	return c.JSON(groupRoundResponseHelper(round, guess, p.UserID))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroupRoundWord submits the picker's word for the current round
func (r *Router) createGroupRoundWord(c *fiber.Ctx) error {
	p, ctx := principalAndCtx(c)
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
	p, ctx := principalAndCtx(c)
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
	var rows []guessRowResponse
	_ = json.Unmarshal([]byte(guess.RowsJSON), &rows)
	rows = append(rows, guessRowResponse{Word: word, Result: result})
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

	return c.JSON(&groupRoundGuessResponse{
		Result:   result,
		Attempt:  guess.AttemptsUsed,
		Won:      won,
		Finished: guess.Finished,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupRoundReveal returns the completed round word and guess summaries
func (r *Router) getGroupRoundReveal(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)
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
	resp := &groupRoundRevealResponse{
		Status:       round.Status,
		PickerUserID: round.PickerUserID,
		Guesses:      []groupRoundRevealGuessResponse{},
	}
	if picker != nil {
		resp.PickerDisplayName = picker.DisplayName
	}

	if round.Status == models.RoundCompleted && round.WordPlain != nil {
		resp.Word = *round.WordPlain
	}

	guesses, _ := r.appDao.ListGuessesForRound(ctx, round.ID)
	for _, g := range guesses {
		if g.UserID == round.PickerUserID {
			continue
		}

		u, _ := r.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: g.UserID}))
		name := g.UserID
		if u != nil {
			name = u.DisplayName
		}

		resp.Guesses = append(resp.Guesses, groupRoundRevealGuessResponse{
			UserID: g.UserID, DisplayName: name,
			AttemptsUsed: g.AttemptsUsed, Solved: g.Solved, Score: g.Score,
		})
	}

	return c.JSON(resp)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundSummary returns a summary of the current round for a group member
func (r *Router) roundSummary(ctx context.Context, g *models.Group, userID string) *groupRoundSummaryResponse {
	roundDate := time.Now().Format("2006-01-02")
	round, _ := r.appDao.GetCurrentRound(ctx, g.ID, roundDate)
	if round == nil {
		return &groupRoundSummaryResponse{Status: "none"}
	}

	m, _ := r.appDao.GetGroupMember(ctx, g.ID, userID)
	yourRole := "guesser"
	if round.PickerUserID == userID {
		yourRole = "picker"
	}

	out := &groupRoundSummaryResponse{
		Status:   string(round.Status),
		YourRole: yourRole,
	}
	if round.Status == models.RoundCompleted || round.Status == models.RoundSkipped {
		out.CanReveal = true
	}

	if m != nil {
		out.GroupRole = m.GroupRole
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
