package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joliverstrom-cmd/gator_boot/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {

	if len(cmd.args) < 2 {
		return fmt.Errorf("Not enough arguments, usage: addfeed <name> <URL>")
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Couldn't fetch current user from DB: %w", err)
	}

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't add feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create feed follow: %w", err)
	}

	printFeed(feed)
	return nil

}

func handlerFeeds(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't read feeds from DB: %w", err)
	}

	for _, feed := range feeds {
		feedUser, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Error getting the user by User ID: %w", err)
		}
		fmt.Printf("Name: %v - URL: %v - Submitted by: %v\n", feed.Name, feed.Url, feedUser.Name)
	}

	return nil

}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}
