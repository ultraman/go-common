package rest

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"
)

type TransportConfig struct {
	Username string
	Passowrd string

	UserAgent string

	Transport     http.RoundTripper
	WrapTransport WrapperFunc
	Dial          func(ctx context.Context, network, address string) (net.Conn, error)
	Proxy         func(*http.Request) (*url.URL, error)
}

func (c *TransportConfig) HasBasicAuth() bool {
	return len(c.Username) != 0
}

func (c *TransportConfig) Wrap(fn WrapperFunc) {
	c.WrapTransport = Wrappers(c.WrapTransport, fn)
}

type WrapperFunc func(rt http.RoundTripper) http.RoundTripper

func Wrappers(fns ...WrapperFunc) WrapperFunc {
	if len(fns) == 0 {
		return nil
	}
	if len(fns) == 2 && fns[0] == nil {
		return fns[1]
	}

	return func(rt http.RoundTripper) http.RoundTripper {
		base := rt
		for _, fn := range fns {
			if fn != nil {
				base = fn(base)
			}
		}
		return base
	}
}

func NewTransport(config *TransportConfig) (http.RoundTripper, error) {

	var (
		rt http.RoundTripper
	)
	// 可以考虑cache transport
	if config.Transport != nil {
		rt = config.Transport
	} else {
		rt = NewDefaultTransport(config)
	}

	return HTTPWrappersForConfig(config, rt)
}

func TransportFor(config *Config) (http.RoundTripper, error) {
	conf, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	return NewTransport(conf)
}

func HTTPWrappersForConfig(config *TransportConfig, rt http.RoundTripper) (http.RoundTripper, error) {
	if config.WrapTransport != nil {
		rt = config.WrapTransport(rt)
	}
	if len(config.UserAgent) > 0 {
		rt = NewUserAgentRoundTripper(config.UserAgent, rt)
	}
	switch {
	case config.HasBasicAuth():
		rt = NewBasicAuthRoundTripper(config.Username, config.Passowrd, rt)
	}

	return rt, nil
}

func NewDefaultTransport(config *TransportConfig) http.RoundTripper {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, time.Second*3)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   128,
		MaxIdleConns:          2048,
		IdleConnTimeout:       time.Second * 90,
		ExpectContinueTimeout: 5 * time.Second,
	}
}

type userAgentRoundTripper struct {
	ua string
	rt http.RoundTripper
}

func NewUserAgentRoundTripper(agent string, rt http.RoundTripper) http.RoundTripper {
	return &userAgentRoundTripper{agent, rt}
}
func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("User-Agent")) != 0 {
		return rt.rt.RoundTrip(req)
	}
	req = CloneRequest(req)
	req.Header.Set("User-Agent", rt.ua)
	return rt.rt.RoundTrip(req)
}

type basicAuthRoundTripper struct {
	username string
	password string `datapolicy:"password"`
	rt       http.RoundTripper
}

func NewBasicAuthRoundTripper(username, password string, rt http.RoundTripper) http.RoundTripper {
	return &basicAuthRoundTripper{username, password, rt}
}

func (rt *basicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("Authorization")) != 0 {
		return rt.rt.RoundTrip(req)
	}
	req = CloneRequest(req)
	req.SetBasicAuth(rt.username, rt.password)
	return rt.rt.RoundTrip(req)
}
