package rest

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Host string
	Path string

	Username string
	Passowrd string

	BearerToken string

	UserAgent string

	Transport     http.RoundTripper
	WrapTransport WrapperFunc
	Dial          func(ctx context.Context, network, address string) (net.Conn, error)
	Proxy         func(*http.Request) (*url.URL, error)

	QPS         float32
	Burst       int
	RateLimiter RateLimiter
	Timeout     time.Duration

	AuthConfig   AuthConfig
	AuthProvider AuthProvider
	ContentConfig
	// TODO:// Metric and Tracing
}

type ContentConfig struct {
	AcceptContentTypes string
	ContentType        string
	Codec              Marshaler
	Timeout            time.Duration
}

func (c *Config) TransportConfig() (*TransportConfig, error) {
	conf := &TransportConfig{
		Username:      c.Username,
		Passowrd:      c.Passowrd,
		UserAgent:     c.UserAgent,
		Transport:     c.Transport,
		WrapTransport: c.WrapTransport,
		Dial:          c.Dial,
		Proxy:         c.Proxy,
	}
	// 自定义认证
	if c.AuthProvider != nil {

		conf.Wrap(c.AuthProvider.WrapTransport)
	}
	return conf, nil
}

func ParseUrl(config *Config) (*url.URL, error) {
	host := config.Host
	if host == "" {
		host = "localhost"
	}
	return defaultServerURL(host)
}

func defaultServerURL(host string) (*url.URL, error) {
	base := host
	hostURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	return hostURL, err
}
