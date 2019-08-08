package util

import (
	"encoding/json"
	"strconv"
)

func StrToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func StrToInt64IgnoreErr(str string) int64 {
	i, _ := StrToInt64(str)
	return i
}

func ToJsonIgnoreErr(obj interface{}) string {
	resBytes, _ := json.Marshal(obj)
	return string(resBytes)
}
