package api

import (
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
)

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

type adminUserResponse struct {
	ID          string                     `json:"id"`
	Username    string                     `json:"username"`
	DisplayName string                     `json:"displayName"`
	SiteRole    types.SiteRole             `json:"siteRole"`
	GroupCount  int                        `json:"groupCount"`
	Groups      []*userGroupSummaryResponse `json:"groups"`
}

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

type userGroupSummaryResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	MemberCount int             `json:"memberCount"`
	GroupRole   types.GroupRole `json:"groupRole,omitempty"`
}

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

func userGroupSummariesByUserID(rows []*models.UserGroupSummaryRow) map[string][]*models.UserGroupSummaryRow {
	byUser := make(map[string][]*models.UserGroupSummaryRow)
	for _, row := range rows {
		byUser[row.UserID] = append(byUser[row.UserID], row)
	}

	return byUser
}

type groupSearchResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"memberCount"`
	IsMember    bool   `json:"isMember"`
	JoinPending bool   `json:"joinPending"`
}

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

func stringSet(ids []string) map[string]struct{} {
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}

	return set
}

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

type adminGroupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"memberCount"`
}

func adminGroupResponseHelper(groups []*models.AdminGroupListRow) []*adminGroupResponse {
	responses := make([]*adminGroupResponse, 0, len(groups))
	for _, g := range groups {
		responses = append(responses, &adminGroupResponse{
			ID: g.ID, Name: g.Name, MemberCount: g.MemberCount,
		})
	}

	return responses
}
