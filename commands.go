package main

import (
	"fmt"
	"github.com/paul39-33/gator/internal/config"
	"github.com/paul39-33/gator/internal/database"
	"github.com/google/uuid"
	"time"
	"os"
	"context"
	"database/sql"
	"strconv"
	"github.com/lib/pq"
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
	time_between_reqs, err := time.ParseDuration(cmd.arg[0])
	if err != nil {
		return fmt.Errorf("Agg function needs valid time duration e.g. '10s'")
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	//set a ticker based on time_between_reqs duration
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error during scraping: %v\n", err)
		}
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error getting current user: %v", err)
		}
		return handler(s, cmd, user)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arg) < 2 {
		return fmt.Errorf("addfeed command requires a name and url!")
	}

	name := cmd.arg[0]
	url := cmd.arg[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("Error creating feed: %v", err)
	}

	//add feed and user to the feed follow table
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: feed.UserID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}

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
		fmt.Printf("Creator: %v\n", userName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	//check for the url argument
	if len(cmd.arg) < 1 {
		return fmt.Errorf("follow command needs a url!")
	}

	//get feed by URL
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.arg[0])
	if err != nil {
		return fmt.Errorf("Error getting feed from URL: %v", err)
	}


	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID:		user.ID,
		FeedID:		feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error performing feed follow: %v", err)
	}

	fmt.Printf("%v now follows %v\n", feed_follow.UserName, feed_follow.FeedName)
	return nil
	
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error getting feed follows for user: %v", err)
	}

	for _, feed_follow := range feed_follows{
		fmt.Println(feed_follow.FeedName)
	}

	return nil
}

func handlerUnfollow( s *state, cmd command, user database.User) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Unfollow command need feed url!")
	}

	err := s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		Url:	cmd.arg[0],
	})
	if err != nil {
		return fmt.Errorf("Error unfollowing feed: %v", err)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	//check if the user put an optional limit
	var limit int32
	if len(cmd.arg) == 0 {
		limit = 2
	} else {
		// Attempt to parse the string argument into an integer
		parsedLimit, err := strconv.Atoi(cmd.arg[0])
		if err != nil {
			// Handle the error if the user enters something that's not a number
			return fmt.Errorf("Invalid limit argument: %v. Please provide a number.", err)
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID:		user.ID,
		Limit:	limit,
	})
	if err != nil {
		return fmt.Errorf("Error performing browse: %v", err)
	}

	for _, post := range posts{
		fmt.Println(post.Title)
		fmt.Println(post.Url)
		fmt.Println(post.Description.String)
		fmt.Println(post.PublishedAt.Time)
		fmt.Printf("\n")
	}

	return nil

}

func scrapeFeeds(s *state) error {
	//get the next feed to fetch
	nextFeedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting next feed to fetch: %v", err)
	}

	//mark the feed as fetched
	err = s.db.MarkFeedFetched(context.Background(), nextFeedToFetch.ID)
	if err != nil {
		return fmt.Errorf("Error marking feed as fetched: %v", err)
	}

	//fetch feed
	fetchedFeed, err := fetchFeed(context.Background(), nextFeedToFetch.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %v", err)
	}
	
	//loop over fetched feed Channel.Item
	for _, item := range fetchedFeed.Channel.Item{
		//Check if the Description is empty or not and store it to a sql.NullString type
		var description sql.NullString
		if item.Description != "" {
			description = sql.NullString{
				String: item.Description,
				Valid: true,
			} 
		} else {
			description = sql.NullString{
				Valid: false,
			}
		}

		//Change PublishedAt format from string to time.Time
		parsedPubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Printf("Error parsing time format: %v", err)
			continue
		}
		publishedAt := sql.NullTime{
			Time: parsedPubDate,
			Valid: true,
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			Title:			item.Title,
			Url:			item.Link,
			Description:	description,
			PublishedAt:	publishedAt,
			FeedID:			nextFeedToFetch.ID,
		})
		if err != nil {
			// Check if the error is a *pq.Error and if its Code is "23505"
            if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
                // This is a unique constraint violation (duplicate URL)
                // The assignment says to ignore this error
                // You could optionally log it for debugging, but don't return an error
                fmt.Printf("Info: Post with URL '%s' already exists, ignoring.\n", item.Link)
            } else {
                // This is a different kind of error, so log it
                fmt.Printf("Error creating new post: %v for URL '%s'\n", err, item.Link)
            }
		}
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


