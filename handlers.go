package main

import (
	"context"
	"fmt"
	"strconv"
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
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("not enough arguments passed, needs to specify a time between requests")
	}

	command := cmd.arguments[0]
	timeBetweenReqs, err := time.ParseDuration(command)
	if err != nil {
		return fmt.Errorf("error parsing durationg: %v", err)
	}

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("not enough arguments passed, needs name and url")
	}

	name := cmd.arguments[0]
	url := cmd.arguments[1]
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

	feedFollow, err := createFollow(s, url, user)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("not enough arguments passed, need an url")
	}

	url := cmd.arguments[0]
	feedFollow, err := createFollow(s, url, user)
	if err != nil {
		return fmt.Errorf("error creating follow: %v", err)
	}

	fmt.Printf("User %s followed the %s feed!\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follows for user: %v", err)
	}

	fmt.Printf("Feeds followed by %s\n", user.Name)
	for _, row := range feedFollows {
		fmt.Printf(" - %s\n", row.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("not enough arguments passed, needs an url")
	}

	url := cmd.arguments[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}

	if err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("error deleting feed follow: %v", err)
	}

	fmt.Printf("Feed %s unfollowed by user %s\n", feed.Name, user.Name)
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	if len(cmd.arguments) < 1 {
		limit = 2
	} else {
		var err error
		limit, err = strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return fmt.Errorf("can't parse limit passed, pass an integer: %v", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error getting posts for user: %v", err)
	}

	for i, post := range posts {
		fmt.Printf("Post %d: %s\n", i+1, post.Title)
		fmt.Println("-------------------Link-------------------")
		fmt.Println(post.Url)
		fmt.Printf("------------%v------------\n\n", post.PublishedAt)
	}
	return nil
}
