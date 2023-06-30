package strings

import (
	"fmt"
	"strconv"
)

func AnyTo(v any) string {
	newV := ""

	switch _v := v.(type) {
	case string:
		newV = _v
	case int:
		newV = strconv.Itoa(_v)
	case float64:
		newV = strconv.FormatFloat(_v, 'f', -1, 64)
	default:
		panic(fmt.Sprintf("convert string type error[%T %v]", v, v))
	}

	return newV
}
