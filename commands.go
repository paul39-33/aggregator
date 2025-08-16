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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Login command needs a name!")
	}

	//Check if user with given name exists
	if _, err := s.db.GetUser(context.Background(), cmd.arg[0]); err != nil {
		fmt.Println("User with such name doesn't exist!")
		os.Exit(1)
	}

	err := s.cfg.SetUser(cmd.arg[0])
	if err != nil {
		return fmt.Errorf("Error setting the username: %v", err)
	}

	fmt.Println("The user has been set.")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Register command needs a name!")
	}

	name := cmd.arg[0]

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
	s.cfg.SetUser(name)
	fmt.Println("User created successfully!")
	fmt.Println(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUser(context.Background())
	if err != nil {
		return fmt.Errorf("Error resetting users database: %v", err)
	}

	fmt.Println("Users database reset successful!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	//get current user
	currentUser := s.cfg.CurrentUserName

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting list of users: %v", err)
	}

	for _, user := range users{
		if user.Name == currentUser {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	rssfeed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error running agg: %v", err)
	}

	fmt.Println(rssfeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arg) < 2 {
		return fmt.Errorf("addfeed command requires a name and url!")
	}

	name := cmd.arg[0]
	url := cmd.arg[1]
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error getting current user: %v", err)
	}
	CurrentUserID := currentUser.ID

	fmt.Println(currentUser)

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: CurrentUserID,
	})

	if err != nil {
		return fmt.Errorf("Error creating feed: %v", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feedsrow, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting feeds: %v", err)
	}

	for _, feed := range feedsrow{
		//get user's name from id
		userName, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Error getting user's name: %v", err)
		}

		fmt.Printf("Name: %v\n", feed.Name)
		fmt.Printf("Url: %v\n", feed.Url)
		fmt.Printf("User name: %v\n", userName)
	}

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


