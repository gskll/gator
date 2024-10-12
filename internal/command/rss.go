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

func handlerAgg(s *state.State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %w", err)
	}
	fmt.Printf("%+v\n", feed)
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

	printFeed(feed)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:		%v\n", feed.ID)
	fmt.Printf("* User ID:	%v\n", feed.UserID)
	fmt.Printf("* Name:		%v\n", feed.Name)
	fmt.Printf("* URL:		%v\n", feed.Url)
}
