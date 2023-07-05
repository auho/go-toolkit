package maps

import (
	"errors"
	"fmt"
	"github.com/auho/go-toolkit/farmtools/convert/types/reflects/values/floats"
	"github.com/auho/go-toolkit/farmtools/convert/types/reflects/values/ints"
	"github.com/auho/go-toolkit/farmtools/convert/types/reflects/values/strings"
	"github.com/auho/go-toolkit/farmtools/convert/types/reflects/values/structs"

	"reflect"
)

// MapStringAnyToStruct convert string any map to any struct
// s any must be a pointer
func MapStringAnyToStruct(s any, from map[string]any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%s", r))
		}
	}()

	sRef := reflect.ValueOf(s)
	if sRef.Kind() != reflect.Ptr {
		return errors.New("input is not pointer")
	}

	if sRef.IsNil() {
		return errors.New("input is nil")
	}

	sRefElem := sRef.Elem()
	if sRefElem.Kind() != reflect.Struct {
		return errors.New("input is not struct")
	}

	sRefElemType := sRefElem.Type()
	fieldNum := sRefElem.NumField()

	var fieldRef reflect.Value
	var fieldType reflect.Kind
	var fieldName string

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("[%s] %s", fieldName, r))
		}
	}()

	for i := 0; i < fieldNum; i++ {
		fieldRef = sRefElem.Field(i)
		fieldType = fieldRef.Kind()
		fieldName = sRefElemType.Field(i).Tag.Get("json")

		if v, ok := from[fieldName]; ok {
			switch fieldType {
			case reflect.Int, reflect.Int64:
				err = ints.FromAny(fieldRef, v)
			case reflect.Float32, reflect.Float64:
				err = floats.FromAny(fieldRef, v)
			case reflect.String:
				err = strings.FromAny(fieldRef, v)
			case reflect.Struct:
				err = structs.FromAny(fieldRef, v)
			default:
				err = errors.New(fmt.Sprintf("unknow type %s[%T %v]", fieldType.String(), v, v))
			}
		}
	}

	return err
}
