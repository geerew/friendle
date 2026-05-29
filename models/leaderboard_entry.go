package models

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// LeaderboardEntry is a computed score row for a group member
type LeaderboardEntry struct {
	UserID       string `db:"user_id"`
	DisplayName  string `db:"display_name"`
	TotalScore   int    `db:"total_score"`
	RoundsPlayed int    `db:"rounds_played"`
}
