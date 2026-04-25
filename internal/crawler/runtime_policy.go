package crawler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type HTTPRuntimePolicy struct {
	Timeout     time.Duration
	RateLimit   rate.Limit
	RateBurst   int
	MaxInFlight int
}

func (p HTTPRuntimePolicy) normalized() HTTPRuntimePolicy {
	if p.Timeout <= 0 {
		p.Timeout = 30 * time.Second
	}
	if p.RateLimit <= 0 {
		p.RateLimit = rate.Limit(2)
	}
	if p.RateBurst <= 0 {
		p.RateBurst = 2
	}
	if p.MaxInFlight <= 0 {
		p.MaxInFlight = 1
	}
	return p
}

func WrapTransportWithPolicy(base http.RoundTripper, policy HTTPRuntimePolicy) http.RoundTripper {
	policy = policy.normalized()
	if base == nil {
		base = http.DefaultTransport
	}

	return &protectedTransport{
		base:     base,
		limiter:  rate.NewLimiter(policy.RateLimit, policy.RateBurst),
		inflight: make(chan struct{}, policy.MaxInFlight),
	}
}

type protectedTransport struct {
	base     http.RoundTripper
	limiter  *rate.Limiter
	inflight chan struct{}
}

func (t *protectedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	ctx := req.Context()
	if err := t.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	if err := acquireInflight(ctx, t.inflight); err != nil {
		return nil, err
	}
	defer func() { <-t.inflight }()

	return t.base.RoundTrip(req)
}

func acquireInflight(ctx context.Context, inflight chan struct{}) error {
	select {
	case inflight <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
