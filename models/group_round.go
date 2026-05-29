package models

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupRound is a round with all player participations for a group
type GroupRound struct {
	Round   *Round
	Players []*UserRound
}
