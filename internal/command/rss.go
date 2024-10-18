package command

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/gskll/gator/internal/database"
	"github.com/gskll/gator/internal/rss"
	"github.com/gskll/gator/internal/state"
)

func handlerBrowse(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("Usage: %s [number_of_posts]", cmd.Name)
	}

	var err error
	numPosts := 2
	if len(cmd.Args) == 1 {
		numPosts, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("number_of_posts must be an int")
		}
		if numPosts < 1 {
			return fmt.Errorf("Usage: %s [number_of_posts]", cmd.Name)
		}
	}

	posts, err := s.Db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{UserID: user.ID, Limit: int32(numPosts)},
	)
	if err != nil {
		return fmt.Errorf("Error getting posts: %w", err)
	}

	for _, post := range posts {
		printPost(post)
		fmt.Println()
	}

	return nil
}

func handlerAgg(s *state.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s time_between_reqs", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error parsing time_between_reqs: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenReqs)

	ctx := context.Background()
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = rss.ScrapeFeeds(ctx, s)
		if err != nil {
			return fmt.Errorf("Error scraping feeds: %w", err)
		}
	}
}

func handlerFollowing(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	following, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error getting user following: %w", err)
	}

	for _, feed := range following {
		fmt.Printf("* %s\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s feed_url", cmd.Name)
	}

	feedUrl := cmd.Args[0]
	feed, err := s.Db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Error getting feed: %w", err)
	}

	err = s.Db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("Error deleting feed: %w", err)
	}

	fmt.Printf("User '%s' has unfollowed feed '%s'\n", user.Name, feed.Name)
	return nil
}

func handlerFollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s feed_url", cmd.Name)
	}

	feedUrl := cmd.Args[0]
	feed, err := s.Db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Error getting feed: %w", err)
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

func handlerAddFeed(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s name feed_url", cmd.Name)
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]
	now := time.Now()

	feed, err := s.Db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			UserID:    user.ID,
			Name:      feedName,
			Url:       feedUrl,
			UpdatedAt: now,
			CreatedAt: now,
		},
	)
	if err != nil {
		return fmt.Errorf("Error creating feed: %w", err)
	}

	_, err = s.Db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{ID: uuid.New(), UserID: user.ID, FeedID: feed.ID, CreatedAt: now, UpdatedAt: now},
	)
	if err != nil {
		return fmt.Errorf("Error following new feed: %w", err)
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

func printPost(post database.GetPostsForUserRow) {
	fmt.Printf("* ID:			%v\n", post.ID)
	fmt.Printf("* Feed:			%v\n", post.FeedName)
	fmt.Printf("* Title:		%v\n", post.Title)
	fmt.Printf("* URL:			%v\n", post.Url)
}
