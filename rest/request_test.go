package rest

import (
	"mime"
	"testing"
)

func TestContentType(t *testing.T) {
	contentType := "text/html; charset=utf-8"
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(mediaType)
	t.Log(params)

}
