package session

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
	fs "github.com/gofiber/fiber/v2/middleware/session"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Storage extends the `fiber.Storage` interface
type Storage interface {
	fiber.Storage
	DeleteUser(userId string) error
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SessionManager is a thin wrapper around the fiber session store that enables storing
// sessions in a database, instead of in-memory
type SessionManager struct {
	dao        *dao.DAO
	fiberStore *fs.Store
	storage    Storage
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// New creates a new session manager. It is essentially a wrapper around the fiber session store
func New(db database.Database, config fs.Config, storage Storage) *SessionManager {
	if storage != nil {
		config.Storage = storage
	}

	sessionManager := &SessionManager{
		dao:        dao.New(db),
		fiberStore: fs.New(config),
		storage:    storage,
	}

	return sessionManager
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Get gets the session for a user
func (s *SessionManager) Get(c *fiber.Ctx) (*fs.Session, error) {
	return s.fiberStore.Get(c)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SetSession sets the session for a user
func (s *SessionManager) SetSession(c *fiber.Ctx, userId string, userRole types.SiteRole) error {
	session, err := s.Get(c)
	if err != nil {
		return err
	}

	if err := session.Regenerate(); err != nil {
		return err
	}

	session.Set("id", userId)
	session.Set("role", userRole.String())

	return session.Save()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteAllSessions deletes all sessions from the storage
func (s *SessionManager) DeleteAllSessions() error {
	return s.fiberStore.Storage.Reset()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteSession deletes a single session based on the session ID
func (s *SessionManager) DeleteSession(c *fiber.Ctx) error {
	session, err := s.Get(c)
	if err != nil {
		return err
	}

	return session.Destroy()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteUserSessions deletes all sessions for a user
func (s *SessionManager) DeleteUserSessions(id string) error {
	return s.storage.DeleteUser(id)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateSessionRoleForUser updates the role for all sessions belonging to a user
func (s *SessionManager) UpdateSessionRoleForUser(userID string, newRole types.SiteRole) error {
	if userID == "" {
		return utils.ErrUserId
	}

	ctx := context.Background()
	sessions, err := s.dao.ListSessions(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		models.SESSION_TABLE_USER_ID: userID,
	}))
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		return nil
	}

	updatedSessions := make([]*models.Session, 0, len(sessions))
	for _, session := range sessions {
		var values map[string]interface{}
		buf := bytes.NewBuffer(session.Data)
		if err := gob.NewDecoder(buf).Decode(&values); err != nil {
			continue
		}

		values["role"] = newRole.String()

		var out bytes.Buffer
		if err := gob.NewEncoder(&out).Encode(values); err != nil {
			continue
		}

		session.Data = out.Bytes()
		updatedSessions = append(updatedSessions, session)
	}

	return s.dao.BulkUpdateSessions(ctx, updatedSessions)
}
