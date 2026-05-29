package api

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
	"github.com/geerew/friendle/utils/wordgame"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initRoundRoutes initializes round routes under groups
func (r *Router) initRoundRoutes() {
	g := r.apiGroup("groups")

	g.Get("/:id/rounds", r.requireAccess(accessGroupMemberScope), r.listGroupRounds)
	g.Get("/:id/rounds/current", r.requireAccess(accessGroupMemberScope), r.getGroupRound)
	g.Post("/:id/rounds/current/word", r.requireAccess(accessGroupMemberScope), r.createGroupRoundWord)
	g.Post("/:id/rounds/current/guesses", r.requireAccess(accessGroupMemberScope), r.createGroupRoundGuess)
	g.Get("/:id/rounds/current/reveal", r.requireAccess(accessGroupMemberScope), r.getGroupRoundReveal)
	g.Get("/:id/rounds/:roundId", r.requireAccess(accessGroupMemberScope), r.getGroupRoundByID)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listGroupRounds returns paginated round history for a group
func (r *Router) listGroupRounds(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)
	groupID := c.Params("id")

	dbOpts := dao.NewOptions().WithPagination(paginationFromCtx(c))
	rounds, err := r.appDao.ListRoundsForGroup(ctx, groupID, dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "List failed", err)
	}

	items := make([]*roundSummaryResponse, 0, len(rounds))
	for _, round := range rounds {
		items = append(items, roundSummaryResponseHelper(round, false))
	}

	pResult, err := dbOpts.Pagination.BuildResult(items)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupRoundByID returns a round with role-filtered player data
func (r *Router) getGroupRoundByID(c *fiber.Ctx) error {
	p, ctx := principalAndCtx(c)
	groupID := c.Params("id")
	roundID := c.Params("roundId")

	round, err := r.appDao.GetRound(ctx, roundID)
	if err != nil || round == nil || round.GroupID != groupID {
		return errorResponse(c, fiber.StatusNotFound, "Round not found", nil)
	}

	gr, err := r.appDao.GetGroupRound(ctx, roundID, dao.RoundDetailLoad{
		Participations: true,
		Guesses:        true,
		Picker:         dao.RoundIsRevealed(round.Status),
	})
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Load failed", err)
	}

	return c.JSON(groupRoundDetailResponseHelper(gr, p.UserID, false))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupRound returns the current round state for the caller
func (r *Router) getGroupRound(c *fiber.Ctx) error {
	p, ctx := principalAndCtx(c)
	groupID := c.Params("id")
	roundDate := time.Now().Format("2006-01-02")

	round, err := r.appDao.GetCurrentRound(ctx, groupID, roundDate)
	if err != nil || round == nil {
		return c.JSON(groupRoundResponseHelper(nil, nil, nil, p.UserID))
	}

	participation, _ := r.appDao.GetRoundParticipation(ctx, round.ID, p.UserID)
	var guesses []*models.Guess
	if participation != nil {
		guesses, _ = r.appDao.ListGuessesForRoundUser(ctx, round.ID, p.UserID)
	}

	return c.JSON(groupRoundResponseHelper(round, participation, guesses, p.UserID))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroupRoundWord submits the picker's word for the current round
func (r *Router) createGroupRoundWord(c *fiber.Ctx) error {
	p, ctx := principalAndCtx(c)
	groupID := c.Params("id")
	roundDate := time.Now().Format("2006-01-02")

	round, _ := r.appDao.GetCurrentRound(ctx, groupID, roundDate)
	if round == nil || round.Status != types.RoundAwaitingWord {
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
	round.Status = types.RoundActive
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
	if round == nil || round.Status != types.RoundActive {
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

	participation, _ := r.appDao.GetRoundParticipation(ctx, round.ID, p.UserID)
	if participation == nil {
		now := time.Now().UTC().Format(time.RFC3339)
		participation = &models.RoundParticipation{
			RoundID: round.ID, UserID: p.UserID, FirstGuessAt: &now,
		}
		_ = r.appDao.CreateRoundParticipation(ctx, participation)
	}

	if participation.Finished {
		return errorResponse(c, fiber.StatusBadRequest, "Already finished", nil)
	}

	attempts, err := r.appDao.CountGuessesForRoundUser(ctx, round.ID, p.UserID)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Load failed", err)
	}
	if attempts >= 6 {
		return errorResponse(c, fiber.StatusBadRequest, "No attempts left", nil)
	}

	answer := *round.WordPlain
	result := wordgame.Grade(word, answer)
	won := wordgame.IsWin(result)
	resultJSON, _ := json.Marshal(result)

	guess := &models.Guess{
		RoundID:       round.ID,
		UserID:        p.UserID,
		AttemptNumber: attempts + 1,
		Word:          word,
		Result:        string(resultJSON),
	}
	if err := r.appDao.CreateGuess(ctx, guess); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Guess failed", err)
	}

	participation.Solved = won
	if won || guess.AttemptNumber >= 6 {
		participation.Finished = true
		participation.Score = wordgame.ScoreForAttempt(guess.AttemptNumber, won)
		now := time.Now().UTC().Format(time.RFC3339)
		participation.CompletedAt = &now
	}
	_ = r.appDao.UpdateRoundParticipation(ctx, participation)

	if r.allGuessersDone(ctx, round) {
		round.Status = types.RoundCompleted
		_ = r.appDao.UpdateRound(ctx, round)
	}

	return c.JSON(&groupRoundGuessResponse{
		Result:   result,
		Attempt:  guess.AttemptNumber,
		Won:      won,
		Finished: participation.Finished,
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

	if !dao.RoundIsRevealed(round.Status) {
		return errorResponse(c, fiber.StatusBadRequest, "Round not finished", nil)
	}

	gr, err := r.appDao.GetGroupRound(ctx, round.ID, dao.RoundDetailLoad{
		Participations: true,
		Guesses:        false,
		Picker:         true,
	})
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Load failed", err)
	}

	return c.JSON(groupRoundRevealResponseHelper(gr))
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
	if dao.RoundIsRevealed(round.Status) {
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
	participations, _ := r.appDao.ListRoundParticipationsForRound(ctx, round.ID)

	doneByUser := make(map[string]bool, len(participations))
	for _, p := range participations {
		if p.Finished {
			doneByUser[p.UserID] = true
		}
	}

	done, need := 0, 0
	for _, m := range members {
		if m.UserID == round.PickerUserID {
			continue
		}

		need++
		if doneByUser[m.UserID] {
			done++
		}
	}

	return need > 0 && done >= need
}
