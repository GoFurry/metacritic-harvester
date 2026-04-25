package cli

import (
	"database/sql"

	"github.com/GoFurry/metacritic-harvester/internal/domain"
	sqlcgen "github.com/GoFurry/metacritic-harvester/internal/storage/sqlcgen"
)

type reviewView struct {
	ReviewKey         string `json:"review_key"`
	ExternalReviewID  string `json:"external_review_id,omitempty"`
	WorkHref          string `json:"work_href"`
	Category          string `json:"category"`
	ReviewType        string `json:"review_type"`
	PlatformKey       string `json:"platform_key,omitempty"`
	ReviewURL         string `json:"review_url,omitempty"`
	ReviewDate        string `json:"review_date,omitempty"`
	Score             string `json:"score,omitempty"`
	Quote             string `json:"quote,omitempty"`
	PublicationName   string `json:"publication_name,omitempty"`
	PublicationSlug   string `json:"publication_slug,omitempty"`
	AuthorName        string `json:"author_name,omitempty"`
	AuthorSlug        string `json:"author_slug,omitempty"`
	SeasonLabel       string `json:"season_label,omitempty"`
	Username          string `json:"username,omitempty"`
	UserSlug          string `json:"user_slug,omitempty"`
	ThumbsUp          string `json:"thumbs_up,omitempty"`
	ThumbsDown        string `json:"thumbs_down,omitempty"`
	VersionLabel      string `json:"version_label,omitempty"`
	SpoilerFlag       string `json:"spoiler_flag,omitempty"`
	SourcePayloadJSON string `json:"source_payload_json,omitempty"`
	SourceCrawlRunID  string `json:"source_crawl_run_id"`
	LastCrawledAt     string `json:"last_crawled_at"`
}

type reviewCompareView struct {
	ReviewKey        string `json:"review_key"`
	WorkHref         string `json:"work_href"`
	Category         string `json:"category"`
	ReviewType       string `json:"review_type"`
	PlatformKey      string `json:"platform_key"`
	FromScore        string `json:"from_score,omitempty"`
	ToScore          string `json:"to_score,omitempty"`
	ScoreDiff        string `json:"score_diff,omitempty"`
	FromQuote        string `json:"from_quote,omitempty"`
	ToQuote          string `json:"to_quote,omitempty"`
	FromThumbsUp     string `json:"from_thumbs_up,omitempty"`
	ToThumbsUp       string `json:"to_thumbs_up,omitempty"`
	FromThumbsDown   string `json:"from_thumbs_down,omitempty"`
	ToThumbsDown     string `json:"to_thumbs_down,omitempty"`
	FromVersionLabel string `json:"from_version_label,omitempty"`
	ToVersionLabel   string `json:"to_version_label,omitempty"`
	FromSpoilerFlag  string `json:"from_spoiler_flag,omitempty"`
	ToSpoilerFlag    string `json:"to_spoiler_flag,omitempty"`
	ChangeType       string `json:"change_type"`
}

func mapLatestReviews(rows []sqlcgen.LatestReview) []reviewView {
	result := make([]reviewView, 0, len(rows))
	for _, row := range rows {
		result = append(result, reviewView{
			ReviewKey:         row.ReviewKey,
			ExternalReviewID:  nullStringValue(row.ExternalReviewID),
			WorkHref:          row.WorkHref,
			Category:          row.Category,
			ReviewType:        row.ReviewType,
			PlatformKey:       row.PlatformKey,
			ReviewURL:         nullStringValue(row.ReviewUrl),
			ReviewDate:        nullStringValue(row.ReviewDate),
			Score:             nullFloat64Value(row.Score),
			Quote:             nullStringValue(row.Quote),
			PublicationName:   nullStringValue(row.PublicationName),
			PublicationSlug:   nullStringValue(row.PublicationSlug),
			AuthorName:        nullStringValue(row.AuthorName),
			AuthorSlug:        nullStringValue(row.AuthorSlug),
			SeasonLabel:       nullStringValue(row.SeasonLabel),
			Username:          nullStringValue(row.Username),
			UserSlug:          nullStringValue(row.UserSlug),
			ThumbsUp:          nullInt64Value(row.ThumbsUp),
			ThumbsDown:        nullInt64Value(row.ThumbsDown),
			VersionLabel:      nullStringValue(row.VersionLabel),
			SpoilerFlag:       nullBoolIntValue(row.SpoilerFlag),
			SourcePayloadJSON: row.SourcePayloadJson,
			SourceCrawlRunID:  row.SourceCrawlRunID,
			LastCrawledAt:     row.LastCrawledAt,
		})
	}
	return result
}

func mapReviewCompareRows(rows []sqlcgen.CompareReviewSnapshotsRow) []reviewCompareView {
	result := make([]reviewCompareView, 0, len(rows))
	for _, row := range rows {
		result = append(result, reviewCompareView{
			ReviewKey:        row.ReviewKey,
			WorkHref:         row.WorkHref,
			Category:         row.Category,
			ReviewType:       row.ReviewType,
			PlatformKey:      row.PlatformKey,
			FromScore:        nullFloat64Value(row.FromScore),
			ToScore:          nullFloat64Value(row.ToScore),
			ScoreDiff:        interfaceValueString(row.ScoreDiff),
			FromQuote:        nullStringValue(row.FromQuote),
			ToQuote:          nullStringValue(row.ToQuote),
			FromThumbsUp:     nullInt64Value(row.FromThumbsUp),
			ToThumbsUp:       nullInt64Value(row.ToThumbsUp),
			FromThumbsDown:   nullInt64Value(row.FromThumbsDown),
			ToThumbsDown:     nullInt64Value(row.ToThumbsDown),
			FromVersionLabel: nullStringValue(row.FromVersionLabel),
			ToVersionLabel:   nullStringValue(row.ToVersionLabel),
			FromSpoilerFlag:  nullBoolIntValue(row.FromSpoilerFlag),
			ToSpoilerFlag:    nullBoolIntValue(row.ToSpoilerFlag),
			ChangeType:       row.ChangeType,
		})
	}
	return result
}

func reviewScoreString(value *float64) string {
	if value == nil {
		return ""
	}
	return nullFloat64Value(sql.NullFloat64{Float64: *value, Valid: true})
}

func reviewInt64String(value *int64) string {
	if value == nil {
		return ""
	}
	return nullInt64Value(sql.NullInt64{Int64: *value, Valid: true})
}

func reviewBoolString(value *bool) string {
	if value == nil {
		return ""
	}
	return nullBoolIntValue(sql.NullInt64{Int64: func() int64 {
		if *value {
			return 1
		}
		return 0
	}(), Valid: true})
}

func reviewReviewTypeValue(value domain.ReviewType) string {
	return string(value)
}
