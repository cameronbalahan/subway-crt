package mta

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(timeout time.Duration) *Client {
	apiKey := os.Getenv("MTA_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("MTA_KEY")
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        12,
				MaxIdleConnsPerHost: 6,
				IdleConnTimeout:     45 * time.Second,
				DisableCompression: false,
			},
		},
		apiKey: apiKey,
	}
}

// FetchFeed handles decompression, headers, timeouts, and unmarshaling
func (c *Client) FetchFeed(ctx context.Context, feedURL string) (*gtfs.FeedMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("request build failed: %w", err)
	}

	req.Header.Set("Accept", "application/x-protobuf, application/octet-stream, */*")
	req.Header.Set("User-Agent", "subway-crt/1.0 (Raspberry Pi/Terminal Kiosk; +https://github.com/charmbracelet/bubbletea)")

	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network connection error for %s: %w", feedURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MTA API returned status %d for %s", resp.StatusCode, feedURL)
	}

	var reader io.Reader = io.LimitReader(resp.Body, 8*1024*1024)
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(reader)
		if err == nil {
			defer gz.Close()
			reader = gz
		}
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	feed := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		return nil, fmt.Errorf("protobuf unmarshal failed: %w", err)
	}

	return feed, nil
}
