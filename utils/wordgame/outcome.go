package wordgame

// GuessOutcomeLabel returns correct, partial, or incorrect for a graded row
func GuessOutcomeLabel(result []TileState) string {
	if IsWin(result) {
		return "correct"
	}

	for _, r := range result {
		if r == TileCorrect || r == TilePresent {
			return "partial"
		}
	}

	return "incorrect"
}
