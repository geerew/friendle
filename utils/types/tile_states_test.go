package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTileStates(t *testing.T) {
	// Test successfully marshaling and scanning stored JSON
	t.Run("scan value", func(t *testing.T) {
		original := TileStates{TileAbsent, TilePresent, TileCorrect, TileAbsent, TileAbsent}
		val, err := original.Value()
		require.NoError(t, err)

		var scanned TileStates
		require.NoError(t, scanned.Scan(val))
		require.Equal(t, original, scanned)
	})

	// Test successfully unmarshaling API JSON
	t.Run("json", func(t *testing.T) {
		var states TileStates
		require.NoError(t, states.UnmarshalJSON([]byte(`["absent","present","correct","absent","absent"]`)))
		require.Equal(t, TileStates{TileAbsent, TilePresent, TileCorrect, TileAbsent, TileAbsent}, states)
	})

	// Test successfully detecting a winning row
	t.Run("is win", func(t *testing.T) {
		require.True(t, TileStates{TileCorrect, TileCorrect, TileCorrect, TileCorrect, TileCorrect}.IsWin())
		require.False(t, TileStates{TileCorrect, TilePresent, TileCorrect, TileCorrect, TileCorrect}.IsWin())
	})
}
