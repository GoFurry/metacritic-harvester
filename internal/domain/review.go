package domain

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type ReviewType string

const (
	ReviewTypeCritic ReviewType = "critic"
	ReviewTypeUser   ReviewType = "user"
	ReviewTypeAll    ReviewType = "all"
)

type ReviewSentiment string

const (
	ReviewSentimentAll      ReviewSentiment = "all"
	ReviewSentimentPositive ReviewSentiment = "positive"
	ReviewSentimentNeutral  ReviewSentiment = "neutral"
	ReviewSentimentNegative ReviewSentiment = "negative"
)

type ReviewSort string

const (
	ReviewSortDate        ReviewSort = "date"
	ReviewSortScore       ReviewSort = "score"
	ReviewSortPublication ReviewSort = "publication"
)

func ParseReviewType(raw string) (ReviewType, error) {
	switch ReviewType(strings.TrimSpace(strings.ToLower(raw))) {
	case ReviewTypeCritic:
		return ReviewTypeCritic, nil
	case ReviewTypeUser:
		return ReviewTypeUser, nil
	case ReviewTypeAll:
		return ReviewTypeAll, nil
	default:
		return "", fmt.Errorf("invalid review type %q: must be one of critic|user|all", raw)
	}
}

func ParseReviewSentiment(raw string) (ReviewSentiment, error) {
	switch ReviewSentiment(strings.TrimSpace(strings.ToLower(raw))) {
	case "", ReviewSentimentAll:
		return ReviewSentimentAll, nil
	case ReviewSentimentPositive:
		return ReviewSentimentPositive, nil
	case ReviewSentimentNeutral:
		return ReviewSentimentNeutral, nil
	case ReviewSentimentNegative:
		return ReviewSentimentNegative, nil
	default:
		return "", fmt.Errorf("invalid review sentiment %q: must be one of all|positive|neutral|negative", raw)
	}
}

func ParseReviewSort(raw string) (ReviewSort, error) {
	switch ReviewSort(strings.TrimSpace(strings.ToLower(raw))) {
	case "":
		return "", nil
	case ReviewSortDate:
		return ReviewSortDate, nil
	case ReviewSortScore:
		return ReviewSortScore, nil
	case ReviewSortPublication:
		return ReviewSortPublication, nil
	default:
		return "", fmt.Errorf("invalid review sort %q: must be one of date|score|publication", raw)
	}
}

type ReviewTask struct {
	Category    Category
	WorkHref    string
	Limit       int
	Force       bool
	Concurrency int
	ReviewType  ReviewType
	Sentiment   ReviewSentiment
	Sort        ReviewSort
	Platform    string
	PageSize    int
	MaxPages    int
	Debug       bool
}

type ReviewScope struct {
	WorkHref    string
	Category    Category
	ReviewType  ReviewType
	PlatformKey string
}

func (s ReviewScope) Key() string {
	return fmt.Sprintf("%s|%s|%s|%s", NormalizeWorkHref(s.WorkHref, ""), s.Category, s.ReviewType, strings.TrimSpace(s.PlatformKey))
}

type ReviewRecord struct {
	ReviewKey         string
	ExternalReviewID  string
	CrawlRunID        string
	WorkHref          string
	Category          Category
	ReviewType        ReviewType
	PlatformKey       string
	ReviewURL         string
	ReviewDate        string
	Score             *float64
	Quote             string
	PublicationName   string
	PublicationSlug   string
	AuthorName        string
	AuthorSlug        string
	SeasonLabel       string
	Username          string
	UserSlug          string
	ThumbsUp          *int64
	ThumbsDown        *int64
	VersionLabel      string
	SpoilerFlag       *bool
	SourcePayloadJSON string
	CrawledAt         time.Time
}

func BuildCriticReviewKey(workHref string, category Category, platformKey string, publicationSlug string, reviewDate string, quote string) string {
	return buildReviewKey(category, ReviewTypeCritic, workHref, platformKey, publicationSlug, reviewDate, "", quote)
}

func BuildUserReviewKey(workHref string, category Category, platformKey string, externalID string, author string, reviewDate string, score *float64, quote string) string {
	if strings.TrimSpace(externalID) != "" {
		return fmt.Sprintf("%s|%s|%s|%s|%s", category, ReviewTypeUser, NormalizeWorkHref(workHref, ""), strings.TrimSpace(platformKey), strings.TrimSpace(externalID))
	}

	scoreText := ""
	if score != nil {
		scoreText = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.4f", *score), "0"), ".")
	}
	return buildReviewKey(category, ReviewTypeUser, workHref, platformKey, author, reviewDate, scoreText, quote)
}

func buildReviewKey(category Category, reviewType ReviewType, workHref string, platformKey string, actor string, reviewDate string, scoreText string, quote string) string {
	base := strings.Join([]string{
		string(category),
		string(reviewType),
		NormalizeWorkHref(workHref, ""),
		strings.TrimSpace(platformKey),
		strings.TrimSpace(strings.ToLower(actor)),
		strings.TrimSpace(reviewDate),
		strings.TrimSpace(scoreText),
		hashNormalizedQuote(quote),
	}, "|")
	return base
}

func hashNormalizedQuote(quote string) string {
	normalized := strings.Join(strings.Fields(strings.TrimSpace(strings.ToLower(quote))), " ")
	sum := sha1.Sum([]byte(normalized))
	return hex.EncodeToString(sum[:])
}
