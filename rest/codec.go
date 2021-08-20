package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type ParameterCodec interface {
	DecodeParameters(parameters url.Values, obj interface{}) error
	EncodeParameters(obj interface{}) (url.Values, error)
}

func NewParameterCodec() ParameterCodec {
	return &parameterCodec{}
}

var (
	ErrStruct = errors.New("Unmarshal() expects struct input. ")
)

type parameterCodec struct{}

// DecodeParameters  将url.Values填充到obj中
func (p parameterCodec) DecodeParameters(parameters url.Values, obj interface{}) error {
	// check obj is struct
	val := reflect.ValueOf(obj)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ErrStruct
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return ErrStruct
	}
	return reflectUrlValuesToStruct(parameters, val)
}

// EncodeParameters  将obj转为url.Values
func (p parameterCodec) EncodeParameters(obj interface{}) (url.Values, error) {
	values := &url.Values{}
	val := reflect.ValueOf(obj)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return *values, ErrStruct
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return *values, ErrStruct
	}
	return reflectStructToUrlValues(values, val)
}

func reflectStructToUrlValues(values *url.Values, val reflect.Value) (url.Values, error) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		kt := typ.Field(i)
		// 获取Tag
		tag := kt.Tag.Get("param")
		if tag == "-" {
			continue
		}
		name, _ := parseTag(tag)
		sv := val.Field(i)
		// 如果vlaue没有赋值 默认跳过
		if sv.IsZero() {
			continue
		}
		value := ""
		switch sv.Kind() {
		case reflect.String:
			value = sv.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(sv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			value = strconv.FormatUint(sv.Uint(), 10)
		case reflect.Bool:
			value = strconv.FormatBool(sv.Bool())
		}
		values.Set(name, value)
	}
	return *values, nil
}

func reflectUrlValuesToStruct(values url.Values, val reflect.Value) error {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		kt := typ.Field(i)
		tag := kt.Tag.Get("param")
		if tag == "-" {
			continue
		}
		sv := val.Field(i)
		// 根据tag的key 从url.Values里获取对应的value
		uv := getValByTag(values, tag)
		switch sv.Kind() {
		case reflect.String:
			sv.SetString(uv)
		case reflect.Bool:
			b, err := strconv.ParseBool(uv)
			if err != nil {
				return errors.New(fmt.Sprintf("cast bool has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
			}
			sv.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(uv, 10, 64)
			if err != nil || sv.OverflowUint(n) {
				return errors.New(fmt.Sprintf("cast uint has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
			}
			sv.SetUint(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(uv, 10, 64)
			if err != nil || sv.OverflowInt(n) {
				return errors.New(fmt.Sprintf("cast int has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
			}
			sv.SetInt(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(uv, sv.Type().Bits())
			if err != nil || sv.OverflowFloat(n) {
				return errors.New(fmt.Sprintf("cast float has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
			}
			sv.SetFloat(n)
		default:
			return errors.New(fmt.Sprintf("unsupported type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
		}
	}
	return nil
}

func getValByTag(parameters url.Values, tag string) string {
	name, opts := parseTag(tag)
	// 根据key 查找对应指
	v := parameters.Get(name)
	if len(opts) > 0 {
		// 如果有的话 返回默认值
		if len(opts) == 1 && v == "" {
			v = opts[0]
		}
	}
	return v
}

// parseTag 如果 param:"name,omitempty"` 那么 tag = "name,omitempty" 结果 s[0] = name s[1:] = [omitempty,]
func parseTag(tag string) (string, []string) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

var _ ParameterCodec = &parameterCodec{}

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type jsonMarshaler struct{}

var _ Marshaler = &jsonMarshaler{}

func NewJsonMarshaler() Marshaler {
	return &jsonMarshaler{}
}

func (j jsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	bytes, err := json.Marshal(v)
	return bytes, err
}

func (j jsonMarshaler) Unmarshal(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}
