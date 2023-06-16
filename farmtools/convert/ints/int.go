package ints

import (
	"fmt"
	"strconv"
)

func AnyTo(v any) int {
	newV := 0

	switch _v := v.(type) {
	case int:
		newV = _v
	case int8:
		newV = int(_v)
	case int16:
		newV = int(_v)
	case int32:
		newV = int(_v)
	case int64:
		newV = int(_v)
	case float32:
		newV = int(_v)
	case float64:
		newV = int(_v)
	case string:
		newV, _ = strconv.Atoi(_v)
	default:
		panic(fmt.Sprintf("convert int type error[%T %v]", v, v))
	}

	return newV
}
