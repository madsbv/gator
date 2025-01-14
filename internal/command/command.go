package command

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
	if err := cmd.verify("login", 1); err != nil {
		return err
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
	if err := cmd.verify("register", 1); err != nil {
		return err
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
	if err := cmd.verify("reset", 0); err != nil {
		return err
	}

	return s.Db.DeleteAllUsers(context.Background())
}

func HandlerAgg(s *state.State, cmd Command) error {
	if err := cmd.verify("agg", 1); err != nil {
		return err
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return errors.New(fmt.Sprintf("Error parsing passed in time duration: %s\n", err))
	}

	fmt.Printf("Collecting feeds every %s\n", duration)

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		err = rss.ScrapeFeeds(s)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func HandlerAddFeed(s *state.State, cmd Command, user database.User) error {
	if err := cmd.verify("addfeed", 2); err != nil {
		return err
	}

	feedParams := database.CreateFeedParams{
		Name:   sql.NullString{String: cmd.Args[0], Valid: true},
		Url:    cmd.Args[1],
		UserID: user.ID,
	}

	feed, err := s.Db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating feed: %s", err))
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating feed_follow entry for user %s and feed %s", user.Name, feed.Name))
	}

	fmt.Println(feed)
	return nil
}

func HandlerGetAllFeeds(s *state.State, cmd Command) error {
	if err := cmd.verify("feeds", 0); err != nil {
		return err
	}

	feeds, err := s.Db.GetAllFeedsWithUsernames(context.Background())
	if err != nil {
		return errors.New(fmt.Sprintf("Error fetching feeds from database: %s", err))
	}

	for _, f := range feeds {
		fmt.Printf("Name: %s\nURL: %s\nAdded by: %s\n\n", f.Name.String, f.Url, f.UserName)
	}
	return nil
}

func HandlerFollow(s *state.State, cmd Command, user database.User) error {
	if err := cmd.verify("follow", 1); err != nil {
		return err
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return errors.New(fmt.Sprintf("Error getting feed from database, url: %s", cmd.Args[0]))
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})

	if err != nil {
		return errors.New(fmt.Sprintf("Error following feed %s as user %s", feed.Name.String, user.Name))
	}

	fmt.Printf("User %s successfully followed feed %s\n", feedFollow.UserName, feedFollow.FeedName.String)
	return nil
}

func HandlerFollowing(s *state.State, cmd Command, user database.User) error {
	if err := cmd.verify("following", 0); err != nil {
		return err
	}

	follows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return errors.New(fmt.Sprintf("Error getting feed follows for user %s from database: %s", user.Name, err))
	}

	fmt.Printf("User %s follows the following feeds:\n", user.Name)
	for _, f := range follows {
		fmt.Println(f.FeedName.String)
	}

	return nil
}

func HandlerUnfollow(s *state.State, cmd Command, user database.User) error {
	if err := cmd.verify("unfollow", 1); err != nil {
		return err
	}

	err := s.Db.DeleteFeedFollowByUrl(context.Background(), database.DeleteFeedFollowByUrlParams{Url: cmd.Args[0], UserID: user.ID})
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Command) verify(name string, numArgs int) error {
	if cmd.Name != name || len(cmd.Args) != numArgs {
		return errors.New(fmt.Sprintf("Command name mismatch or wrong number of arguments, expected %d arguments for command name %s", numArgs, name))
	}

	return nil
}
