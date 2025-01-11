package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/madsbv/gator/internal/command"
	"github.com/madsbv/gator/internal/config"
	"github.com/madsbv/gator/internal/database"
	"github.com/madsbv/gator/internal/state"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	conf, err := config.Read()
	if err != nil {
		return errors.Join(errors.New("Error reading config file"), err)
	}

	db, err := sql.Open("postgres", conf.DbUrl)
	dbQueries := database.New(db)

	s := state.State{Config: conf, Db: *dbQueries}
	cmds := command.Commands{Map: make(map[string]func(*state.State, command.Command) error)}
	cmds.Register("login", command.HandlerLogin)
	cmds.Register("register", command.HandlerRegister)
	cmds.Register("reset", command.HandlerReset)
	cmds.Register("agg", command.HandlerAgg)
	cmds.Register("addfeed", command.HandlerAddFeed)

	args := os.Args
	if len(args) < 2 {
		return errors.New("Not enough arguments: missing command name")
	}

	cmd := command.Command{
		Name: args[1],
		Args: args[2:],
	}

	return cmds.Run(&s, cmd)
}
