package api

import (
	"encoding/json"

	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
	"github.com/geerew/friendle/utils/wordgame"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userRequest is the body for site-admin user create and update endpoints
type userRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Role        string `json:"siteRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userResponse is a user returned by the API
type userResponse struct {
	ID          string         `json:"id"`
	Username    string         `json:"username"`
	DisplayName string         `json:"displayName"`
	SiteRole    types.SiteRole `json:"siteRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminUserResponse is a user row enriched for the site admin list
type adminUserResponse struct {
	ID          string                      `json:"id"`
	Username    string                      `json:"username"`
	DisplayName string                      `json:"displayName"`
	SiteRole    types.SiteRole              `json:"siteRole"`
	GroupCount  int                         `json:"groupCount"`
	Groups      []*userGroupSummaryResponse `json:"groups"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminUserResponseHelper maps admin user rows to API responses
func adminUserResponseHelper(users []*dao.AdminUserRow, groupsByUser map[string][]*dao.UserGroupSummaryRow) []*adminUserResponse {
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

// userGroupSummaryResponse is a group summary for a user's memberships
type userGroupSummaryResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	MemberCount int             `json:"memberCount"`
	GroupRole   types.GroupRole `json:"groupRole,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userGroupSummaryResponsesFromRows maps group summary rows to API responses
func userGroupSummaryResponsesFromRows(rows []*dao.UserGroupSummaryRow) []*userGroupSummaryResponse {
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userGroupSummariesByUserID groups summary rows by user ID
func userGroupSummariesByUserID(rows []*dao.UserGroupSummaryRow) map[string][]*dao.UserGroupSummaryRow {
	byUser := make(map[string][]*dao.UserGroupSummaryRow)
	for _, row := range rows {
		byUser[row.UserID] = append(byUser[row.UserID], row)
	}

	return byUser
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroupRequest is the body for creating a group
type createGroupRequest struct {
	Name string `json:"name"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateGroupRequest is the body for updating a group
type updateGroupRequest struct {
	Name          *string `json:"name"`
	IntervalHours *int    `json:"intervalHours"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroupJoinRequest is the optional body for a join request
type createGroupJoinRequest struct {
	UserID string `json:"userId"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminAddGroupMemberRequest is the body for a site-admin direct member add
type adminAddGroupMemberRequest struct {
	UserID    string `json:"userId"`
	GroupRole string `json:"groupRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// submitRoundWordRequest is the body for submitting the picker's word
type submitRoundWordRequest struct {
	Word string `json:"word"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// submitRoundGuessRequest is the body for submitting a guess
type submitRoundGuessRequest struct {
	Word string `json:"word"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupResponse is a group returned by the API
type groupResponse struct {
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	IntervalHours int                        `json:"intervalHours"`
	Timezone      string                     `json:"timezone"`
	GroupRole     types.GroupRole            `json:"groupRole"`
	Round         *groupRoundSummaryResponse `json:"round,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupResponseHelper maps a group to an API response
func groupResponseHelper(g *models.Group, role types.GroupRole) *groupResponse {
	return &groupResponse{
		ID:            g.ID,
		Name:          g.Name,
		IntervalHours: g.IntervalHours,
		Timezone:      g.Timezone,
		GroupRole:     role,
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupRoundSummaryResponse is a summary of today's round on a group detail response
type groupRoundSummaryResponse struct {
	Status    string          `json:"status"`
	YourRole  string          `json:"yourRole,omitempty"`
	CanReveal bool            `json:"canReveal,omitempty"`
	GroupRole types.GroupRole `json:"groupRole,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// joinRequestResponse is a join request returned by the API
type joinRequestResponse struct {
	ID     string                   `json:"id,omitempty"`
	Status types.JoinRequestStatus `json:"status"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminGroupMemberResponse is a member added directly by a site admin
type adminGroupMemberResponse struct {
	ID        string          `json:"id"`
	UserID    string          `json:"userId"`
	GroupID   string          `json:"groupId"`
	GroupRole types.GroupRole `json:"groupRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// guessRowResponse is one submitted guess row for the current round
type guessRowResponse struct {
	Word   string               `json:"word"`
	Result []wordgame.TileState `json:"result"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupRoundResponse is the current round state for the caller
type groupRoundResponse struct {
	RoundID      string             `json:"roundId,omitempty"`
	Status       string             `json:"status"`
	YourRole     string             `json:"yourRole,omitempty"`
	AttemptsUsed int                `json:"attemptsUsed,omitempty"`
	Finished     bool               `json:"finished,omitempty"`
	Solved       bool               `json:"solved,omitempty"`
	Rows         []guessRowResponse `json:"rows,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupRoundResponseHelper builds the current round response for a member
func groupRoundResponseHelper(round *models.Round, guess *models.Guess, userID string) *groupRoundResponse {
	if round == nil {
		return &groupRoundResponse{Status: "none"}
	}

	yourRole := "guesser"
	if round.PickerUserID == userID {
		yourRole = "picker"
	}

	resp := &groupRoundResponse{
		RoundID:      round.ID,
		Status:       string(round.Status),
		YourRole:     yourRole,
		AttemptsUsed: 0,
		Finished:     false,
		Rows:         []guessRowResponse{},
	}
	if guess != nil {
		var rows []guessRowResponse
		_ = json.Unmarshal([]byte(guess.RowsJSON), &rows)
		resp.AttemptsUsed = guess.AttemptsUsed
		resp.Finished = guess.Finished
		resp.Solved = guess.Solved
		resp.Rows = rows
	}

	return resp
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupRoundGuessResponse is returned after submitting a guess
type groupRoundGuessResponse struct {
	Result   []wordgame.TileState `json:"result"`
	Attempt  int                  `json:"attempt"`
	Won      bool                 `json:"won"`
	Finished bool                 `json:"finished"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupRoundRevealGuessResponse is one guesser's result on a completed round
type groupRoundRevealGuessResponse struct {
	UserID       string `json:"userId"`
	DisplayName  string `json:"displayName"`
	AttemptsUsed int    `json:"attemptsUsed"`
	Solved       bool   `json:"solved"`
	Score        int    `json:"score"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupRoundRevealResponse is the reveal payload for a finished round
type groupRoundRevealResponse struct {
	Status            types.RoundStatus               `json:"status"`
	PickerUserID      string                          `json:"pickerUserId"`
	PickerDisplayName string                          `json:"pickerDisplayName,omitempty"`
	Word              string                          `json:"word,omitempty"`
	Guesses           []groupRoundRevealGuessResponse `json:"guesses"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// versionResponse is the application version payload
type versionResponse struct {
	Version string `json:"version"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupSearchResponse is a group row enriched for name search results
type groupSearchResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"memberCount"`
	IsMember    bool   `json:"isMember"`
	JoinPending bool   `json:"joinPending"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupSearchResponsesFromRows maps search rows to API responses with membership flags
func groupSearchResponsesFromRows(rows []*dao.GroupSearchRow, memberGroupIDs map[string]struct{}, pendingGroupIDs map[string]struct{}) []*groupSearchResponse {
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// stringSet converts a slice of IDs into a lookup set
func stringSet(ids []string) map[string]struct{} {
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}

	return set
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// signupStatusResponse reports whether self-service registration is enabled
type signupStatusResponse struct {
	Enabled bool `json:"enabled"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// registerRequest is the body for user registration
type registerRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// loginRequest is the body for user login
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// selfUpdateRequest is the body for updating the authenticated user's profile
type selfUpdateRequest struct {
	DisplayName     string `json:"displayName"`
	CurrentPassword string `json:"currentPassword"`
	Password        string `json:"password"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// selfDeleteRequest is the body for deleting the authenticated user's account
type selfDeleteRequest struct {
	CurrentPassword string `json:"currentPassword"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminRecoveryRequest is the body for POST /api/admin/recovery
type adminRecoveryRequest struct {
	Token string `json:"token"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminGroupResponse is a group row enriched for the site admin list
type adminGroupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"memberCount"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminGroupResponseHelper maps admin group rows to API responses
func adminGroupResponseHelper(groups []*dao.AdminGroupRow) []*adminGroupResponse {
	responses := make([]*adminGroupResponse, 0, len(groups))
	for _, g := range groups {
		responses = append(responses, &adminGroupResponse{
			ID: g.ID, Name: g.Name, MemberCount: g.MemberCount,
		})
	}

	return responses
}
