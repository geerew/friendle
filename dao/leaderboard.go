package dao

import (
	"context"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// LeaderboardEntry is a computed score row for a group member
type LeaderboardEntry struct {
	UserID       string `db:"user_id"`
	DisplayName  string `db:"display_name"`
	TotalScore   int    `db:"total_score"`
	RoundsPlayed int    `db:"rounds_played"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListLeaderboardEntries returns scored members for a group
func (dao *DAO) ListLeaderboardEntries(ctx context.Context, groupID string) ([]*LeaderboardEntry, error) {
	query := `
SELECT u.id AS user_id, u.display_name, COALESCE(SUM(rp.score), 0) AS total_score,
       COUNT(rp.id) AS rounds_played
FROM group_members gm
JOIN users u ON u.id = gm.user_id
LEFT JOIN rounds r ON r.group_id = gm.group_id AND r.status = 'completed'
LEFT JOIN round_participations rp ON rp.round_id = r.id AND rp.user_id = gm.user_id
WHERE gm.group_id = ?
GROUP BY u.id, u.display_name
ORDER BY total_score DESC`

	rows, err := dao.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*LeaderboardEntry
	for rows.Next() {
		entry := &LeaderboardEntry{}
		if err := rows.Scan(&entry.UserID, &entry.DisplayName, &entry.TotalScore, &entry.RoundsPlayed); err != nil {
			return nil, err
		}
		out = append(out, entry)
	}

	return out, rows.Err()
}
