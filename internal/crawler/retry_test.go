package crawler

import (
	"sync"
	"testing"
)

func TestRetryTrackerConcurrentAccess(t *testing.T) {
	t.Parallel()

	tracker := NewRetryTracker(1000)
	const workers = 32
	const iterations = 100

	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				tracker.Next("https://example.com/page")
				if i%10 == 0 {
					tracker.Reset("https://example.com/page")
				}
			}
		}()
	}
	wg.Wait()
}
