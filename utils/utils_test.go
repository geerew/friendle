package utils

import (
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
