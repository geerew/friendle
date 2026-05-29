package models

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UserRound is one user's play in a round with guess attempts loaded
type UserRound struct {
	Participation *RoundParticipation
	Guesses       []*Guess
	User          *User
}
