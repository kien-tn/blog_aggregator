package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kien-tn/blog_aggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("a feed url is required")
	}
	// Fetch the feed
	rssFeed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	// Insert the feed
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    rssFeed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Feed %v successfully followed by user %v\n", rssFeed.Name, user.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("error fetching follows: %w", err)
	}
	for _, follow := range follows {
		fmt.Fprintf(os.Stdout, "Feed Name: %v\n", follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("a feed url is required")
	}
	err := s.db.DropFeedFollowsForUrlCurrentUser(context.Background(), database.DropFeedFollowsForUrlCurrentUserParams{
		Url:  cmd.arguments[0],
		Name: s.config.CurrentUserName,
	})
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Feed %v successfully unfollowed for current user %v\n", cmd.arguments[0], s.config.CurrentUserName)
	return nil
}
