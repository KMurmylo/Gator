package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/database"
	"html"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type command struct {
	name      string
	arguments []string
}
type commands struct {
	commandList map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.commandList[cmd.name]
	if !ok {
		return fmt.Errorf("unrecognized command")
	}
	return function(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("needs a login name")
	}
	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err == sql.ErrNoRows {
		return fmt.Errorf("no user found with that name")
	} else if err != nil {
		return fmt.Errorf("unexpected database error: %w", err)
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}
func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("needs a name")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return fmt.Errorf("user already in database")
			}
			return fmt.Errorf("postgres error code: %s, message: %s", pgErr.Code, pgErr.Message)
		} else {
			// For any other generic error:
			return fmt.Errorf("an unexpected error occurred: %e", err)
		}
	}
	s.config.SetUser(cmd.arguments[0])
	fmt.Printf("User %s \ncreated at %s \nupdated at %s \nUUd is %s\n", user.Name, user.CreatedAt, user.UpdatedAt, user.ID)
	return nil
}
func handlerResetUsers(s *state, cmd command) error {
	err := s.db.ResetUser(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Users cleaned")
	return nil
}
func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.config.UserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil

}
func handlerAgg(s *state, cmd command) error {
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i, value := range rss.Channel.Item {
		rss.Channel.Item[i].Description = html.UnescapeString(value.Description)
		rss.Channel.Item[i].Title = html.UnescapeString(value.Title)
	}
	fmt.Println(rss)
	return nil
}
func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("not enough arguments, needs name and url")
	}
	user, err := s.db.GetUser(context.Background(), s.config.UserName)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no user found with that name")
	} else if err != nil {
		return fmt.Errorf("unexpected database error: %w", err)
	}

	feed, err := s.db.InsertFeed(context.Background(), database.InsertFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to insert feed: %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to update following: %w", err)
	}
	return nil
}
func handlerListFeeds(s *state, cmd command) error {

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch feeds: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds available.")
		return nil
	}
	fmt.Println("-----------------------------")
	for _, value := range feeds {
		fmt.Printf("Name: %s\n", value.Feedname)
		fmt.Printf("URL: %s\n", value.Url)
		fmt.Printf("User: %s\n", value.Username)
		fmt.Println("-----------------------------")

	}
	return nil

}
func handlerFollow(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("No url provided")
	}

	feed, err := s.db.GetFeedURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	user, err := s.db.GetUser(context.Background(), s.config.UserName)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch create a follow: %w", err)
	}
	fmt.Printf("Created follow :\n")
	fmt.Printf("Url: %s\n", cmd.arguments[0])
	fmt.Printf("User: %s\n", user.Name)
	fmt.Printf("Feed: %s\n", feed.Name)

	return nil
}
func handlerFollowing(s *state, cmd command) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), s.config.UserName)
	if err != nil {
		return fmt.Errorf("failed to fetch feeds: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Printf("%s isn't following any feeds", s.config.UserName)
	}
	fmt.Printf("%s is following:\n", s.config.UserName)
	for _, value := range feeds {
		fmt.Printf("* %s\n", value)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.UserName)
		if err == sql.ErrNoRows {
			return fmt.Errorf("no user found with that name")
		} else if err != nil {
			return fmt.Errorf("unexpected database error: %w", err)
		}
		return handler(s, cmd, user)
	}
}
