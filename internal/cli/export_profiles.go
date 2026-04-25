package cli

import (
	"sort"
	"strings"
)

const (
	exportProfileRaw     = "raw"
	exportProfileFlat    = "flat"
	exportProfileSummary = "summary"
)

func isValidExportProfile(profile string) bool {
	switch profile {
	case exportProfileRaw, exportProfileFlat, exportProfileSummary:
		return true
	default:
		return false
	}
}

func joinCSVValues(values []string) string {
	filtered := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		filtered = append(filtered, trimmed)
	}
	sort.Strings(filtered)
	return strings.Join(filtered, ", ")
}
