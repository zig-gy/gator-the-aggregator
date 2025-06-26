package main

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/google/uuid"
	"github.com/zig-gy/gator-the-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no username passed for login")
	}

	username := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("Username %s set successfuly.\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no username passed for register")
	}

	username := cmd.arguments[0]
	createdUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: username,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error writing user to config file, try logging in: %v", err)
	}

	fmt.Println("User registered succesfully")
	fmt.Println(createdUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetUser(context.Background()); err != nil {
		return fmt.Errorf("error resetting user table: %v", err)
	}
	fmt.Println("Users deleted succesfully")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %v", err)
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if s.cfg.CurrentUsername == user.Name {
			fmt.Print(" (current)")
		}
		fmt.Println("")
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error buscando feed: %v", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Description = html.UnescapeString(item.Description)
		item.Title = html.UnescapeString(item.Title)
	}

	fmt.Println(*feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("not enough arguments passed, needs name and url")
	}

	name := cmd.arguments[0]
	url := cmd.arguments[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("error getting user from database: %v", err)
	}

	feed , err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to database: %v", err)
	}

	fmt.Println(feed)

	feedFollow, err := createFollow(s, url)
	if err != nil {
		return fmt.Errorf("error creating follow record: %v", err)
	}
	fmt.Printf("User %s followed the %s feed!\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds from the database: %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf(" - %s@%s (%s)\n", feed.FeedName, feed.Url, feed.UserName)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("not enough arguments passed, need an url")
	}

	url := cmd.arguments[0]
	feedFollow, err := createFollow(s, url)
	if err != nil {
		return fmt.Errorf("error creating follow: %v", err)
	}

	fmt.Printf("User %s followed the %s feed!\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("error getting user by name: %v", err)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follows for user: %v", err)
	}

	fmt.Printf("Feeds followed by %s\n", feedFollows[0].UserName)
	for _, row := range feedFollows {
		fmt.Printf(" - %s\n", row.FeedName)
	}
	return nil
}
