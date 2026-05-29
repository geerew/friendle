package dao

import (
	"context"

	"github.com/geerew/friendle/models"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListLeaderboardEntries returns scored members for a group
func (dao *DAO) ListLeaderboardEntries(ctx context.Context, groupID string) ([]*models.LeaderboardEntry, error) {
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

	var out []*models.LeaderboardEntry
	for rows.Next() {
		entry := &models.LeaderboardEntry{}
		if err := rows.Scan(&entry.UserID, &entry.DisplayName, &entry.TotalScore, &entry.RoundsPlayed); err != nil {
			return nil, err
		}
		out = append(out, entry)
	}

	return out, rows.Err()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Leaderboard returns leaderboard entries as maps for backward compatibility
func (dao *DAO) Leaderboard(ctx context.Context, groupID string) ([]map[string]interface{}, error) {
	entries, err := dao.ListLeaderboardEntries(ctx, groupID)
	if err != nil {
		return nil, err
	}

	out := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		out = append(out, map[string]interface{}{
			"userId": e.UserID, "displayName": e.DisplayName,
			"totalScore": e.TotalScore, "roundsPlayed": e.RoundsPlayed,
		})
	}

	return out, nil
}
