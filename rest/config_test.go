package rest

import (
	"path"
	"strings"
	"testing"
)

func TestDefaultServerURL(t *testing.T) {
	host := "http://192.168.100.200:8080/user/api/v1/login?user=hello&age=18"
	url, _ := defaultServerURL(host)
	t.Log(url.Host)     // 192.168.100.200:8080
	t.Log(url.Path)     // /user/api/v1/login
	t.Log(url.RawQuery) // user=hello&age=18
	t.Log(url.Scheme)   // http
}

func TestDefaultStrings(t *testing.T) {
	host := "http://192.168.100.200:8080"
	base, _ := defaultServerURL(host)
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	t.Log(base.Path)
}

func TestDefaultPath(t *testing.T) {
	host := "http://192.168.100.200:8080"
	base, _ := defaultServerURL(host)
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	versionAPIPath := "v1"
	pathPrefix := path.Join("/", base.Path, versionAPIPath)
	t.Log(pathPrefix)

	pathPrefix = path.Join(pathPrefix, path.Join("a", "b", "c"))
	t.Log(pathPrefix)
}
