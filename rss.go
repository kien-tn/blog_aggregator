package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kien-tn/blog_aggregator/internal/database"
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
	// Fetch the feed
	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// Parse the feed
	rssFeed := &RSSFeed{}
	err = xml.Unmarshal(body, rssFeed)
	if err != nil {
		return nil, err
	}
	return rssFeed, nil

}

func handlerFetchFeed(s *state, cmd command) error {
	// do something
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("a feed URL is required")
	}
	feedURL := cmd.arguments[0]
	feed, err := fetchFeed(context.Background(), feedURL)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Title: %v\n", feed.Channel.Title)
	fmt.Fprintf(os.Stdout, "Description: %v\n", feed.Channel.Description)
	// fmt.Fprintf(os.Stdout, "Feed fetched successfully: %v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("addfeed requires 2 args: a name and a URL")
	}
	name := cmd.arguments[0]
	url := cmd.arguments[1]
	user, err := s.db.GetUserByName(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:      name,
		Url:       url,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Feed %v successfully created\n", feed)
	return nil
}
