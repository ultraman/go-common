package rest

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
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

func (a *nullAuthProvider) RoundTrip(request *http.Request) (*http.Response, error) {
	return a.rt.RoundTrip(request)
}

func (a *nullAuthProvider) WrapTransport(rt http.RoundTripper) http.RoundTripper {
	return a
}

func (a nullAuthProvider) Login() error {
	return nil
}

func NewNullAuthProvider() AuthProvider {
	return &nullAuthProvider{}
}

type serviceAuthTokenProvider struct {
	config *AuthConfig
	rt     http.RoundTripper
}

func (p *serviceAuthTokenProvider) RoundTrip(request *http.Request) (*http.Response, error) {
	token := p.getToken()
	request.Header.Set("server-token", token)
	return p.rt.RoundTrip(request)
}

func (p *serviceAuthTokenProvider) getToken() string {
	ak := p.config.Config["ak"]
	sk := p.config.Config["sk"]
	nonce := uuid.NewV4()
	timestamp := time.Now().Unix()
	sign := fmt.Sprintf("%s/%d/%s", sk, timestamp, nonce)
	signMd5 := md5Util(sign)
	token := fmt.Sprintf("%s/%d/%s/%s/", ak, timestamp, nonce, signMd5)
	return token
}

func (p *serviceAuthTokenProvider) WrapTransport(rt http.RoundTripper) http.RoundTripper {
	p.rt = rt
	return p
}

func (p *serviceAuthTokenProvider) Login() error {
	return nil
}

func NewServiceAuthTokenProvider(config *AuthConfig) AuthProvider {
	return &serviceAuthTokenProvider{config: config}
}

func md5Util(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
