package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joliverstrom-cmd/gator_boot/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.args) < 1 {
		return fmt.Errorf("No URL supplied, usage: follow <URL>")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't follow the feed: %w", err)
	}

	fmt.Printf("User %v just subscribed to feed %v", feedFollowRow.UserName, feedFollowRow.FeedName)
	return nil
}

func handlerFollows(s *state, cmd command, user database.User) error {

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Couldn't fetch followed feeds: %w", err)
	}

	for _, feedFollow := range followedFeeds {
		fmt.Printf("* %v\n", feedFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Couldn't fetch feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't delete feed: %w", err)
	}

	return nil
}
