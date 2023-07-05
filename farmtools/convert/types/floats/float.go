package floats

import (
	"fmt"
	"strconv"
)

func FromAny(v any) float64 {
	newV := float64(0)

	switch _v := v.(type) {
	case int:
		newV = float64(_v)
	case int64:
		newV = float64(_v)
	case float32:
		newV = float64(_v)
	case float64:
		newV = _v
	case string:
		newV, _ = strconv.ParseFloat(_v, 64)
	default:
		panic(fmt.Sprintf("convert float type error[%T %v]", v, v))
	}

	return newV
}
