package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Couldn't create new request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Couldn't post the request: %w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Couldn't read the response body: %w", err)
	}

	myRSS := RSSFeed{}

	err = xml.Unmarshal(data, &myRSS)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Couldn't unmarshal the feed: %w", err)
	}

	myRSS.Channel.Title = html.UnescapeString(myRSS.Channel.Title)
	myRSS.Channel.Title = html.UnescapeString(myRSS.Channel.Description)

	for i := 0; i < len(myRSS.Channel.Item); i++ {
		myRSS.Channel.Item[i].Title = html.UnescapeString(myRSS.Channel.Item[i].Title)
		myRSS.Channel.Item[i].Description = html.UnescapeString(myRSS.Channel.Item[i].Description)
	}

	return &myRSS, nil

}
