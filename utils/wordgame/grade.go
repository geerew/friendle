package wordgame

import "strings"

type TileState string

const (
	TileCorrect TileState = "correct"
	TilePresent TileState = "present"
	TileAbsent  TileState = "absent"
)

func Grade(guess, answer string) []TileState {
	guess = strings.ToUpper(guess)
	answer = strings.ToUpper(answer)
	result := make([]TileState, 5)
	remaining := make(map[byte]int)
	for i := 0; i < 5; i++ {
		if guess[i] == answer[i] {
			result[i] = TileCorrect
		} else {
			remaining[answer[i]]++
		}
	}
	for i := 0; i < 5; i++ {
		if result[i] == TileCorrect {
			continue
		}
		ch := guess[i]
		if remaining[ch] > 0 {
			result[i] = TilePresent
			remaining[ch]--
		} else {
			result[i] = TileAbsent
		}
	}
	return result
}

func IsWin(result []TileState) bool {
	for _, r := range result {
		if r != TileCorrect {
			return false
		}
	}
	return true
}
