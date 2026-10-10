package config

import (
	"os"
	"time"
)

type Config struct {
	MTAFeedL        string
	MTAFeedAlerts   string
	MTAFeedG        string
	MTAFeedJZ       string
	MTAFeedACE      string
	MTAFeedBDFM     string
	MTAFeed123456   string
	MTAFeedNQRW     string
	StopBedfordN    string // L10N: Bedford Av Manhattan-bound (8 Av)
	StopBedfordS    string // L10S: Bedford Av Canarsie-bound (Rockaway Pkwy)
	PollInterval    time.Duration
	HTTPTimeout     time.Duration
	DegradedTimeout time.Duration
	TargetWidth     int
	TargetHeight    int
}

func Load() *Config {
	pollSecs := 20 * time.Second
	if envPoll := os.Getenv("SUBWAY_CRT_POLL_SECS"); envPoll != "" {
		if d, err := time.ParseDuration(envPoll + "s"); err == nil {
			pollSecs = d
		}
	}

	getEnv := func(key, fallback string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		return fallback
	}

	return &Config{
		MTAFeedL:        getEnv("MTA_FEED_L", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-l"),
		MTAFeedAlerts:   getEnv("MTA_FEED_ALERTS", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/camsys%2Fsubway-alerts"),
		MTAFeedG:        getEnv("MTA_FEED_G", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-g"),
		MTAFeedJZ:       getEnv("MTA_FEED_JZ", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-jz"),
		MTAFeedACE:      getEnv("MTA_FEED_ACE", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace"),
		MTAFeedBDFM:     getEnv("MTA_FEED_BDFM", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-bdfm"),
		MTAFeed123456:   getEnv("MTA_FEED_123456", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs"),
		MTAFeedNQRW:     getEnv("MTA_FEED_NQRW", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-nqrw"),
		StopBedfordN:    "L10N",
		StopBedfordS:    "L10S",
		PollInterval:    pollSecs,
		HTTPTimeout:     8 * time.Second,
		DegradedTimeout: 45 * time.Second,
		TargetWidth:     240,
		TargetHeight:    67,
	}
}
