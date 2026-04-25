package app

import (
	"time"

	"github.com/GoFurry/metacritic-harvester/internal/crawler"
	"golang.org/x/time/rate"
)

const (
	defaultHTTPTimeout = 30 * time.Second
	defaultRunRateRPS  = 2
	defaultRunBurst    = 2
)

func listRuntimePolicy() crawler.HTTPRuntimePolicy {
	return crawler.HTTPRuntimePolicy{
		Timeout:     defaultHTTPTimeout,
		RateLimit:   rate.Limit(defaultRunRateRPS),
		RateBurst:   defaultRunBurst,
		MaxInFlight: 1,
	}
}

func detailRuntimePolicy(concurrency int) crawler.HTTPRuntimePolicy {
	if concurrency <= 0 {
		concurrency = 1
	}
	return crawler.HTTPRuntimePolicy{
		Timeout:     defaultHTTPTimeout,
		RateLimit:   rate.Limit(defaultRunRateRPS),
		RateBurst:   defaultRunBurst,
		MaxInFlight: concurrency,
	}
}

func reviewRuntimePolicy(concurrency int) crawler.HTTPRuntimePolicy {
	if concurrency <= 0 {
		concurrency = 1
	}
	return crawler.HTTPRuntimePolicy{
		Timeout:     defaultHTTPTimeout,
		RateLimit:   rate.Limit(defaultRunRateRPS),
		RateBurst:   defaultRunBurst,
		MaxInFlight: concurrency,
	}
}
