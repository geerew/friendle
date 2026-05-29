package wordgame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test successfully round-tripping tile states through JSON and SQL value helpers
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
}
