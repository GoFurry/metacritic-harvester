package domain

import (
	"fmt"
	"sort"
	"strings"
)

type Filter struct {
	ReleaseYearMin *int
	ReleaseYearMax *int
	Platforms      []string
	Networks       []string
	Genres         []string
	ReleaseTypes   []string
}

func (f Filter) Key() string {
	parts := []string{
		fmt.Sprintf("yearMin=%s", intPointerString(f.ReleaseYearMin)),
		fmt.Sprintf("yearMax=%s", intPointerString(f.ReleaseYearMax)),
		fmt.Sprintf("platforms=%s", strings.Join(normalizeValues(f.Platforms), ",")),
		fmt.Sprintf("networks=%s", strings.Join(normalizeValues(f.Networks), ",")),
		fmt.Sprintf("genres=%s", strings.Join(normalizeValues(f.Genres), ",")),
		fmt.Sprintf("releaseTypes=%s", strings.Join(normalizeValues(f.ReleaseTypes), ",")),
	}
	return strings.Join(parts, "|")
}

func normalizeValues(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func intPointerString(v *int) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%d", *v)
}
