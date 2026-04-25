package domain

import "time"

type WorkDetail struct {
	WorkHref             string
	Category             Category
	Title                string
	Summary              string
	ReleaseDate          string
	Metascore            string
	MetascoreSentiment   string
	MetascoreReviewCount int
	UserScore            string
	UserScoreSentiment   string
	UserScoreCount       int
	Rating               string
	Duration             string
	Tagline              string
	Details              WorkDetailExtras
	LastFetchedAt        time.Time
}

type WorkDetailExtras struct {
	CurrentPlatform     string          `json:"current_platform,omitempty"`
	ESRBRating          string          `json:"esrb_rating,omitempty"`
	ESRBDescription     string          `json:"esrb_description,omitempty"`
	Platforms           []string        `json:"platforms,omitempty"`
	Developers          []string        `json:"developers,omitempty"`
	Publishers          []string        `json:"publishers,omitempty"`
	Genres              []string        `json:"genres,omitempty"`
	PlatformScores      []PlatformScore `json:"platform_scores,omitempty"`
	Directors           []string        `json:"directors,omitempty"`
	Writers             []string        `json:"writers,omitempty"`
	ProductionCompanies []string        `json:"production_companies,omitempty"`
	Awards              []AwardSummary  `json:"awards,omitempty"`
	Seasons             []SeasonSummary `json:"seasons,omitempty"`
	NumberOfSeasons     string          `json:"number_of_seasons,omitempty"`
	WhereToBuy          []BuyOption     `json:"where_to_buy,omitempty"`
	WhereToWatch        []WatchGroup    `json:"where_to_watch,omitempty"`
}

type PlatformScore struct {
	Platform          string `json:"platform,omitempty"`
	Href              string `json:"href,omitempty"`
	Metascore         string `json:"metascore,omitempty"`
	CriticReviewCount int    `json:"critic_review_count,omitempty"`
}

type AwardSummary struct {
	Event   string `json:"event,omitempty"`
	Details string `json:"details,omitempty"`
}

type SeasonSummary struct {
	Label     string `json:"label,omitempty"`
	Episodes  string `json:"episodes,omitempty"`
	Year      string `json:"year,omitempty"`
	Href      string `json:"href,omitempty"`
	Metascore string `json:"metascore,omitempty"`
}

type BuyOption struct {
	GroupName          string   `json:"group_name,omitempty"`
	Store              string   `json:"store,omitempty"`
	LinkURL            string   `json:"link_url,omitempty"`
	Price              *float64 `json:"price,omitempty"`
	OriginalPrice      *float64 `json:"original_price,omitempty"`
	DiscountedPrice    *float64 `json:"discounted_price,omitempty"`
	DiscountPercentage *float64 `json:"discount_percentage,omitempty"`
	ImageURL           string   `json:"image_url,omitempty"`
	PurchaseType       string   `json:"purchase_type,omitempty"`
	LowestPrice        *float64 `json:"lowest_price,omitempty"`
}

type WatchGroup struct {
	GroupName          string        `json:"group_name,omitempty"`
	ProviderName       string        `json:"provider_name,omitempty"`
	ProviderID         string        `json:"provider_id,omitempty"`
	ProviderIcon       string        `json:"provider_icon,omitempty"`
	LinkURL            string        `json:"link_url,omitempty"`
	Monetization       string        `json:"monetization,omitempty"`
	OfferType          string        `json:"offer_type,omitempty"`
	QualityType        string        `json:"quality_type,omitempty"`
	OptionCurrency     string        `json:"option_currency,omitempty"`
	OptionCurrencyCode string        `json:"option_currency_code,omitempty"`
	NumberOfSeasons    int           `json:"number_of_seasons,omitempty"`
	Options            []WatchOption `json:"options,omitempty"`
}

type WatchOption struct {
	OfferType          string   `json:"offer_type,omitempty"`
	QualityType        string   `json:"quality_type,omitempty"`
	Monetization       string   `json:"monetization,omitempty"`
	LinkURL            string   `json:"link_url,omitempty"`
	OptionCurrency     string   `json:"option_currency,omitempty"`
	OptionCurrencyCode string   `json:"option_currency_code,omitempty"`
	OptionPrice        *float64 `json:"option_price,omitempty"`
}
