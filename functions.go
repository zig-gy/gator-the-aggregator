package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zig-gy/gator-the-aggregator/internal/database"
)

func createFollow(s *state, url string, user database.User) (feedFollow database.CreateFeedFollowRow, err error) {
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		err =  fmt.Errorf("error finding feed by url: %v", err)
		return
	}

	feedFollow, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		err = fmt.Errorf("error creating follow record: %v", err)
	}
	return
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return fmt.Errorf("error finding user by name: %v", err)
		}

		if err := handler(s, cmd, user); err != nil {
			return fmt.Errorf("error executing command: %v", err)
		}
		return nil
	}
}

func scrapeFeeds(s *state) error {
	feedRecord, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %v", err)
	}

	if err := s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID: feedRecord.ID,
	}); err != nil {
		return fmt.Errorf("error marking feed %s as fetched: %v", feedRecord.Name, err)
	}

	feed, err := fetchFeed(context.Background(), feedRecord.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	fmt.Printf("Feed %s:\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		fmt.Printf(" - %s\n",item.Title)
	}
	return nil
}