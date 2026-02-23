package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joliverstrom-cmd/gator_boot/internal/database"
)

func handlerAggs(s *state, cmd command) error {

	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: agg <time between requests> time between requests being e.g. 1s, 1m, 1h...")
	}

	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Invalid duration supplied: %w", err)
	}

	fmt.Printf("Collecting feeds every %v", time_between_reqs.String())

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("Something went wrong: %w", err)
		}
	}

}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't get next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't update the feed's last fetched timestamp: %w", err)
	}

	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Couldn't fetch the feed from it's source: %w", err)
	}
	fmt.Printf("Scraping feed: %v\n", fetchedFeed.Channel.Title)
	err = scrapeRSSFeed(s, fetchedFeed, feed.ID)
	if err != nil {
		return fmt.Errorf("Couldn't scrape: %w", err)
	}
	return nil

}

func addPostToDB(s *state, post RSSItem, feedID uuid.UUID) (database.Post, error) {

	addedPost, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       post.Title,
		Url:         post.Link,
		Description: sql.NullString{String: post.Description, Valid: true},
		PublishedAt: sql.NullString{String: post.PubDate, Valid: true},
		FeedID:      feedID,
	})
	if err != nil {
		return database.Post{}, fmt.Errorf("error is: %w", err)
	}

	return addedPost, nil

}

func scrapeRSSFeed(s *state, feed *RSSFeed, feedID uuid.UUID) error {
	fmt.Printf("Scraping feed: %v\n", feed.Channel.Title)

	if len(feed.Channel.Item) == 0 {
		fmt.Println("No items in feed")
		return nil
	}

	fmt.Println("")
	fmt.Printf("Feed: %v\n", feed.Channel.Title)

	for _, item := range feed.Channel.Item {
		post, err := addPostToDB(s, item, feedID)
		if err != nil {

			if strings.Contains(err.Error(), "duplicate key value") {
				continue
			}

			return fmt.Errorf("Couldn't add post to DB: %w", err)
		}
		fmt.Printf("* Post added – Title: %v, URL: %v\n", post.Title, post.Url)
	}
	fmt.Println("")
	fmt.Println("******************************************")

	return nil
}
