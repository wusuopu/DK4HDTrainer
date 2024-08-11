package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/valyala/fastjson"
)

// lodash 风格的 GET 操作
func GetJSONInt64(val *fastjson.Value, keys string) (int64, error) {
	k := strings.Split(keys, ".")
	if !val.Exists(k...) {
		return 0, fmt.Errorf("key %s not exists", keys)
	}
	item := val.Get(k...)
	if item.Type() == fastjson.TypeNumber {
		return item.Int64()
	}
	if item.Type() == fastjson.TypeString {
		strconv.ParseInt(item.String(), 10, 64)
	}
	return 0, fmt.Errorf("key %s in not number", keys)
}
func GetJSONFloat64(val *fastjson.Value, keys string) (float64, error) {
	k := strings.Split(keys, ".")
	if !val.Exists(k...) {
		return 0, fmt.Errorf("key %s not exists", keys)
	}
	item := val.Get(k...)
	if item.Type() == fastjson.TypeNumber {
		return item.Float64()
	}
	if item.Type() == fastjson.TypeString {
		strconv.ParseFloat(item.String(), 64)
	}
	return 0, fmt.Errorf("key %s in not number", keys)
}
func GetJSONItem(val *fastjson.Value, keys string) *fastjson.Value {
	return val.Get(strings.Split(keys, ".")...)
}
func GetJSONString(val *fastjson.Value, keys string) string {
	v := GetJSONItem(val, keys)
	if v.Type() != fastjson.TypeString {
		return ""
	}

	ret := v.String()
	ret, err := strconv.Unquote(ret)
	if err != nil {
		return ""
	}
	return ret
}
