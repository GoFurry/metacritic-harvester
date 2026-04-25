package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type finderMappingKind string

const (
	finderMappingKindGenre    finderMappingKind = "genre"
	finderMappingKindPlatform finderMappingKind = "platform"
	finderMappingKindNetwork  finderMappingKind = "network"
)

type finderMappingError struct {
	Kind  finderMappingKind
	Value string
}

func (e *finderMappingError) Error() string {
	return fmt.Sprintf("finder api does not support unresolved %s %q", e.Kind, e.Value)
}

func isFinderMappingError(err error) bool {
	var target *finderMappingError
	return errors.As(err, &target)
}

func IsFinderMappingError(err error) bool {
	return isFinderMappingError(err)
}

func mapFinderGenres(values []string) ([]string, error) {
	result := make([]string, 0, len(values))
	for _, raw := range values {
		value := strings.TrimSpace(raw)
		if value == "" {
			continue
		}
		normalized := normalizeFinderMappingKey(value)
		if normalized == "" {
			continue
		}
		if mapped, ok := knownFinderGenres[normalized]; ok {
			result = append(result, mapped)
			continue
		}
		result = append(result, titleizeFinderWords(normalized))
	}
	return result, nil
}

func mapFinderPlatformIDs(values []string) ([]string, error) {
	return mapFinderIDs(values, knownFinderPlatformIDs, finderMappingKindPlatform)
}

func mapFinderNetworkIDs(values []string) ([]string, error) {
	return mapFinderIDs(values, knownFinderNetworkIDs, finderMappingKindNetwork)
}

func mapFinderIDs(values []string, known map[string]string, kind finderMappingKind) ([]string, error) {
	result := make([]string, 0, len(values))
	for _, raw := range values {
		value := strings.TrimSpace(raw)
		if value == "" {
			continue
		}
		normalized := normalizeFinderMappingKey(value)
		if mapped, ok := known[normalized]; ok {
			result = append(result, mapped)
			continue
		}
		if _, err := strconv.Atoi(value); err == nil {
			result = append(result, value)
			continue
		}
		return nil, &finderMappingError{Kind: kind, Value: raw}
	}
	return result, nil
}

func normalizeFinderMappingKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}
	replacer := strings.NewReplacer(
		"_", " ",
		"-", " ",
		"/", " ",
		"|", " ",
		"&", " and ",
	)
	value = replacer.Replace(value)
	return strings.Join(strings.Fields(value), " ")
}

func titleizeFinderWords(value string) string {
	parts := strings.Fields(strings.TrimSpace(value))
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}

var knownFinderGenres = map[string]string{
	"action":       "Action",
	"adventure":    "Adventure",
	"comedy":       "Comedy",
	"drama":        "Drama",
	"fantasy":      "Fantasy",
	"history":      "History",
	"horror":       "Horror",
	"rpg":          "RPG",
	"role playing": "RPG",
	"shooter":      "Shooter",
	"sports":       "Sports",
	"western rpg":  "Western RPG",
}

var knownFinderPlatformIDs = map[string]string{
	"pc":                  "1500000019",
	"playstation 4":       "1500000120",
	"ps4":                 "1500000120",
	"xbox one":            "1500000121",
	"switch":              "1500000122",
	"playstation 5":       "1500000128",
	"ps5":                 "1500000128",
	"xbox series x":       "1500000129",
	"xbox series xs":      "1500000129",
	"xbox series x s":     "1500000129",
	"xbox series x and s": "1500000129",
}

var knownFinderNetworkIDs = map[string]string{
	"netflix": "1943",
}
