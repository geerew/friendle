package api

import (
	"net/http"
	"testing"

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
