package service

import (
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/database"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Service is the application service root. Each field covers one domain area
type Service struct {
	Auth   *Auth
	Groups *Groups
	Users  *Users
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deps holds shared dependencies passed to domain services
type deps struct {
	dao *dao.DAO
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// New creates a Service backed by the data database
func New(db database.Database) *Service {
	d := deps{
		dao: dao.New(db),
	}

	return &Service{
		Auth:   newAuth(d),
		Groups: newGroups(d),
		Users:  newUsers(d),
	}
}
