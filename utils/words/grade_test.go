package words

import (
	"testing"

	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

func TestGrade(t *testing.T) {
	// Test successfully grading an exact match
	t.Run("win", func(t *testing.T) {
		result := Grade("ABOUT", "ABOUT")
		require.True(t, result.IsWin())
		require.Equal(t, types.TileStates{
			types.TileCorrect, types.TileCorrect, types.TileCorrect, types.TileCorrect, types.TileCorrect,
		}, result)
	})

	// Test successfully grading letters that match in the wrong position
	t.Run("partial", func(t *testing.T) {
		result := Grade("AROSE", "ABOUT")
		require.Equal(t, types.TileStates{
			types.TileCorrect, types.TileAbsent, types.TileCorrect, types.TileAbsent, types.TileAbsent,
		}, result)
	})
}

func TestScoreForAttempt(t *testing.T) {
	// Test successfully scoring solved attempts
	t.Run("solved", func(t *testing.T) {
		require.Equal(t, 100, ScoreForAttempt(1, true))
		require.Equal(t, 20, ScoreForAttempt(6, true))
	})

	// Test no score when the round was not solved
	t.Run("unsolved", func(t *testing.T) {
		require.Equal(t, 0, ScoreForAttempt(3, false))
	})
}
