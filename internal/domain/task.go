package domain

type ListTask struct {
	Category Category
	Metric   Metric
	Filter   Filter
	MaxPages int
	Debug    bool
}

type DetailTask struct {
	Category    Category
	WorkHref    string
	Limit       int
	Force       bool
	Debug       bool
	Concurrency int
}
