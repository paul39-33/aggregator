package main

import (
	"fmt"
	"github.com/paul39-33/aggregator/internal/config"
)

type state struct {
	cfg	*config.Config
}

type command struct {
	name		string
	arg			[]string
}

type commands struct {
	commandList	map[string]func(*state, command) error
}

var cfg config.Config

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("login command needs a username!")
	}

	err := cfg.SetUser(cmd.arg[0])
	if err != nil {
		return fmt.Errorf("Error setting the username: %v", err)
	}

	fmt.Println("The user has been set.")

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	if command, exists := c.commandList[cmd.name]; exists {
		return command(s, cmd)
	}
	return fmt.Errorf("Command doesn't exist!")
}

func (c *commands) register(name string, f func(*state, command) error ) error {
	c.commandList[name] = f
	return nil
}


