package crawler

import "sync"

type RetryTracker struct {
	mu     sync.Mutex
	counts map[string]int
	max    int
}

func NewRetryTracker(max int) *RetryTracker {
	return &RetryTracker{
		counts: make(map[string]int),
		max:    max,
	}
}

func (r *RetryTracker) Next(url string) (attempt int, shouldRetry bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counts[url]++
	attempt = r.counts[url]
	return attempt, attempt <= r.max
}

func (r *RetryTracker) Reset(url string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.counts, url)
}
