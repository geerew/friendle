package dao

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_CreateRound(t *testing.T) {
	// Test successfully inserting a round record
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Round Group", user.ID)

		round := &models.Round{
			GroupID:      group.ID,
			RoundDate:    "2026-01-01",
			PickerUserID: user.ID,
			Status:       types.RoundAwaitingWord,
		}
		require.NoError(t, dao.CreateRound(ctx, round))
		require.NotEmpty(t, round.ID)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_GetRound(t *testing.T) {
	// Test successfully retrieving a round row
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Get Round Group", user.ID)

		round := &models.Round{
			GroupID:      group.ID,
			RoundDate:    "2026-01-02",
			PickerUserID: user.ID,
			Status:       types.RoundActive,
		}
		require.NoError(t, dao.CreateRound(ctx, round))

		record, err := dao.GetRound(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: round.ID}))
		require.NoError(t, err)
		require.Equal(t, round.ID, record.ID)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_UpdateRound(t *testing.T) {
	// Test successfully updating a round
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Update Round", user.ID)

		round := &models.Round{
			GroupID:      group.ID,
			RoundDate:    "2026-01-05",
			PickerUserID: user.ID,
			Status:       types.RoundAwaitingWord,
		}
		require.NoError(t, dao.CreateRound(ctx, round))

		word := "hello"
		round.Status = types.RoundActive
		round.WordPlain = &word
		require.NoError(t, dao.UpdateRound(ctx, round))

		record, err := dao.GetRound(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: round.ID}))
		require.NoError(t, err)
		require.Equal(t, types.RoundActive, record.Status)
		require.NotNil(t, record.WordPlain)
		require.Equal(t, "hello", *record.WordPlain)
	})

	// Test error due to missing ID
	t.Run("missing id", func(t *testing.T) {
		dao, ctx := setup(t)

		require.ErrorIs(t, dao.UpdateRound(ctx, &models.Round{}), utils.ErrId)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_DeleteRounds(t *testing.T) {
	// Test successfully deleting a round
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Delete Round", user.ID)

		round := &models.Round{
			GroupID:      group.ID,
			RoundDate:    "2026-01-06",
			PickerUserID: user.ID,
			Status:       types.RoundAwaitingWord,
		}
		require.NoError(t, dao.CreateRound(ctx, round))

		require.NoError(t, dao.DeleteRounds(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: round.ID})))

		record, err := dao.GetRound(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: round.ID}))
		require.NoError(t, err)
		require.Nil(t, record)
	})

	// Test error due to missing where clause
	t.Run("missing where", func(t *testing.T) {
		dao, ctx := setup(t)

		require.ErrorIs(t, dao.DeleteRounds(ctx, nil), utils.ErrWhere)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_UpdateRoundEntry(t *testing.T) {
	// Test successfully updating a round entry
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Entry Group", user.ID)

		round := &models.Round{
			GroupID:      group.ID,
			RoundDate:    "2026-01-07",
			PickerUserID: user.ID,
			Status:       types.RoundActive,
		}
		require.NoError(t, dao.CreateRound(ctx, round))

		entry := &models.RoundEntry{
			RoundID: round.ID,
			UserID:  user.ID,
		}
		require.NoError(t, dao.CreateRoundEntry(ctx, entry))

		entry.Score = 5
		entry.Solved = true
		require.NoError(t, dao.UpdateRoundEntry(ctx, entry))

		record, err := dao.GetRoundEntry(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: entry.ID}))
		require.NoError(t, err)
		require.Equal(t, 5, record.Score)
		require.True(t, record.Solved)
	})

	// Test error due to missing ID
	t.Run("missing id", func(t *testing.T) {
		dao, ctx := setup(t)

		require.ErrorIs(t, dao.UpdateRoundEntry(ctx, &models.RoundEntry{}), utils.ErrId)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_CreateGuess(t *testing.T) {
	// Test successfully inserting and listing guesses in attempt order
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Guess Group", user.ID)

		round := &models.Round{
			GroupID:      group.ID,
			RoundDate:    "2026-01-08",
			PickerUserID: user.ID,
			Status:       types.RoundActive,
		}
		require.NoError(t, dao.CreateRound(ctx, round))

		require.NoError(t, dao.CreateGuess(ctx, &models.Guess{
			RoundID: round.ID,
			UserID:  user.ID,
			Attempt: 2,
			Word:    "world",
			Result:  types.TileStates{types.TileAbsent, types.TileAbsent, types.TileAbsent, types.TileAbsent, types.TileAbsent},
			Outcome: types.GuessOutcomeIncorrect,
		}))
		require.NoError(t, dao.CreateGuess(ctx, &models.Guess{
			RoundID: round.ID,
			UserID:  user.ID,
			Attempt: 1,
			Word:    "crane",
			Result:  types.TileStates{types.TileCorrect, types.TileAbsent, types.TileAbsent, types.TileAbsent, types.TileAbsent},
			Outcome: types.GuessOutcomePartial,
		}))

		guesses, err := dao.ListGuesses(ctx, NewOptions().WithWhere(squirrel.Eq{models.GUESS_ROUND_ID: round.ID}))
		require.NoError(t, err)
		require.Len(t, guesses, 2)
		require.Equal(t, 1, guesses[0].Attempt)
		require.Equal(t, 2, guesses[1].Attempt)

		count, err := dao.CountGuesses(ctx, NewOptions().WithWhere(squirrel.Eq{models.GUESS_ROUND_ID: round.ID}))
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_DeleteGuesses(t *testing.T) {
	// Test error due to missing where clause
	t.Run("missing where", func(t *testing.T) {
		dao, ctx := setup(t)

		require.ErrorIs(t, dao.DeleteGuesses(ctx, nil), utils.ErrWhere)
	})
}
