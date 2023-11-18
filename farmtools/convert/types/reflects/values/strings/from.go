package strings

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func FromAny(fieldRef reflect.Value, from interface{}) error {
	var err error
	tempString := ""

	if from != nil {
		switch tmpVal := from.(type) {
		case string:
			tempString = tmpVal
		case float64:
			tempString = strconv.FormatFloat(tmpVal, 'f', -1, 64)
		default:
			err = errors.New(fmt.Sprintf("convert string type error[%T %v]", from, from))
		}

		if err != nil {
			return err
		}
	}

	fieldRef.SetString(tempString)

	return nil
}
