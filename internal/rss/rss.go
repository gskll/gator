package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/gskll/gator/internal/state"
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
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	req.Header.Add("User-Agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling xml: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Title)

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func ScrapeFeeds(ctx context.Context, s *state.State) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Error getting next feed to fetch: %w", err)
	}

	err = s.Db.MarkFeedFectched(ctx, nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Error marking feed fetched '%s': %w", nextFeed.Name, err)
	}

	rss, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed '%s': %w", nextFeed.Name, err)
	}

	fmt.Printf("Fetched: %s\n", nextFeed.Name)
	for _, item := range rss.Channel.Item {
		fmt.Printf("* %s\n", item.Title)
	}
	return nil
}
