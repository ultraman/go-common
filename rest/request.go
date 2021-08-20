package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"live-works/pkg/codec"
	"live-works/pkg/logger"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Request struct {
	c *Client

	timeout     time.Duration
	rateLimiter RateLimiter
	// url params
	verb       string
	pathPrefix string
	params     url.Values
	headers    http.Header
	// output
	err   error
	body  io.Reader
	retry WithRetry

	coder codec.Marshaler
}

// Verb 请求类型
func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

// SetHeader 设置请求头
func (r *Request) SetHeader(key string, values ...string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Del(key)
	for _, value := range values {
		r.headers.Add(key, value)
	}
	return r
}

// Timeout 设置超时
func (r *Request) Timeout(d time.Duration) *Request {
	if r.err != nil {
		return r
	}
	r.timeout = d

	return r
}

func (r *Request) Prefix(segments ...string) *Request {
	if r.err != nil {
		return r
	}
	// segments = [a,b,c]
	// pathPrefix = /api/v1/ + "/a/b/c"
	r.pathPrefix = path.Join(r.pathPrefix, path.Join(segments...))
	return r
}

// RequestURI 完整请求地址 + 参数 覆盖现有地址与参数
func (r *Request) RequestURI(uri string) *Request {
	if r.err != nil {
		return r
	}
	locator, err := url.Parse(uri)
	if err != nil {
		r.err = err
		return r
	}

	r.pathPrefix = locator.Path

	if len(locator.Query()) > 0 {
		if r.params == nil {
			r.params = make(url.Values)
		}
		for k, v := range locator.Query() {
			r.params[k] = v
		}
	}
	return r
}

// Path 设置请求路径
func (r *Request) Path(url string) *Request {
	r.pathPrefix = url
	return r
}

// Param 请求参数
func (r *Request) Param(name, value string) *Request {
	if r.err != nil {
		return r
	}
	return r.setParam(name, value)
}

// Params 请求参数
func (r *Request) Params(obj interface{}, codec ParameterCodec) *Request {
	if obj == nil {
		return r
	}
	params, err := codec.EncodeParameters(obj)
	if err != nil {
		r.err = err
		return r
	}
	if r.params == nil {
		r.params = make(url.Values)
	}
	for key, values := range params {
		r.params[key] = append(r.params[key], values...)
	}
	return r
}

func (r *Request) setParam(name, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[name] = append(r.params[name], value)
	return r
}

// URL 生成URL对象
func (r *Request) URL() *url.URL {
	path := r.pathPrefix

	finalURL := &url.URL{}
	if r.c.base != nil {
		*finalURL = *r.c.base
	}
	finalURL.Path = path

	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	finalURL.RawQuery = query.Encode()
	return finalURL
}

// Body 请求体
func (r *Request) Body(obj interface{}) *Request {
	if r.err != nil {
		return r
	}

	switch t := obj.(type) {
	case []byte:
		r.body = bytes.NewReader(t)
	case io.Reader:
		r.body = t
	default:
		b, err := r.coder.Marshal(obj)
		if err != nil {
			r.err = err
		}
		r.body = bytes.NewReader(b)
	}
	return r
}

// Do 发起请求
func (r *Request) Do(ctx context.Context) Result {
	var result Result
	err := r.request(ctx, func(request *http.Request, response *http.Response) {
		// 返回结果
		result = r.transformResponse(response, request)
	})
	if err != nil {
		return Result{err: err}
	}
	return result
}

// DoRaw 返回字节与错误
func (r *Request) DoRaw(ctx context.Context) ([]byte, error) {
	result := r.Do(ctx)
	return result.Raw()
}

// newHTTPRequest 构造http.Request
func (r *Request) newHTTPRequest(ctx context.Context) (*http.Request, error) {
	u := r.URL().String()
	req, err := http.NewRequest(r.verb, u, r.body)
	if err != nil {
		return nil, err
	}
	req.Header = r.headers
	req.WithContext(ctx)
	return req, nil
}

func (r *Request) request(ctx context.Context, fn func(*http.Request, *http.Response)) error {
	start := time.Now()
	defer func() {
		latency := time.Since(start)
		logger.Infof("latency: %d", latency)
	}()
	if r.err != nil {
		return r.err
	}

	client := r.c.Client
	if client == nil {
		client = http.DefaultClient
	}
	// 超时取消 context
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}
	// 构造Request
	req, err := r.newHTTPRequest(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	// TODO:// 重试 限流 Metric
	done := func() bool {

		f := func(req *http.Request, resp *http.Response) {
			if resp == nil {
				return
			}
			fn(req, resp)
		}

		f(req, resp)

		return true
	}()
	if done {
		return err
	}
	return nil
}

// transformResponse 处理返回结果
func (r *Request) transformResponse(resp *http.Response, req *http.Request) Result {
	var body []byte
	// 先读取Body
	if resp.Body != nil {
		out, err := ioutil.ReadAll(resp.Body)
		switch err.(type) {
		case nil:
			body = out
		default:
			err = fmt.Errorf("unexpected error when reading response body. Please retry. Original error: %w", err)
			return Result{err: err}
		}
	}

	coder := r.coder

	contentType := resp.Header.Get("Content-Type")
	if len(contentType) == 0 {
		contentType = r.c.Config.ContentType
	}
	return Result{
		body:        body,
		contentType: contentType,
		statusCode:  resp.StatusCode,
		codecer:     coder,
	}
}

func NewRequest(c *Client) *Request {

	var timeout time.Duration

	if c.Client != nil {
		timeout = c.Client.Timeout
	}
	coder := c.Config.Codec

	var pathPrefix string

	r := &Request{
		c:           c,
		rateLimiter: c.rateLimiter,
		timeout:     timeout,
		pathPrefix:  pathPrefix,
		coder:       coder,
		retry:       NewWithRetry(10),
	}
	switch {
	case len(c.Config.AcceptContentTypes) > 0:
		r.SetHeader("Accept", c.Config.AcceptContentTypes+", */*")
	case len(c.Config.ContentType) > 0:
		r.SetHeader("Accept", c.Config.ContentType+", */*")
	}
	return r
}
