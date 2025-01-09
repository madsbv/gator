package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madsbv/gator/internal/database"
	"github.com/madsbv/gator/internal/rss"
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

	name := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		return errors.New(fmt.Sprintf("User %s not found\n", name))

	}

	err = s.Config.SetUser(name)
	if err == nil {
		fmt.Printf("User %s logged in\n", name)
	}
	return err
}

func HandlerRegister(s *state.State, cmd Command) error {
	if cmd.Name != "register" || len(cmd.Args) != 1 {
		return errors.New("Missing user name to register")
	}

	id := uuid.New()
	currentTime := time.Now()
	params := database.CreateUserParams{ID: id, CreatedAt: currentTime, UpdatedAt: currentTime, Name: cmd.Args[0]}

	user, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		return errors.New(fmt.Sprintf("User %s already exists", cmd.Args[0]))
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return errors.New("Error setting user")
	}

	fmt.Printf("User %s created\n", user.Name)
	fmt.Println(user)

	return nil
}

func HandlerReset(s *state.State, cmd Command) error {
	if cmd.Name != "reset" || len(cmd.Args) != 0 {
		return errors.New("Unexpected arguments")
	}

	return s.Db.DeleteAllUsers(context.Background())
}

func HandlerAgg(s *state.State, cmd Command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return errors.New(fmt.Sprintf("Error fetching feed at %s: %s", url, err))
	}
	fmt.Println(feed)
	return nil
}
