package domain

import "strings"

func NormalizeWorkHref(raw string, baseURL string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "/") {
		return strings.TrimRight(baseURL, "/") + raw
	}
	return raw
}
