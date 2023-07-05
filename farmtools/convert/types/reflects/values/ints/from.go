package ints

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func FromAny(fieldRef reflect.Value, from interface{}) error {
	var err error
	tempInt := int64(0)

	switch tmpVal := from.(type) {
	case int:
		tempInt = int64(tmpVal)
	case int64:
		tempInt = tmpVal
	case float32:
		tempInt = int64(tmpVal)
	case float64:
		tempInt = int64(tmpVal)
	case string:
		tempInt, _ = strconv.ParseInt(tmpVal, 10, 64)
	default:
		err = errors.New(fmt.Sprintf("convert int type error[%T %v]", from, from))
	}

	if err != nil {
		return err
	}

	fieldRef.SetInt(tempInt)

	return nil
}
