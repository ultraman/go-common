package rest

import (
	"net/http"
)

type AuthConfig struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config,omitempty"`
}

type AuthProvider interface {
	WrapTransport(http.RoundTripper) http.RoundTripper
	Login() error
}

type nullAuthProvider struct {
	rt http.RoundTripper
}

func (p *nullAuthProvider) RoundTrip(request *http.Request) (*http.Response, error) {
	return p.rt.RoundTrip(request)
}

func (p *nullAuthProvider) WrapTransport(rt http.RoundTripper) http.RoundTripper {
	return p
}

func (p nullAuthProvider) Login() error {
	return nil
}

func NewNullAuthProvider() AuthProvider {
	return &nullAuthProvider{}
}
