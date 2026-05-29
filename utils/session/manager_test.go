package session

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"testing"
	"time"

	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/filesystem"
	"github.com/geerew/friendle/utils/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	fs "github.com/gofiber/fiber/v2/middleware/session"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func setup(tb testing.TB) (database.Database, context.Context) {
	tb.Helper()

	dbManager, err := database.NewSQLite(&database.SQLiteConfig{
		DataDir: "./oc_data",
		FS:   filesystem.New(afero.NewMemMapFs()),
		Testing: true,
	})

	require.NoError(tb, err)
	require.NotNil(tb, dbManager)

	return dbManager.DataDb, context.Background()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_UpdateSessionRoleForUser(t *testing.T) {
	// Test successfully updating the session role for a user
	t.Run("success", func(t *testing.T) {
		db, ctx := setup(t)
		manager := New(db, fs.Config{}, NewSqliteStorage(db, time.Hour))

		for i := range 3 {
			session := &models.Session{
				ID:      fmt.Sprintf("session-%d", i),
				UserId:  "user-123",
				Expires: time.Now().Add(24 * time.Hour).Unix(),
			}

			values := map[string]interface{}{
				"role": types.SiteRoleUser.String(),
			}

			var out bytes.Buffer
			require.NoError(t, gob.NewEncoder(&out).Encode(values))
			session.Data = out.Bytes()

			require.NoError(t, manager.dao.CreateOrReplaceSession(ctx, session))
		}

		require.NoError(t, manager.UpdateSessionRoleForUser("user-123", types.SiteRoleAdmin))

		records, err := manager.dao.ListSessions(ctx, nil)
		require.NoError(t, err)
		require.Len(t, records, 3)

		for _, record := range records {
			buf := bytes.NewBuffer(record.Data)
			var values map[string]interface{}
			require.NoError(t, gob.NewDecoder(buf).Decode(&values))
			require.Equal(t, types.SiteRoleAdmin.String(), values["role"])
		}
	})

	// Test error due to invalid user ID
	t.Run("invalid user ID", func(t *testing.T) {
		db, _ := setup(t)
		manager := New(db, fs.Config{}, NewSqliteStorage(db, time.Hour))

		require.ErrorIs(t, manager.UpdateSessionRoleForUser("", types.SiteRoleAdmin), utils.ErrUserId)
	})
}
