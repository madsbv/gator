package command

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

	userName := cmd.Args[0]

	user, err := s.Db.CreateUser(context.Background(), userName)
	if err != nil {
		return errors.New(fmt.Sprintf("User %s already exists", userName))
	}

	if userName != user.Name {
		return errors.New(fmt.Sprintf("username %s was saved to the database, but the database returned username %s instead.", userName, user.Name))
	}

	err = s.Config.SetUser(userName)
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
	if cmd.Name != "agg" {
		return errors.New("Command name mismatch")
	}

	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return errors.New(fmt.Sprintf("Error fetching feed at %s: %s", url, err))
	}
	fmt.Println(feed)
	return nil
}

func HandlerAddFeed(s *state.State, cmd Command) error {
	if cmd.Name != "addfeed" || len(cmd.Args) != 2 {
		return errors.New("Command name mismatch or wrong number of arguments. The addfeeds commmand takes two arguments, the name and the url of the feed to add.")
	}

	user_name := s.Config.CurrentUserName
	ctx := context.Background()
	user, err := s.Db.GetUser(ctx, user_name)
	if err != nil {
		return errors.New("Error retrieving currently logged in user from database")
	}

	feedParams := database.CreateFeedParams{
		Name:   sql.NullString{String: cmd.Args[0], Valid: true},
		Url:    cmd.Args[1],
		UserID: user.ID,
	}

	feed, err := s.Db.CreateFeed(ctx, feedParams)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating feed: %s", err))
	}

	fmt.Println(feed)
	return nil
}

func HandlerGetAllFeeds(s *state.State, cmd Command) error {
	if cmd.Name != "feeds" {
		return errors.New("Command name mismatch")
	}

	feeds, err := s.Db.GetAllFeedsWithUsernames(context.Background())
	if err != nil {
		return errors.New(fmt.Sprintf("Error fetching feeds from database: %s", err))
	}

	fmt.Println(feeds)
	return nil
}
