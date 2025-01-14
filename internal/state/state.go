package state

import (
	"context"
	"errors"
	"fmt"

	"github.com/madsbv/gator/internal/config"
	"github.com/madsbv/gator/internal/database"
)

type State struct {
	Config config.Config
	Db     database.Queries
}

func (s *State) CurrentUser() (database.User, error) {
	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error retrieving currently logged in user from database: %s", err))
	}

	return user, err
}
