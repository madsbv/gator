package rss

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/madsbv/gator/internal/state"
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

func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{
		CheckRedirect: nil,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func ScrapeFeeds(s *state.State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return errors.New(fmt.Sprintf("Error getting next feed to fetch: %s\n", err))
	}

	feed_content, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return errors.New(fmt.Sprintf("Error fetching feed %s: %s\n", feed.Name.String, err))
	}

	s.Db.MarkFeedFetched(context.Background(), feed.ID)

	for _, item := range feed_content.Channel.Item {
		fmt.Println(item.Title)
	}

	return nil
}
