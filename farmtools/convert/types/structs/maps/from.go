package maps

import (
	"errors"
	"fmt"
	"reflect"
)

// MapStringAnyFromStruct convert struct to string any map
func MapStringAnyFromStruct(s interface{}) (m map[string]any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%s", r))
		}
	}()

	var sRef, sRefElem reflect.Value
	var sRefElemType reflect.Type
	var fieldNum int

	sRef = reflect.ValueOf(s)
	if sRef.Kind() == reflect.Ptr {
		if sRef.IsNil() {
			return nil, errors.New("input is nil")
		}

		sRefElem = sRef.Elem()
	} else {
		sRefElem = sRef
	}

	if sRefElem.Kind() != reflect.Struct {
		return nil, errors.New("input is not struct")
	}

	sRefElemType = sRefElem.Type()
	fieldNum = sRefElem.NumField()

	m = make(map[string]interface{})

	for i := 0; i < fieldNum; i++ {
		fieldRef := sRefElem.Field(i)
		fieldName := sRefElemType.Field(i).Tag.Get("json")
		if fieldName == "" {
			fieldName = sRefElemType.Field(i).Name
		}

		m[fieldName] = fieldRef.Interface()
	}

	return
}
