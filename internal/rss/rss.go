package rss

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/gskll/gator/internal/database"
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

	for _, item := range rss.Channel.Item {
		if item.Title == "" || item.Link == "" {
			fmt.Printf("%s: missing title or link\n", nextFeed.Name)
			continue
		}

		pubDate, err := parsePubDate(item.PubDate)
		if err != nil {
			fmt.Printf("%s: error parsing pub date: %w\n", nextFeed.Name, err)
		}

		desc := sql.NullString{String: item.Description, Valid: item.Description != ""}

		err = s.Db.CreatePost(
			ctx,
			database.CreatePostParams{
				ID:          uuid.New(),
				FeedID:      nextFeed.ID,
				Title:       item.Title,
				Url:         item.Link,
				PublishedAt: pubDate,
				Description: desc,
			},
		)
		if err != nil {
			if !isUniqueConstraintValidation(err) {
				fmt.Printf("%s: error saving in db: %w\n", nextFeed.Name, err)
			}
		}
	}

	fmt.Printf("* Fetched: %s\n", nextFeed.Name)

	return nil
}

func isUniqueConstraintValidation(err error) bool {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return pqErr.Code == "23505"
}

func parsePubDate(dateStr string) (sql.NullTime, error) {
	if dateStr == "" {
		return sql.NullTime{}, nil
	}

	dateStr = strings.TrimSpace(dateStr)
	formats := []string{
		time.RFC1123Z,                    // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,                     // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC822Z,                     // "02 Jan 06 15:04 -0700"
		time.RFC822,                      // "02 Jan 06 15:04 MST"
		"2006-01-02T15:04:05Z07:00",      // ISO8601 with timezone
		"2006-01-02T15:04:05Z",           // ISO8601 in UTC
		"2006-01-02 15:04:05 -0700",      // Custom format sometimes used
		"2006-01-02 15:04:05",            // Custom format without timezone
		"Mon, 2 Jan 2006 15:04:05 -0700", // RFC1123Z with one-digit day
		"2 Jan 2006 15:04:05 -0700",      // RFC1123Z without weekday
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return sql.NullTime{Time: t, Valid: true}, nil
		}
	}

	return sql.NullTime{}, fmt.Errorf("unable to parse RSS date: %s", dateStr)
}
