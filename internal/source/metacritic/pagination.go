package metacritic

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ParsePagination(e *colly.HTMLElement) int {
	maxFoundPage := 1
	e.ForEach(SelectorPaginationItem, func(_ int, el *colly.HTMLElement) {
		text := strings.TrimSpace(el.Text)
		if text == "" {
			return
		}

		n, err := strconv.Atoi(text)
		if err == nil && n > maxFoundPage {
			maxFoundPage = n
		}
	})
	return maxFoundPage
}
