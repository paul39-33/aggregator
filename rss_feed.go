package main

import (
	"context"
	"net/http"
	"encoding/xml"
	"io"
	"fmt"
	"html"
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
	newRSSFeed := RSSFeed{}

	//create new HTTP GET req with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &newRSSFeed, fmt.Errorf("Error creating request: %v", err)
	}

	//add header
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return &newRSSFeed, fmt.Errorf("Error performing request: %v", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &newRSSFeed, fmt.Errorf("Error performing read: %v", err)
	}

	if err = xml.Unmarshal(data, &newRSSFeed); err != nil {
		return &newRSSFeed, fmt.Errorf("Error performing unmarshal: %v", err)
	}

	for i, item := range newRSSFeed.Channel.Item{
		newRSSFeed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		newRSSFeed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	newRSSFeed.Channel.Title = html.UnescapeString(newRSSFeed.Channel.Title)
	newRSSFeed.Channel.Description = html.UnescapeString(newRSSFeed.Channel.Description)

	return &newRSSFeed, nil
}