package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zig-gy/gator-the-aggregator/internal/database"
)

func createFollow(s *state, url string) (feedFollow database.CreateFeedFollowRow, err error) {
	feedId, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		err =  fmt.Errorf("error finding feed by url: %v", err)
		return
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		err = fmt.Errorf("error finding user by name: %v", err)
		return 
	}

	feedFollow, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feedId,
	})
	if err != nil {
		err = fmt.Errorf("error creating follow record: %v", err)
	}
	return
}