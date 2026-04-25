package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
)

type Config struct {
	AllowedDomains []string
	Debug          bool
	ProxyURLs      []string
	MaxRetries     int
	Transport      *http.Transport
	ProxyRotator   *ProxyRotator
}

func NewCollector(cfg Config) (*colly.Collector, *RetryTracker, error) {
	options := []colly.CollectorOption{
		colly.AllowURLRevisit(),
	}
	if len(cfg.AllowedDomains) > 0 {
		options = append(options, colly.AllowedDomains(cfg.AllowedDomains...))
	}

	c := colly.NewCollector(options...)
	c.SetRequestTimeout(30 * time.Second)

	rotator := cfg.ProxyRotator
	if rotator == nil {
		var err error
		rotator, err = NewProxyRotator(cfg.ProxyURLs)
		if err != nil {
			return nil, nil, err
		}
	}

	transport := cfg.Transport
	if transport == nil {
		transport = NewHTTPTransport(cfg.Debug, rotator)
	}
	c.WithTransport(transport)

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       800 * time.Millisecond,
		RandomDelay: 400 * time.Millisecond,
	})
	if err != nil {
		return nil, nil, err
	}

	retryTracker := NewRetryTracker(cfg.MaxRetries)
	return c, retryTracker, nil
}

func NewHTTPTransport(debug bool, rotator *ProxyRotator) *http.Transport {
	transport := &http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
		IdleConnTimeout:       30 * time.Second,
	}
	transport.Proxy = func(req *http.Request) (*url.URL, error) {
		if rotator == nil {
			return nil, nil
		}
		proxy, err := rotator.Next()
		if err != nil {
			return nil, err
		}
		if debug && proxy != nil {
			fmt.Printf("using proxy: %s -> %s\n", proxy.String(), req.URL.String())
		}
		return proxy, nil
	}
	return transport
}

func SetDefaultRequestHeaders(r *colly.Request) {
	r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	r.Headers.Set("Referer", "https://www.metacritic.com/")
}
