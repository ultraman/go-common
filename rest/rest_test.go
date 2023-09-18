package rest

import (
	"context"
	"testing"
)

type UserResq struct{}
type UserResp struct{}

func TestClient_Post(t *testing.T) {
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
		Param("user", "man").
		Param("age", "18").
		Do(ctx).
		Into(&resp)
	if err != nil {
		t.Log(err)
		return
	}
}
