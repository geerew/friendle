package words

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ScoreForAttempt returns points for a solved guess on the given attempt
func ScoreForAttempt(attempt int, solved bool) int {
	if !solved {
		return 0
	}

	switch attempt {
	case 1:
		return 100
	case 2:
		return 80
	case 3:
		return 60
	case 4:
		return 45
	case 5:
		return 30
	case 6:
		return 20
	default:
		return 0
	}
}
