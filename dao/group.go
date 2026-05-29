package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ensure types import used

var defaultAdminGroupsListOrderBy = []string{models.GROUP_TABLE + "." + models.BASE_CREATED_AT + " desc"}

func (dao *DAO) CreateGroup(ctx context.Context, g *models.Group) error {
	if g.ID == "" {
		g.RefreshId()
	}
	g.RefreshCreatedAt()
	g.RefreshUpdatedAt()
	return createGeneric(ctx, dao, *newBuilderOptions(models.GROUP_TABLE).WithData(map[string]interface{}{
		models.BASE_ID: g.ID, "name": g.Name, "created_by": g.CreatedBy,
		"interval_hours": g.IntervalHours, "timezone": g.Timezone,
		models.BASE_CREATED_AT: g.CreatedAt, models.BASE_UPDATED_AT: g.UpdatedAt,
	}))
}

func (dao *DAO) GetGroup(ctx context.Context, id string) (*models.Group, error) {
	return getGeneric[models.Group](ctx, dao, *newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})).
		WithLimit(1))
}

func (dao *DAO) ListGroups(ctx context.Context, dbOpts *Options) ([]*models.Group, error) {
	return listGeneric[models.Group](ctx, dao, *newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).SetDbOpts(dbOpts))
}

// ListAdminGroups returns groups with member counts for the site admin list
func (dao *DAO) ListAdminGroups(ctx context.Context, dbOpts *Options) ([]*AdminGroupRow, error) {
	g := models.GROUP_TABLE
	gm := models.GROUP_MEMBER_TABLE

	applyDefaultOrderBy(dbOpts, defaultAdminGroupsListOrderBy)

	return listGeneric[AdminGroupRow](ctx, dao, *newBuilderOptions(g).
		WithColumns(
			g+"."+models.BASE_ID+" AS id",
			g+".name AS name",
			"COUNT("+gm+".id) AS member_count",
		).
		WithLeftJoin(gm, gm+".group_id = "+g+"."+models.BASE_ID).
		WithGroupBy(g+"."+models.BASE_ID, g+".name").
		SetDbOpts(dbOpts))
}

func (dao *DAO) SearchGroupSummaries(ctx context.Context, q string, dbOpts *Options) ([]*GroupSearchRow, error) {
	g := models.GROUP_TABLE
	gm := models.GROUP_MEMBER_TABLE
	like := "%" + q + "%"

	if dbOpts == nil {
		dbOpts = NewOptions()
	}

	searchOpts := NewOptions().
		WithWhere(squirrel.Like{"LOWER(" + g + ".name)": like}).
		WithOrderByClause(searchGroupRelevanceOrder(g, q))

	if dbOpts.Pagination != nil {
		searchOpts = searchOpts.WithPagination(dbOpts.Pagination)
	}

	return listGeneric[GroupSearchRow](ctx, dao, *newBuilderOptions(g).
		WithColumns(
			g+"."+models.BASE_ID+" AS id",
			g+".name AS name",
			"(SELECT COUNT(*) FROM "+gm+" gm_count WHERE gm_count.group_id = "+g+"."+models.BASE_ID+") AS member_count",
		).
		SetDbOpts(searchOpts))
}

// searchGroupRelevanceOrder ranks exact name matches first, then prefix matches, then other
// substring matches, then shorter names, then name ascending, then newest created
func searchGroupRelevanceOrder(g, q string) squirrel.Sqlizer {
	return squirrel.Expr(
		`CASE WHEN LOWER(`+g+`.name) = ? THEN 0 WHEN LOWER(`+g+`.name) LIKE ? THEN 1 ELSE 2 END, LENGTH(`+g+`.name), LOWER(`+g+`.name), `+g+`.`+models.BASE_CREATED_AT+` DESC`,
		q, q+"%",
	)
}

func (dao *DAO) SearchGroups(ctx context.Context, q string) ([]*models.Group, error) {
	like := "%" + q + "%"
	return listGeneric[models.Group](ctx, dao, *newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Like{"LOWER(name)": like})))
}

func (dao *DAO) UpdateGroup(ctx context.Context, g *models.Group) error {
	g.RefreshUpdatedAt()
	_, err := updateGeneric(ctx, dao, *newBuilderOptions(models.GROUP_TABLE).WithData(map[string]interface{}{
		"name": g.Name, "interval_hours": g.IntervalHours, "timezone": g.Timezone,
		models.BASE_UPDATED_AT: g.UpdatedAt,
	}).SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: g.ID})))
	return err
}

func (dao *DAO) CreateGroupMember(ctx context.Context, m *models.GroupMember) error {
	if m.ID == "" {
		m.RefreshId()
	}
	m.RefreshCreatedAt()
	m.RefreshUpdatedAt()
	return createGeneric(ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).WithData(map[string]interface{}{
		models.BASE_ID: m.ID, "group_id": m.GroupID, "user_id": m.UserID, "group_role": m.GroupRole,
		"times_picked": m.TimesPicked, "picker_skips": m.PickerSkips,
		models.BASE_CREATED_AT: m.CreatedAt, models.BASE_UPDATED_AT: m.UpdatedAt,
	}))
}

func (dao *DAO) GetGroupMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	return getGeneric[models.GroupMember](ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(models.GroupMemberColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID, "user_id": userID})).
		WithLimit(1))
}

func (dao *DAO) ListGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	return listGeneric[models.GroupMember](ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(models.GroupMemberColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID})))
}

func (dao *DAO) ListGroupsForUser(ctx context.Context, userID string) ([]*models.GroupMember, error) {
	return listGeneric[models.GroupMember](ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(models.GroupMemberColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"user_id": userID})))
}

// ListUserGroupSummariesForUserIDs returns group summaries for the given user IDs, ordered by
// user then group name
func (dao *DAO) ListUserGroupSummariesForUserIDs(ctx context.Context, userIDs []string) ([]*UserGroupSummaryRow, error) {
	if len(userIDs) == 0 {
		return []*UserGroupSummaryRow{}, nil
	}

	gm := models.GROUP_MEMBER_TABLE
	g := models.GROUP_TABLE

	return listGeneric[UserGroupSummaryRow](ctx, dao, *newBuilderOptions(gm).
		WithColumns(
			gm+".user_id AS user_id",
			g+"."+models.BASE_ID+" AS id",
			g+".name AS name",
			gm+".group_role AS group_role",
			"(SELECT COUNT(*) FROM "+gm+" gm_count WHERE gm_count.group_id = "+g+"."+models.BASE_ID+") AS member_count",
		).
		WithJoin(g, g+"."+models.BASE_ID+" = "+gm+".group_id").
		SetDbOpts(NewOptions().
			WithWhere(squirrel.Eq{gm + ".user_id": userIDs}).
			WithOrderBy(gm+".user_id ASC", g+"."+models.BASE_CREATED_AT+" desc")))
}

func (dao *DAO) DeleteGroupMember(ctx context.Context, groupID, userID string) error {
	if groupID == "" || userID == "" {
		return utils.ErrWhere
	}
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).SetDbOpts(
		NewOptions().WithWhere(squirrel.Eq{"group_id": groupID, "user_id": userID}))
	sqlStr, args, _ := deleteBuilder(*builderOpts)
	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}

func (dao *DAO) IncrementTimesPicked(ctx context.Context, memberID string) error {
	now := types.NowDateTime().String()
	sql := `UPDATE group_members SET times_picked = times_picked + 1, updated_at = ? WHERE id = ?`
	_, err := dao.db.ExecContext(ctx, sql, now, memberID)
	return err
}

func (dao *DAO) CreateJoinRequest(ctx context.Context, r *models.GroupJoinRequest) error {
	if r.ID == "" {
		r.RefreshId()
	}
	r.RefreshCreatedAt()
	r.RefreshUpdatedAt()
	if r.Status == "" {
		r.Status = types.JoinPending
	}
	return createGeneric(ctx, dao, *newBuilderOptions(models.JOIN_REQUEST_TABLE).WithData(map[string]interface{}{
		models.BASE_ID: r.ID, "group_id": r.GroupID, "user_id": r.UserID, "status": r.Status,
		models.BASE_CREATED_AT: r.CreatedAt, models.BASE_UPDATED_AT: r.UpdatedAt,
	}))
}

func (dao *DAO) GetJoinRequest(ctx context.Context, id string) (*models.GroupJoinRequest, error) {
	return getGeneric[models.GroupJoinRequest](ctx, dao, *newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.GroupJoinRequestColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})).WithLimit(1))
}

func (dao *DAO) GetJoinRequestByUser(ctx context.Context, groupID, userID string) (*models.GroupJoinRequest, error) {
	return getGeneric[models.GroupJoinRequest](ctx, dao, *newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.GroupJoinRequestColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID, "user_id": userID})).WithLimit(1))
}

// ListMemberGroupIDsForUser returns group IDs the user belongs to from the given set
func (dao *DAO) ListMemberGroupIDsForUser(ctx context.Context, userID string, groupIDs []string) ([]string, error) {
	if len(groupIDs) == 0 {
		return []string{}, nil
	}

	type groupIDRow struct {
		GroupID string `db:"group_id"`
	}

	rows, err := listGeneric[groupIDRow](ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns("group_id").
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{
			"user_id":  userID,
			"group_id": groupIDs,
		})))
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.GroupID)
	}

	return ids, nil
}

// ListPendingJoinGroupIDsForUser returns group IDs with a pending join request for the user
func (dao *DAO) ListPendingJoinGroupIDsForUser(ctx context.Context, userID string, groupIDs []string) ([]string, error) {
	if len(groupIDs) == 0 {
		return []string{}, nil
	}

	type groupIDRow struct {
		GroupID string `db:"group_id"`
	}

	rows, err := listGeneric[groupIDRow](ctx, dao, *newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns("group_id").
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{
			"user_id":  userID,
			"group_id": groupIDs,
			"status":   types.JoinPending,
		})))
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.GroupID)
	}

	return ids, nil
}

func (dao *DAO) ListPendingJoinRequests(ctx context.Context, groupID string) ([]*models.GroupJoinRequest, error) {
	return listGeneric[models.GroupJoinRequest](ctx, dao, *newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.GroupJoinRequestColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID, "status": types.JoinPending})))
}

func (dao *DAO) UpdateJoinRequestStatus(ctx context.Context, id string, status types.JoinRequestStatus) error {
	now := types.NowDateTime().String()
	_, err := updateGeneric(ctx, dao, *newBuilderOptions(models.JOIN_REQUEST_TABLE).WithData(map[string]interface{}{
		"status": status, models.BASE_UPDATED_AT: now,
	}).SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})))
	return err
}

// DeletePendingJoinRequest removes a pending join request for the given group and user
func (dao *DAO) DeletePendingJoinRequest(ctx context.Context, groupID, userID string) error {
	if groupID == "" || userID == "" {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).SetDbOpts(
		NewOptions().WithWhere(squirrel.Eq{
			"group_id": groupID,
			"user_id":  userID,
			"status":   types.JoinPending,
		}))
	sqlStr, args, _ := deleteBuilder(*builderOpts)
	_, err := dao.db.ExecContext(ctx, sqlStr, args...)

	return err
}

func (dao *DAO) CreateRound(ctx context.Context, r *models.Round) error {
	if r.ID == "" {
		r.RefreshId()
	}
	r.RefreshCreatedAt()
	r.RefreshUpdatedAt()
	data := map[string]interface{}{
		models.BASE_ID: r.ID, "group_id": r.GroupID, "round_date": r.RoundDate,
		"picker_user_id": r.PickerUserID, "status": r.Status,
		models.BASE_CREATED_AT: r.CreatedAt, models.BASE_UPDATED_AT: r.UpdatedAt,
	}
	if r.WordHash != nil {
		data["word_hash"] = *r.WordHash
	}
	if r.WordPlain != nil {
		data["word_plain"] = *r.WordPlain
	}
	return createGeneric(ctx, dao, *newBuilderOptions(models.ROUND_TABLE).WithData(data))
}

func (dao *DAO) GetRound(ctx context.Context, id string) (*models.Round, error) {
	return getGeneric[models.Round](ctx, dao, *newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})).WithLimit(1))
}

func (dao *DAO) GetCurrentRound(ctx context.Context, groupID, roundDate string) (*models.Round, error) {
	return getGeneric[models.Round](ctx, dao, *newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID, "round_date": roundDate})).
		WithLimit(1))
}

func (dao *DAO) UpdateRound(ctx context.Context, r *models.Round) error {
	r.RefreshUpdatedAt()
	data := map[string]interface{}{
		"status": r.Status, models.BASE_UPDATED_AT: r.UpdatedAt,
	}
	if r.WordHash != nil {
		data["word_hash"] = *r.WordHash
	}
	if r.WordPlain != nil {
		data["word_plain"] = *r.WordPlain
	}
	_, err := updateGeneric(ctx, dao, *newBuilderOptions(models.ROUND_TABLE).WithData(data).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: r.ID})))
	return err
}

func (dao *DAO) ListActiveRounds(ctx context.Context) ([]*models.Round, error) {
	return listGeneric[models.Round](ctx, dao, *newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Or{
			squirrel.Eq{"status": types.RoundAwaitingWord},
			squirrel.Eq{"status": types.RoundActive},
		})))
}

func (dao *DAO) CreateGuess(ctx context.Context, g *models.Guess) error {
	if g.ID == "" {
		g.RefreshId()
	}
	g.RefreshCreatedAt()
	g.RefreshUpdatedAt()
	return createGeneric(ctx, dao, *newBuilderOptions(models.GUESS_TABLE).WithData(map[string]interface{}{
		models.BASE_ID: g.ID, "round_id": g.RoundID, "user_id": g.UserID,
		"attempts_used": g.AttemptsUsed, "solved": g.Solved, "rows_json": g.RowsJSON,
		"score": g.Score, "finished": g.Finished,
		"first_guess_at": g.FirstGuessAt, "completed_at": g.CompletedAt,
		models.BASE_CREATED_AT: g.CreatedAt, models.BASE_UPDATED_AT: g.UpdatedAt,
	}))
}

func (dao *DAO) GetGuess(ctx context.Context, roundID, userID string) (*models.Guess, error) {
	return getGeneric[models.Guess](ctx, dao, *newBuilderOptions(models.GUESS_TABLE).
		WithColumns(models.GuessColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"round_id": roundID, "user_id": userID})).
		WithLimit(1))
}

func (dao *DAO) ListGuessesForRound(ctx context.Context, roundID string) ([]*models.Guess, error) {
	return listGeneric[models.Guess](ctx, dao, *newBuilderOptions(models.GUESS_TABLE).
		WithColumns(models.GuessColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"round_id": roundID})))
}

func (dao *DAO) UpdateGuess(ctx context.Context, g *models.Guess) error {
	g.RefreshUpdatedAt()
	_, err := updateGeneric(ctx, dao, *newBuilderOptions(models.GUESS_TABLE).WithData(map[string]interface{}{
		"attempts_used": g.AttemptsUsed, "solved": g.Solved, "rows_json": g.RowsJSON,
		"score": g.Score, "finished": g.Finished,
		"first_guess_at": g.FirstGuessAt, "completed_at": g.CompletedAt,
		models.BASE_UPDATED_AT: g.UpdatedAt,
	}).SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: g.ID})))
	return err
}

func (dao *DAO) Leaderboard(ctx context.Context, groupID string) ([]map[string]interface{}, error) {
	query := `
SELECT u.id AS user_id, u.display_name, COALESCE(SUM(g.score), 0) AS total_score,
       COUNT(g.id) AS rounds_played
FROM group_members gm
JOIN users u ON u.id = gm.user_id
LEFT JOIN rounds r ON r.group_id = gm.group_id AND r.status = 'completed'
LEFT JOIN guesses g ON g.round_id = r.id AND g.user_id = gm.user_id
WHERE gm.group_id = ?
GROUP BY u.id, u.display_name
ORDER BY total_score DESC`
	rows, err := dao.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var userID, displayName string
		var totalScore, roundsPlayed int
		if err := rows.Scan(&userID, &displayName, &totalScore, &roundsPlayed); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"userId": userID, "displayName": displayName,
			"totalScore": totalScore, "roundsPlayed": roundsPlayed,
		})
	}
	return out, rows.Err()
}

func (dao *DAO) ListAllGroups(ctx context.Context) ([]*models.Group, error) {
	return dao.ListGroups(ctx, NewOptions())
}

// DeleteGroups deletes records from the groups table
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteGroups(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.GROUP_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}

// PickNextPicker returns member with minimum times_picked (random tie-break done in service)
func (dao *DAO) MembersWithMinPicks(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	query := `SELECT id, group_id, user_id, group_role, times_picked, picker_skips, created_at, updated_at
FROM group_members WHERE group_id = ? AND times_picked = (
  SELECT MIN(times_picked) FROM group_members WHERE group_id = ?)`
	rows, err := dao.db.QueryContext(ctx, query, groupID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []*models.GroupMember
	for rows.Next() {
		m := &models.GroupMember{}
		if err := rows.Scan(&m.ID, &m.GroupID, &m.UserID, &m.GroupRole, &m.TimesPicked, &m.PickerSkips, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (dao *DAO) CountGroupMembers(ctx context.Context, groupID string) (int, error) {
	return countGeneric(ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID})))
}

