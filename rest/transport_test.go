package rest

import (
	"fmt"
	"net/http"
	"testing"
)

func TestWrappers(t *testing.T) {
	wrap := Wrappers()
	wrap = Wrappers(wrap, func(rt http.RoundTripper) http.RoundTripper {
		fmt.Println("start")
		fmt.Println("end")
		return rt
	})
}
