package rest

import (
	"net/http"
	"net/url"
	"strings"
)

var (
	DefaultContentType = "application/json"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Interface interface {
	GetRateLimiter() RateLimiter
	Verb(verb string) *Request
	Post() *Request
	Put() *Request
	Patch() *Request
	Get() *Request
	Delete() *Request
}

type Client struct {
	base        *url.URL
	rateLimiter RateLimiter
	Config      ContentConfig
	Client      *http.Client
}

func (c *Client) GetRateLimiter() RateLimiter {
	return nil
}

func (c *Client) Post() *Request {
	return c.Verb("POST")
}

func (c *Client) Put() *Request {
	return c.Verb("PUT")
}

func (c *Client) Patch() *Request {
	return c.Verb("PATCH")
}

func (c *Client) Get() *Request {
	return c.Verb("GET")
}

func (c *Client) Delete() *Request {
	return c.Verb("DELETE")
}

func (c *Client) Verb(verb string) *Request {
	return NewRequest(c).Verb(verb)
}

func NewRESTClient(baseURL *url.URL, config ContentConfig, rateLimiter RateLimiter, client *http.Client) (Interface, error) {
	if len(config.ContentType) == 0 {
		config.ContentType = DefaultContentType
	}
	base := *baseURL
	// 判断后缀是否有 /
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	base.RawQuery = ""
	base.Fragment = ""

	return &Client{
		base:   &base,
		Config: config,
		Client: client,
	}, nil
}

func NewRESTClientFor(config *Config) (Interface, error) {
	// 解析Host
	baseURL, err := ParseUrl(config)
	if err != nil {
		return nil, err
	}

	transport, err := TransportFor(config)
	if err != nil {
		return nil, err
	}

	var httpClient *http.Client
	if transport != http.DefaultTransport {
		httpClient = &http.Client{
			Transport: transport,
		}
		if config.Timeout > 0 {
			httpClient.Timeout = config.Timeout
		}
	}
	// TODO: rateLimiter

	return NewRESTClient(baseURL, config.ContentConfig, nil, httpClient)
}
