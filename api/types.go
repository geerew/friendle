package api

import (
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type userRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Role        string `json:"siteRole"`
}

type userResponse struct {
	ID          string         `json:"id"`
	Username    string         `json:"username"`
	DisplayName string         `json:"displayName"`
	SiteRole    types.SiteRole `json:"siteRole"`
}

// userResponseHelper maps users to API responses
func userResponseHelper(users []*models.User) []*userResponse {
	responses := []*userResponse{}
	for _, user := range users {
		responses = append(responses, &userResponse{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			SiteRole:    user.SiteRole,
		})
	}

	return responses
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type adminUserResponse struct {
	ID          string                      `json:"id"`
	Username    string                      `json:"username"`
	DisplayName string                      `json:"displayName"`
	SiteRole    types.SiteRole              `json:"siteRole"`
	GroupCount  int                         `json:"groupCount"`
	Groups      []*userGroupSummaryResponse `json:"groups"`
}

// adminUserResponseHelper maps admin user rows to API responses
func adminUserResponseHelper(
	users []*models.AdminUserListRow,
	groupsByUser map[string][]*models.UserGroupSummaryRow,
) []*adminUserResponse {
	responses := make([]*adminUserResponse, 0, len(users))
	for _, user := range users {
		groupRows := groupsByUser[user.ID]
		groups := userGroupSummaryResponsesFromRows(groupRows)
		responses = append(responses, &adminUserResponse{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			SiteRole:    user.SiteRole,
			GroupCount:  len(groups),
			Groups:      groups,
		})
	}

	return responses
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type userGroupSummaryResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	MemberCount int             `json:"memberCount"`
	GroupRole   types.GroupRole `json:"groupRole,omitempty"`
}

// userGroupSummaryResponsesFromRows maps group summary rows to API responses
func userGroupSummaryResponsesFromRows(rows []*models.UserGroupSummaryRow) []*userGroupSummaryResponse {
	if len(rows) == 0 {
		return []*userGroupSummaryResponse{}
	}

	responses := make([]*userGroupSummaryResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, &userGroupSummaryResponse{
			ID:          row.ID,
			Name:        row.Name,
			MemberCount: row.MemberCount,
			GroupRole:   row.GroupRole,
		})
	}

	return responses
}

// userGroupSummariesByUserID groups summary rows by user ID
func userGroupSummariesByUserID(rows []*models.UserGroupSummaryRow) map[string][]*models.UserGroupSummaryRow {
	byUser := make(map[string][]*models.UserGroupSummaryRow)
	for _, row := range rows {
		byUser[row.UserID] = append(byUser[row.UserID], row)
	}

	return byUser
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type createGroupRequest struct {
	Name string `json:"name"`
}

type updateGroupRequest struct {
	Name          *string `json:"name"`
	IntervalHours *int    `json:"intervalHours"`
}

type createGroupJoinRequest struct {
	UserID string `json:"userId"`
}

type adminAddGroupMemberRequest struct {
	UserID    string `json:"userId"`
	GroupRole string `json:"groupRole"`
}

type submitRoundWordRequest struct {
	Word string `json:"word"`
}

type submitRoundGuessRequest struct {
	Word string `json:"word"`
}

// groupResponseHelper builds a group detail response map
func groupResponseHelper(g *models.Group, role types.GroupRole) fiber.Map {
	return fiber.Map{
		"id": g.ID, "name": g.Name, "intervalHours": g.IntervalHours,
		"timezone": g.Timezone, "groupRole": role,
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type groupSearchResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"memberCount"`
	IsMember    bool   `json:"isMember"`
	JoinPending bool   `json:"joinPending"`
}

// groupSearchResponsesFromRows maps search rows to API responses with membership flags
func groupSearchResponsesFromRows(
	rows []*models.GroupSearchRow,
	memberGroupIDs map[string]struct{},
	pendingGroupIDs map[string]struct{},
) []*groupSearchResponse {
	if len(rows) == 0 {
		return []*groupSearchResponse{}
	}

	responses := make([]*groupSearchResponse, 0, len(rows))
	for _, row := range rows {
		_, isMember := memberGroupIDs[row.ID]
		_, joinPending := pendingGroupIDs[row.ID]
		responses = append(responses, &groupSearchResponse{
			ID:          row.ID,
			Name:        row.Name,
			MemberCount: row.MemberCount,
			IsMember:    isMember,
			JoinPending: joinPending,
		})
	}

	return responses
}

// stringSet converts a slice of IDs into a lookup set
func stringSet(ids []string) map[string]struct{} {
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}

	return set
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type signupStatusResponse struct {
	Enabled bool `json:"enabled"`
}

type registerRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type selfUpdateRequest struct {
	DisplayName     string `json:"displayName"`
	CurrentPassword string `json:"currentPassword"`
	Password        string `json:"password"`
}

type selfDeleteRequest struct {
	CurrentPassword string `json:"currentPassword"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type adminGroupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"memberCount"`
}

// adminGroupResponseHelper maps admin group rows to API responses
func adminGroupResponseHelper(groups []*models.AdminGroupListRow) []*adminGroupResponse {
	responses := make([]*adminGroupResponse, 0, len(groups))
	for _, g := range groups {
		responses = append(responses, &adminGroupResponse{
			ID: g.ID, Name: g.Name, MemberCount: g.MemberCount,
		})
	}

	return responses
}
