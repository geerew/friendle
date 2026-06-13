package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterMap(t *testing.T) {
	// Test successfully keeping mapped values when fn returns ok true
	t.Run("success", func(t *testing.T) {
		out := FilterMap([]int{1, 2, 3, 4}, func(n int) (int, bool) {
			if n%2 == 0 {
				return n * 10, true
			}

			return 0, false
		})

		require.Equal(t, []int{20, 40}, out)
	})

	// Test successfully returning an empty slice when no items pass
	t.Run("empty", func(t *testing.T) {
		out := FilterMap([]int{1, 3, 5}, func(n int) (int, bool) {
			return n, n%2 == 0
		})

		require.Empty(t, out)
	})
}

// Test successfully normalizing and validating a group name
func TestNormalizeGroupName(t *testing.T) {
	const maxLength = 64

	// Test successfully trimming and accepting a valid name
	t.Run("success", func(t *testing.T) {
		name, err := NormalizeGroupName("  Friends  ", maxLength)
		require.NoError(t, err)
		require.Equal(t, "Friends", name)
	})

	// Test error due to an empty name
	t.Run("empty", func(t *testing.T) {
		_, err := NormalizeGroupName("   ", maxLength)
		require.ErrorIs(t, err, ErrGroupName)
	})

	// Test error due to a name that is too long
	t.Run("too long", func(t *testing.T) {
		_, err := NormalizeGroupName(strings.Repeat("a", maxLength+1), maxLength)
		require.ErrorIs(t, err, ErrGroupNameTooLong)
	})
}
