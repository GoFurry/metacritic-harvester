package cli

import (
	"database/sql"

	"github.com/bytedance/sonic"

	"github.com/GoFurry/metacritic-harvester/internal/domain"
	sqlcgen "github.com/GoFurry/metacritic-harvester/internal/storage/sqlcgen"
)

type detailView struct {
	WorkHref             string                  `json:"work_href"`
	Category             string                  `json:"category"`
	Title                string                  `json:"title"`
	Summary              string                  `json:"summary,omitempty"`
	ReleaseDate          string                  `json:"release_date,omitempty"`
	Metascore            string                  `json:"metascore,omitempty"`
	MetascoreSentiment   string                  `json:"metascore_sentiment,omitempty"`
	MetascoreReviewCount int                     `json:"metascore_review_count,omitempty"`
	UserScore            string                  `json:"user_score,omitempty"`
	UserScoreSentiment   string                  `json:"user_score_sentiment,omitempty"`
	UserScoreCount       int                     `json:"user_score_count,omitempty"`
	Rating               string                  `json:"rating,omitempty"`
	Duration             string                  `json:"duration,omitempty"`
	Tagline              string                  `json:"tagline,omitempty"`
	LastFetchedAt        string                  `json:"last_fetched_at"`
	Details              domain.WorkDetailExtras `json:"details"`
	RawDetailsJSON       string                  `json:"-"`
}

type detailCompareView struct {
	WorkHref           string `json:"work_href"`
	Category           string `json:"category"`
	ChangeType         string `json:"change_type"`
	FromTitle          string `json:"from_title,omitempty"`
	ToTitle            string `json:"to_title,omitempty"`
	FromReleaseDate    string `json:"from_release_date,omitempty"`
	ToReleaseDate      string `json:"to_release_date,omitempty"`
	FromMetascore      string `json:"from_metascore,omitempty"`
	ToMetascore        string `json:"to_metascore,omitempty"`
	FromUserScore      string `json:"from_user_score,omitempty"`
	ToUserScore        string `json:"to_user_score,omitempty"`
	FromRating         string `json:"from_rating,omitempty"`
	ToRating           string `json:"to_rating,omitempty"`
	FromDuration       string `json:"from_duration,omitempty"`
	ToDuration         string `json:"to_duration,omitempty"`
	FromTagline        string `json:"from_tagline,omitempty"`
	ToTagline          string `json:"to_tagline,omitempty"`
	DetailsJSONChanged bool   `json:"details_json_changed"`
	FromDetailsJSON    string `json:"from_details_json,omitempty"`
	ToDetailsJSON      string `json:"to_details_json,omitempty"`
}

func mapWorkDetails(rows []sqlcgen.WorkDetail) ([]detailView, error) {
	result := make([]detailView, 0, len(rows))
	for _, row := range rows {
		details, err := unmarshalWorkDetailExtras(row.DetailsJson)
		if err != nil {
			return nil, err
		}
		result = append(result, detailView{
			WorkHref:             row.WorkHref,
			Category:             row.Category,
			Title:                row.Title,
			Summary:              nullStringValue(row.Summary),
			ReleaseDate:          nullStringValue(row.ReleaseDate),
			Metascore:            nullStringValue(row.Metascore),
			MetascoreSentiment:   nullStringValue(row.MetascoreSentiment),
			MetascoreReviewCount: nullIntValue(row.MetascoreReviewCount),
			UserScore:            nullStringValue(row.UserScore),
			UserScoreSentiment:   nullStringValue(row.UserScoreSentiment),
			UserScoreCount:       nullIntValue(row.UserScoreCount),
			Rating:               nullStringValue(row.Rating),
			Duration:             nullStringValue(row.Duration),
			Tagline:              nullStringValue(row.Tagline),
			LastFetchedAt:        row.LastFetchedAt,
			Details:              details,
			RawDetailsJSON:       row.DetailsJson,
		})
	}
	return result, nil
}

func mapDetailCompareRows(rows []sqlcgen.CompareWorkDetailSnapshotsRow) []detailCompareView {
	result := make([]detailCompareView, 0, len(rows))
	for _, row := range rows {
		result = append(result, detailCompareView{
			WorkHref:           row.WorkHref,
			Category:           row.Category,
			ChangeType:         row.ChangeType,
			FromTitle:          interfaceValueString(row.FromTitle),
			ToTitle:            nullStringValue(row.ToTitle),
			FromReleaseDate:    nullStringValue(row.FromReleaseDate),
			ToReleaseDate:      nullStringValue(row.ToReleaseDate),
			FromMetascore:      nullStringValue(row.FromMetascore),
			ToMetascore:        nullStringValue(row.ToMetascore),
			FromUserScore:      nullStringValue(row.FromUserScore),
			ToUserScore:        nullStringValue(row.ToUserScore),
			FromRating:         nullStringValue(row.FromRating),
			ToRating:           nullStringValue(row.ToRating),
			FromDuration:       nullStringValue(row.FromDuration),
			ToDuration:         nullStringValue(row.ToDuration),
			FromTagline:        nullStringValue(row.FromTagline),
			ToTagline:          nullStringValue(row.ToTagline),
			DetailsJSONChanged: row.DetailsJsonChanged != 0,
			FromDetailsJSON:    interfaceValueString(row.FromDetailsJson),
			ToDetailsJSON:      nullStringValue(row.ToDetailsJson),
		})
	}
	return result
}

func unmarshalWorkDetailExtras(raw string) (domain.WorkDetailExtras, error) {
	var details domain.WorkDetailExtras
	if raw == "" {
		raw = "{}"
	}
	if err := sonic.UnmarshalString(raw, &details); err != nil {
		return domain.WorkDetailExtras{}, err
	}
	return details, nil
}

func nullIntValue(value sql.NullInt64) int {
	if !value.Valid {
		return 0
	}
	return int(value.Int64)
}
