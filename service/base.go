package service

import (
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/utils/words"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Service is the application service root. Each field covers one domain area
type Service struct {
	Auth  *Auth
	Users *Users
	Round *Round
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deps holds shared dependencies passed to domain services
type deps struct {
	dao        *dao.DAO
	dictionary *words.Dictionary
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// New creates a Service backed by the data database
func New(db database.Database, dictionary *words.Dictionary) *Service {
	d := deps{
		dao:        dao.New(db),
		dictionary: dictionary,
	}

	return &Service{
		Auth:  newAuth(d),
		Users: newUsers(d),
		Round: newRound(d),
	}
}
