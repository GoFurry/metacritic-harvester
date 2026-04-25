package metacritic

const (
	SelectorCard             = `div[data-testid="filter-results"]`
	SelectorTitle            = `h3[data-testid="product-title"] span:last-child`
	SelectorPagination       = `nav[data-testid="navigation-pagination"]`
	SelectorPaginationItem   = `.c-navigation-pagination__page .c-navigation-pagination__item-content`
	SelectorMetascorePrimary = `[aria-label*="Metascore"] span`
	SelectorUserScorePrimary = `[aria-label*="User Score"] span, [aria-label*="User score"] span`
	SelectorFallbackScore    = `.c-siteReviewScore span`
)
