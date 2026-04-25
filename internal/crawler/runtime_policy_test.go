package crawler

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestWrapTransportWithPolicyHonorsMaxInFlight(t *testing.T) {
	t.Parallel()

	var active int32
	var maxSeen int32
	base := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		current := atomic.AddInt32(&active, 1)
		defer atomic.AddInt32(&active, -1)

		for {
			prev := atomic.LoadInt32(&maxSeen)
			if current <= prev || atomic.CompareAndSwapInt32(&maxSeen, prev, current) {
				break
			}
		}

		time.Sleep(20 * time.Millisecond)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("ok")),
		}, nil
	})

	transport := WrapTransportWithPolicy(base, HTTPRuntimePolicy{
		Timeout:     30 * time.Second,
		RateLimit:   rate.Limit(100),
		RateBurst:   100,
		MaxInFlight: 2,
	})

	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
			if err != nil {
				t.Errorf("NewRequestWithContext() error = %v", err)
				return
			}
			resp, err := transport.RoundTrip(req)
			if err != nil {
				t.Errorf("RoundTrip() error = %v", err)
				return
			}
			_ = resp.Body.Close()
		}()
	}
	wg.Wait()

	if got := atomic.LoadInt32(&maxSeen); got > 2 {
		t.Fatalf("expected max 2 in-flight requests, got %d", got)
	}
}

func TestWrapTransportWithPolicyRateLimitsRequests(t *testing.T) {
	t.Parallel()

	transport := WrapTransportWithPolicy(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("ok")),
		}, nil
	}), HTTPRuntimePolicy{
		Timeout:     30 * time.Second,
		RateLimit:   rate.Limit(10),
		RateBurst:   1,
		MaxInFlight: 1,
	})

	start := time.Now()
	for i := 0; i < 2; i++ {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
		if err != nil {
			t.Fatalf("NewRequestWithContext() error = %v", err)
		}
		resp, err := transport.RoundTrip(req)
		if err != nil {
			t.Fatalf("RoundTrip() error = %v", err)
		}
		_ = resp.Body.Close()
	}

	if elapsed := time.Since(start); elapsed < 80*time.Millisecond {
		t.Fatalf("expected second request to be rate-limited, elapsed=%s", elapsed)
	}
}

func TestWrapTransportWithPolicyRespectsContextWhileWaitingForInflight(t *testing.T) {
	t.Parallel()

	block := make(chan struct{})
	base := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		<-block
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("ok")),
		}, nil
	})

	transport := WrapTransportWithPolicy(base, HTTPRuntimePolicy{
		Timeout:     30 * time.Second,
		RateLimit:   rate.Limit(100),
		RateBurst:   100,
		MaxInFlight: 1,
	})

	firstReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext() error = %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		resp, rtErr := transport.RoundTrip(firstReq)
		if rtErr == nil && resp != nil {
			_ = resp.Body.Close()
		}
	}()

	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	secondReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext() error = %v", err)
	}
	if _, err := transport.RoundTrip(secondReq); err == nil {
		t.Fatal("expected context deadline while waiting for in-flight slot")
	}

	close(block)
	<-done
}
