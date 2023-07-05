package ints

import (
	"errors"
	"fmt"
	"strconv"
)

func FromAny(v any) (int, error) {
	var err error
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
		newV, err = strconv.Atoi(_v)
		err = fmt.Errorf("convert string to int error %w", err)
	default:
		err = errors.New(fmt.Sprintf("convert int type error[%T %v]", v, v))
	}

	return newV, err
}
