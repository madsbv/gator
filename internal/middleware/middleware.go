package middleware

import (
	"github.com/madsbv/gator/internal/command"
	"github.com/madsbv/gator/internal/database"
	"github.com/madsbv/gator/internal/state"
)

func LoggedIn(handler func(s *state.State, cmd command.Command, user database.User) error) func(*state.State, command.Command) error {

	inner := func(s *state.State, cmd command.Command) error {
		user, err := s.CurrentUser()
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}

	return inner
}
