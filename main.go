package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/madsbv/gator/internal/command"
	"github.com/madsbv/gator/internal/config"
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

	s := state.State{Config: conf}
	cmds := command.Commands{Map: make(map[string]func(*state.State, command.Command) error)}
	cmds.Register("login", command.HandlerLogin)

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
