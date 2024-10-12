package command

import (
	"context"
	"fmt"

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
