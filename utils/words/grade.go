package words

import (
	"strings"

	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Grade compares a guess to the answer and returns per-letter tile feedback
func Grade(guess, answer string) types.TileStates {
	guess = strings.ToUpper(guess)
	answer = strings.ToUpper(answer)
	result := make(types.TileStates, 5)
	remaining := make(map[byte]int)

	for i := range 5 {
		if guess[i] == answer[i] {
			result[i] = types.TileCorrect
		} else {
			remaining[answer[i]]++
		}
	}

	for i := range 5 {
		if result[i] == types.TileCorrect {
			continue
		}

		ch := guess[i]
		if remaining[ch] > 0 {
			result[i] = types.TilePresent
			remaining[ch]--
		} else {
			result[i] = types.TileAbsent
		}
	}

	return result
}
