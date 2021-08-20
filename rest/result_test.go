package rest

import (
	"testing"
)

func TestResult_Raw(t *testing.T) {
	result := &Result{
		body:        []byte(`{"name":"yao"}`),
		statusCode:  200,
		err:         nil,
		contentType: DefaultContentType,
		codecer:     NewJsonMarshaler(),
	}
	t.Log(result.Raw())
	user := &struct {
		Name string `json:"name"`
	}{}
	t.Log(result.Into(user))
	t.Log(user.Name)
}
