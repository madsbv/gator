package command

import (
	"errors"
	"fmt"

	"github.com/madsbv/gator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Map map[string]func(*state.State, Command) error
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.Map[name] = f
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	f, exists := c.Map[cmd.Name]
	if !exists {
		return errors.New("Command doesn't exist")
	}

	return f(s, cmd)
}

func HandlerLogin(s *state.State, cmd Command) error {
	if cmd.Name != "login" || len(cmd.Args) != 1 {
		return errors.New("Invalid login command")
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err == nil {
		fmt.Printf("User '%s' logged in\n", cmd.Args[0])
	}
	return err
}
