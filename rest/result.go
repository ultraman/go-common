package rest

import (
	"fmt"
	"live-works/pkg/codec"
)

type IResult interface {
	GetCodec() codec.Marshaler
	Raw() ([]byte, error)
	Into(obj interface{}) error
	StatusCode() int
	Error() error
}

type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
	codecer     codec.Marshaler
}

func (r *Result) GetCodec() codec.Marshaler {
	return r.codecer
}

func (r Result) Raw() ([]byte, error) {
	return r.body, r.err
}

func (r Result) Into(obj interface{}) error {
	if r.err != nil {
		return r.err
	}
	if r.codecer == nil {
		err := fmt.Errorf("serializer for %s doesn't exist", r.contentType)
		r.err = err
		return err
	}
	if len(r.body) == 0 {
		err := fmt.Errorf("0-length response with status code:%d and content type:%s", r.statusCode, r.contentType)
		r.err = err
		return err
	}
	fmt.Println(string(r.body))
	err := r.codecer.Unmarshal(r.body, obj)
	if err != nil {
		r.err = err
		return err
	}
	return nil
}

func (r Result) StatusCode() int {
	return r.statusCode
}
func (r Result) Error() error {
	return r.err
}
