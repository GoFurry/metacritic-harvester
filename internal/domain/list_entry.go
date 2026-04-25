package domain

import "time"

type ListEntry struct {
	CrawlRunID string
	WorkHref   string
	Category   Category
	Metric     Metric
	Page       int
	Rank       int
	Metascore  string
	UserScore  string
	FilterKey  string
	CrawledAt  time.Time
}
