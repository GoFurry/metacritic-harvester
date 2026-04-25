package domain

import (
	"fmt"
	"strings"
)

type Metric string

const (
	MetricMetascore Metric = "metascore"
	MetricUserScore Metric = "userscore"
	MetricNewest    Metric = "newest"
)

func ParseMetric(raw string) (Metric, error) {
	switch Metric(strings.TrimSpace(strings.ToLower(raw))) {
	case MetricMetascore:
		return MetricMetascore, nil
	case MetricUserScore:
		return MetricUserScore, nil
	case MetricNewest:
		return MetricNewest, nil
	default:
		return "", fmt.Errorf("invalid metric %q: must be one of metascore|userscore|newest", raw)
	}
}
