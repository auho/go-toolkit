package redis

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ConvertFromHash(s interface{}, m map[string]interface{}) (err error) {
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

		if v, ok := m[fieldName]; ok {
			switch fieldType {
			case reflect.Int, reflect.Int64:
				_convertInt(fieldRef, v)
			case reflect.Float32, reflect.Float64:
				_convertFloat(fieldRef, v)
			case reflect.String:
				_convertString(fieldRef, v)
			case reflect.Struct:
				_convertStruct(fieldRef, v)
			default:
				panic(fmt.Sprintf("unknow type %s[%T %v]", fieldType.String(), v, v))
			}
		}
	}

	return
}

func _convertInt(fieldRef reflect.Value, v interface{}) {
	tempInt := int64(0)

	switch tmpVal := v.(type) {
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
		panic(fmt.Sprintf("convert int type error[%T %v]", v, v))
	}

	fieldRef.SetInt(tempInt)
}

func _convertFloat(fieldRef reflect.Value, v interface{}) {
	tempFloat := float64(0)

	switch tmpVal := v.(type) {
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
		panic(fmt.Sprintf("convert float type error[%T %v]", v, v))
	}

	fieldRef.SetFloat(tempFloat)
}

func _convertString(fieldRef reflect.Value, v interface{}) {
	tempString := ""

	switch tmpVal := v.(type) {
	case string:
		tempString = tmpVal
	case float64:
		tempString = strconv.FormatFloat(tmpVal, 'E', -1, 64)
	default:
		panic(fmt.Sprintf("convert string type error[%T %v]", v, v))
	}

	fieldRef.SetString(tempString)
}

func _convertStruct(fieldRef reflect.Value, v interface{}) {
	switch fieldRef.Interface().(type) {
	case time.Time:
		switch nv := v.(type) {
		case string:
			timeLayout := "2006-01-02 15:04:05"
			loc, _ := time.LoadLocation("Local")
			nt, _ := time.ParseInLocation(timeLayout, nv, loc)
			fieldRef.Set(reflect.ValueOf(nt))
		case int, float32, float64:
			fieldRef.Set(reflect.ValueOf(time.Unix(0, 0)))
		default:
			panic(fmt.Sprintf("convert time.Time error[%T %v]", v, v))
		}

	default:
		panic(fmt.Sprintf("convert struct type error[%T %v]", v, v))
	}
}
