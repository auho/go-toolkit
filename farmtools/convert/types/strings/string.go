package strings

import (
	"errors"
	"fmt"
	"strconv"
)

func FromAny(v any) (string, error) {
	var err error
	newV := ""

	switch _v := v.(type) {
	case string:
		newV = _v
	case []uint8:
		newV = string(_v)
	case int:
		newV = strconv.Itoa(_v)
	case int64:
		newV = strconv.FormatInt(_v, 10)
	case float64:
		newV = strconv.FormatFloat(_v, 'f', -1, 64)
	default:
		err = errors.New(fmt.Sprintf("convert string type error[%T %v]", v, v))
	}

	return newV, err
}
