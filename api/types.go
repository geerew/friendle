package api

import (
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
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

// adminUserResponseHelper maps admin users to API responses
func adminUserResponseHelper(users []*models.User, groupsByUser map[string][]*models.GroupMember) []*adminUserResponse {
	responses := make([]*adminUserResponse, 0, len(users))
	for _, user := range users {
		members := groupsByUser[user.ID]
		groups := userGroupSummaryResponsesFromMembers(members)
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

// userGroupSummaryResponsesFromMembers maps group memberships to API responses
func userGroupSummaryResponsesFromMembers(members []*models.GroupMember) []*userGroupSummaryResponse {
	if len(members) == 0 {
		return []*userGroupSummaryResponse{}
	}

	responses := make([]*userGroupSummaryResponse, 0, len(members))
	for _, member := range members {
		if member.Group == nil {
			continue
		}

		responses = append(responses, &userGroupSummaryResponse{
			ID:          member.Group.ID,
			Name:        member.Group.Name,
			MemberCount: member.Group.MemberCount,
			GroupRole:   member.GroupRole,
		})
	}

	return responses
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupMembersByUserID groups memberships by user ID
func groupMembersByUserID(members []*models.GroupMember) map[string][]*models.GroupMember {
	byUser := make(map[string][]*models.GroupMember)
	for _, member := range members {
		byUser[member.UserID] = append(byUser[member.UserID], member)
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
	Name *string `json:"name"`
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
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	GroupRole types.GroupRole        `json:"groupRole"`
	Members       []*groupMemberResponse `json:"members,omitempty"`
	JoinRequests  []*joinRequestResponse `json:"joinRequests,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupResponseHelper maps a group to an API response
func groupResponseHelper(g *models.Group, role types.GroupRole) *groupResponse {
	resp := &groupResponse{
		ID:        g.ID,
		Name:      g.Name,
		GroupRole: role,
	}
	if len(g.Members) > 0 {
		resp.Members = groupMemberResponsesFromModels(g.Members)
	}
	if len(g.JoinRequests) > 0 {
		resp.JoinRequests = joinRequestResponsesFromModels(g.JoinRequests)
	}

	return resp
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupMemberResponse is a group member returned by the API
type groupMemberResponse struct {
	ID          string          `json:"id"`
	UserID      string          `json:"userId"`
	DisplayName string          `json:"displayName"`
	GroupRole   types.GroupRole `json:"groupRole"`
	TimesPicked int             `json:"timesPicked"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupMemberResponsesFromModels maps group members to API responses
func groupMemberResponsesFromModels(members []*models.GroupMember) []*groupMemberResponse {
	responses := make([]*groupMemberResponse, 0, len(members))
	for _, m := range members {
		name := m.UserID
		if m.User != nil {
			name = m.User.DisplayName
		}
		responses = append(responses, &groupMemberResponse{
			ID: m.ID, UserID: m.UserID, DisplayName: name,
			GroupRole: m.GroupRole, TimesPicked: m.TimesPicked,
		})
	}

	return responses
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// joinRequestResponsesFromModels maps join requests to API responses
func joinRequestResponsesFromModels(requests []*models.GroupJoinRequest) []*joinRequestResponse {
	responses := make([]*joinRequestResponse, 0, len(requests))
	for _, jr := range requests {
		responses = append(responses, &joinRequestResponse{ID: jr.ID, Status: jr.Status})
	}

	return responses
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundSummaryResponse is a round without player guess details
type roundSummaryResponse struct {
	ID           string            `json:"id"`
	RoundDate    string            `json:"roundDate"`
	Status       types.RoundStatus `json:"status"`
	PickerUserID string            `json:"pickerUserId,omitempty"`
	Word         string            `json:"word,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundSummaryResponseHelper maps a round to a summary response
func roundSummaryResponseHelper(round *models.Round, includeWord bool) *roundSummaryResponse {
	resp := &roundSummaryResponse{
		ID: round.ID, RoundDate: round.RoundDate, Status: round.Status,
	}
	if includeWord && round.WordPlain != nil {
		resp.Word = *round.WordPlain
	}
	if includeWord {
		resp.PickerUserID = round.PickerUserID
	}

	return resp
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
	Word    string               `json:"word"`
	Result  []types.TileState    `json:"result"`
	Outcome types.GuessOutcome   `json:"outcome"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundResponse is a round with optional participations for the API
type roundResponse struct {
	ID             string                      `json:"id,omitempty"`
	GroupID        string                      `json:"groupId,omitempty"`
	RoundDate      string                      `json:"roundDate,omitempty"`
	Status         string                      `json:"status"`
	PickerUserID      string                      `json:"pickerUserId,omitempty"`
	PickerDisplayName string                      `json:"pickerDisplayName,omitempty"`
	Word              string                      `json:"word,omitempty"`
	YourRole       string                      `json:"yourRole,omitempty"`
	Participations []*participationResponse    `json:"participations,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// participationResponse is one member's play in a round
type participationResponse struct {
	ID           string             `json:"id"`
	UserID       string             `json:"userId"`
	DisplayName  string             `json:"displayName,omitempty"`
	Solved       bool               `json:"solved"`
	Finished     bool               `json:"finished"`
	Score        int                `json:"score"`
	FirstGuessAt *string            `json:"firstGuessAt,omitempty"`
	CompletedAt  *string            `json:"completedAt,omitempty"`
	Guesses      []guessRowResponse `json:"guesses,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundResponseHelper maps a round model to an API response
func roundResponseHelper(round *models.Round, viewerID string, siteAdmin, groupAdmin bool) *roundResponse {
	if round == nil {
		return &roundResponse{Status: "none"}
	}

	filterRoundForViewer(round, viewerID, siteAdmin, groupAdmin)

	resp := &roundResponse{
		ID:        round.ID,
		GroupID:   round.GroupID,
		RoundDate: round.RoundDate,
		Status:    string(round.Status),
		YourRole:  yourRoleInRound(round, viewerID),
	}
	if viewerCanSeeRoundWord(round, viewerID, siteAdmin, groupAdmin) {
		resp.Word = *round.WordPlain
		resp.PickerUserID = round.PickerUserID
	} else if round.Status.IsRevealed() {
		resp.PickerUserID = round.PickerUserID
	}
	if round.Picker != nil && (viewerCanSeeRoundWord(round, viewerID, siteAdmin, groupAdmin) || round.Status.IsRevealed()) {
		resp.PickerDisplayName = round.Picker.DisplayName
	}

	if len(round.Participations) > 0 {
		resp.Participations = participationResponsesFromModels(round.Participations)
	}

	return resp
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// participationResponsesFromModels maps participations to API responses
func participationResponsesFromModels(participations []*models.RoundParticipation) []*participationResponse {
	responses := make([]*participationResponse, 0, len(participations))
	for _, p := range participations {
		name := p.UserID
		if p.User != nil {
			name = p.User.DisplayName
		}
		responses = append(responses, &participationResponse{
			ID: p.ID, UserID: p.UserID, DisplayName: name,
			Solved: p.Solved, Finished: p.Finished, Score: p.Score,
			FirstGuessAt: p.FirstGuessAt, CompletedAt: p.CompletedAt,
			Guesses:      guessRowResponsesFromModels(p.Guesses),
		})
	}

	return responses
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// guessRowResponsesFromModels maps guess rows to API responses
func guessRowResponsesFromModels(guesses []*models.Guess) []guessRowResponse {
	rows := make([]guessRowResponse, 0, len(guesses))
	for _, g := range guesses {
		rows = append(rows, guessRowResponse{
			Word: g.Word, Result: []types.TileState(g.Result), Outcome: g.Outcome,
		})
	}

	return rows
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupRoundGuessResponse is returned after submitting a guess
type groupRoundGuessResponse struct {
	Result   []types.TileState `json:"result"`
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

// groupRoundRevealResponseHelper maps a round to a reveal response
func groupRoundRevealResponseHelper(round *models.Round) *groupRoundRevealResponse {
	resp := &groupRoundRevealResponse{
		Status:       round.Status,
		PickerUserID: round.PickerUserID,
		Guesses:      []groupRoundRevealGuessResponse{},
	}
	if round.WordPlain != nil {
		resp.Word = *round.WordPlain
	}

	for _, p := range round.Participations {
		if p.UserID == round.PickerUserID {
			continue
		}

		name := p.UserID
		if p.User != nil {
			name = p.User.DisplayName
		}

		resp.Guesses = append(resp.Guesses, groupRoundRevealGuessResponse{
			UserID: p.UserID, DisplayName: name,
			AttemptsUsed: len(p.Guesses), Solved: p.Solved, Score: p.Score,
		})
	}

	return resp
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

// groupSearchResponsesFromGroups maps groups to search API responses with membership flags
func groupSearchResponsesFromGroups(groups []*models.Group, memberGroupIDs map[string]struct{}, pendingGroupIDs map[string]struct{}) []*groupSearchResponse {
	if len(groups) == 0 {
		return []*groupSearchResponse{}
	}

	responses := make([]*groupSearchResponse, 0, len(groups))
	for _, group := range groups {
		_, isMember := memberGroupIDs[group.ID]
		_, joinPending := pendingGroupIDs[group.ID]
		responses = append(responses, &groupSearchResponse{
			ID:          group.ID,
			Name:        group.Name,
			MemberCount: group.MemberCount,
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

// adminGroupResponseHelper maps groups to admin API responses
func adminGroupResponseHelper(groups []*models.Group) []*adminGroupResponse {
	responses := make([]*adminGroupResponse, 0, len(groups))
	for _, g := range groups {
		responses = append(responses, &adminGroupResponse{
			ID: g.ID, Name: g.Name, MemberCount: g.MemberCount,
		})
	}

	return responses
}
