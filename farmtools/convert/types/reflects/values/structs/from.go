package structs

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func FromAny(fieldRef reflect.Value, from interface{}) error {
	var err error

	switch fieldRef.Interface().(type) {
	case time.Time:
		switch nv := from.(type) {
		case string:
			timeLayout := "2006-01-02 15:04:05"
			loc, _ := time.LoadLocation("Local")
			nt, _ := time.ParseInLocation(timeLayout, nv, loc)
			fieldRef.Set(reflect.ValueOf(nt))
		case int, float32, float64:
			fieldRef.Set(reflect.ValueOf(time.Unix(0, 0)))
		default:
			err = errors.New(fmt.Sprintf("convert time.Time error[%T %v]", from, from))
		}

	default:
		err = errors.New(fmt.Sprintf("convert struct type error[%T %v]", from, from))
	}

	return err
}
