package rest

import (
	"context"
	"testing"
)

type UserResq struct{}
type UserResp struct{}

func TestNewRESTClient(t *testing.T) {
	config := &Config{Host: "http://localhost"}
	client, _ := NewRESTClientFor(config)
	err := client.Post().
		Path("/user/api/v1/info").Body("xxx").Do(context.TODO()).Into(resp)
}

func TestClient_Get(t *testing.T) {
	config := &Config{Host: "http://localhost:8080"}
	client, err := NewRESTClientFor(config)
	if err != nil {
		t.Log(err)
		return
	}
	resp := &UserResp{}
	ctx := context.Background()
	err = client.Post().
		Path("/user/api/v1/info").
		Param("user", "yao").
		Param("age", "18").
		Do(ctx).
		Into(&resp)
	if err != nil {
		t.Log(err)
		return
	}
}
