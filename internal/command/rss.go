package command

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gskll/gator/internal/database"
	"github.com/gskll/gator/internal/rss"
	"github.com/gskll/gator/internal/state"
)

var url = "https://www.wagslane.dev/index.xml"

func handlerFollow(s *state.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s feed_url", cmd.Name)
	}

	feedUrl := cmd.Args[0]
	feed, err := s.Db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Error getting feed: %w", err)
	}
	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error getting user: %w", err)
	}

	now := time.Now()
	follow, err := s.Db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{ID: uuid.New(), FeedID: feed.ID, UserID: user.ID, CreatedAt: now, UpdatedAt: now},
	)
	if err != nil {
		return fmt.Errorf("Error creating follow: %w", err)
	}

	fmt.Printf("User '%s' is following feed '%s'\n", follow.UserName, follow.FeedName)
	return nil
}

func handlerAgg(s *state.State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %w", err)
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerFeeds(s *state.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting feeds: %w", err)
	}

	for _, feed := range feeds {
		printFeed(feed)
	}
	return nil
}

func handlerAddFeed(s *state.State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s name feed_url", cmd.Name)
	}

	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error getting user: %w", err)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	now := time.Now()

	feed, err := s.Db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{ID: uuid.New(), UserID: user.ID, Name: name, Url: url, UpdatedAt: now, CreatedAt: now},
	)
	if err != nil {
		return fmt.Errorf("Error creating feed: %w", err)
	}

	printCreatedFeed(feed)
	return nil
}

func printCreatedFeed(feed database.Feed) {
	fmt.Printf("* ID:		%v\n", feed.ID)
	fmt.Printf("* User ID:	%v\n", feed.UserID)
	fmt.Printf("* Name:		%v\n", feed.Name)
	fmt.Printf("* URL:		%v\n", feed.Url)
}

func printFeed(feed database.GetFeedsRow) {
	fmt.Printf("* ID:		%v\n", feed.ID)
	fmt.Printf("* User:		%v\n", feed.UserName)
	fmt.Printf("* Name:		%v\n", feed.Name)
	fmt.Printf("* URL:		%v\n", feed.Url)
	fmt.Println()
}
