package main

import (
	"fmt"
	"github.com/paul39-33/aggregator/internal/config"
	"github.com/paul39-33/aggregator/internal/database"
	"github.com/google/uuid"
	"time"
	"os"
	"context"
)

type state struct {
	db	*database.Queries
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
		return fmt.Errorf("Login command needs a name!")
	}

	//Check if user with given name exists
	if _, err := s.db.GetUser(context.Background(), cmd.arg[0]); err != nil {
		fmt.Println("User with such name doesn't exist!")
		os.Exit(1)
	}

	err := cfg.SetUser(cmd.arg[0])
	if err != nil {
		return fmt.Errorf("Error setting the username: %v", err)
	}

	fmt.Println("The user has been set.")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	name := cmd.arg[0]
	if len(name) == 0 {
		return fmt.Errorf("Register command needs a name!")
	}

	//Check if user with given name exists
	if _, err := s.db.GetUser(context.Background(), name); err == nil {
		fmt.Println("User with same name already exist!")
		os.Exit(1)
	}

	//create new empty context for new user
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:	uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:	name,
	})
	if err != nil {
		return fmt.Errorf("Error creating new user: %v", err)
	}

	//Set current user with the given name
	cfg.SetUser(name)
	fmt.Println("User created successfully!")
	fmt.Println(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUser(context.Background())
	if err != nil {
		return fmt.Errorf("Error resetting users database!")
	}

	fmt.Println("Users database reset successful!")
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


