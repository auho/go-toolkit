package floats

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func FromAny(fieldRef reflect.Value, from interface{}) error {
	var err error

	tempFloat := float64(0)

	switch tmpVal := from.(type) {
	case int:
		tempFloat = float64(tmpVal)
	case int64:
		tempFloat = float64(tmpVal)
	case float32:
		tempFloat = float64(tmpVal)
	case float64:
		tempFloat = tmpVal
	case string:
		tempFloat, _ = strconv.ParseFloat(tmpVal, 64)
	default:
		err = errors.New(fmt.Sprintf("convert float type error[%T %v]", from, from))
	}

	if err != nil {
		return err
	}

	fieldRef.SetFloat(tempFloat)

	return nil
}
