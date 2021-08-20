package rest

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

func CloneRequest(req *http.Request) *http.Request {
	r := new(http.Request)

	// shallow clone
	*r = *req

	// deep copy headers
	r.Header = CloneHeader(req.Header)

	return r
}

func CloneHeader(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	return out
}

func md5Util(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
