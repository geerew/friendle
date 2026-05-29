package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/geerew/friendle/models"
	"github.com/stretchr/testify/require"
)

func TestListMyGroups(t *testing.T) {
	router, _ := setupUser(t)

	for _, path := range []string{"/api/groups", "/api/groups/"} {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.NoError(t, err)
		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		t.Logf("path=%s status=%d body=%s", path, status, string(body))
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, "[]", string(body))
	}
}

func TestSearchGroupsOrder(t *testing.T) {
	router, ctx := setupUser(t)

	for _, name := range []string{"1", "10", "11", "2"} {
		group := &models.Group{Name: name, CreatedBy: "user", IntervalHours: 24, Timezone: "UTC"}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		time.Sleep(time.Millisecond)
	}

	params := url.Values{"q": {"1"}}
	status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups/search?"+params.Encode(), nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	var result struct {
		Items []struct {
			Name string `json:"name"`
		} `json:"items"`
		TotalItems int `json:"totalItems"`
	}
	require.NoError(t, json.Unmarshal(body, &result))
	require.Equal(t, 3, result.TotalItems)
	require.Equal(t, []string{"1", "10", "11"}, []string{
		result.Items[0].Name,
		result.Items[1].Name,
		result.Items[2].Name,
	})
}

func TestCancelJoinRequest(t *testing.T) {
	router, ctx := setupUser(t)

	group := &models.Group{Name: "Join Me", CreatedBy: "user", IntervalHours: 24, Timezone: "UTC"}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))

	createReq, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests", nil)
	require.NoError(t, err)
	status, _, err := requestHelper(t, router, createReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)

	cancelReq, err := http.NewRequest(http.MethodDelete, "/api/groups/"+group.ID+"/join-requests/me", nil)
	require.NoError(t, err)
	status, _, err = requestHelper(t, router, cancelReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	cancelAgain, err := http.NewRequest(http.MethodDelete, "/api/groups/"+group.ID+"/join-requests/me", nil)
	require.NoError(t, err)
	status, body, err := requestHelper(t, router, cancelAgain)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, status)
	require.Contains(t, string(body), "Pending request not found")
}
