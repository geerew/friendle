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
