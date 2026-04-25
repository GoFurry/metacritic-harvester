package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"sync/atomic"
)

type ProxyRotator struct {
	proxies []*url.URL
	index   uint64
}

func NewProxyRotator(rawURLs []string) (*ProxyRotator, error) {
	var proxies []*url.URL
	for _, raw := range rawURLs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		parsed, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy url %q: %w", raw, err)
		}
		proxies = append(proxies, parsed)
	}

	return &ProxyRotator{proxies: proxies}, nil
}

func (p *ProxyRotator) Next() (*url.URL, error) {
	if len(p.proxies) == 0 {
		return nil, nil
	}

	i := atomic.AddUint64(&p.index, 1)
	return p.proxies[(int(i)-1)%len(p.proxies)], nil
}
