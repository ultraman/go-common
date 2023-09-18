package rest

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestUrlValues(t *testing.T) {
	values := url.Values{}
	values.Set("key", "hello")
	values.Add("key", "world")
	key := values.Get("key")
	t.Log(key)
}

type User struct {
	Name string `json:"name" param:"name"`
	Age  int    `json:"age" param:"age"`
}

func TestReflect(t *testing.T) {
	user := &User{Name: "hello", Age: 18}
	val := reflect.ValueOf(user)
	fmt.Println(val.Kind()) // ptr
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			fmt.Println("val isnil")
			return
		}
		val = val.Elem()
		fmt.Println(val.Kind()) // struct
	}
	typ := val.Type()
	fmt.Println(typ) // User
	// 遍历所有的字段
	for i := 0; i < val.NumField(); i++ {
		// key
		kt := typ.Field(i)
		fmt.Println(kt)
		tag := kt.Tag.Get("param")
		if tag == "-" {
			continue
		}
		fmt.Println("tag", tag)
		// value
		sv := val.Field(i)
		fmt.Println(sv)
		switch sv.Kind() {
		case reflect.String:
			sv.SetString("man")
		}
	}
	fmt.Println(user)
}

func TestParameterCodec_EncodeParameters(t *testing.T) {
	user := &User{Name: "hello", Age: 18}
	codec := parameterCodec{}
	values, err := codec.EncodeParameters(user)
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Println(values)
}

func TestParameterCodec_DecodeParameters(t *testing.T) {
	user := &User{}
	codec := parameterCodec{}
	values := url.Values{}
	values.Set("name", "man")
	values.Set("age", "18")
	err := codec.DecodeParameters(values, user)
	t.Log(err)
	if err == nil {
		t.Log(user)
		t.Log(reflect.TypeOf(user.Age))
	}
}
