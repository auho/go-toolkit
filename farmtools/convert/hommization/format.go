package hommization

import (
	"fmt"
	"reflect"
	"strconv"
)

func Convert(a any) string {
	aRef := reflect.ValueOf(a)

	return Format(aRef)
}

func Format(value reflect.Value) string {
	elemKind := value.Kind()

	s := ""
	switch elemKind {
	case reflect.Bool:
		if value.Bool() {
			s = cBoolTrue
		} else {
			s = CBoolFalse
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		s = strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.Float64:
		s = strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		s = fmt.Sprintf("%v", value)
	case reflect.String:
		s = `"` + value.String() + `"`
	case reflect.Slice, reflect.Array:
		s = formatSlice(value)
	case reflect.Map:
		s = formatMap(value)
	case reflect.Struct:
		s = formatStruct(value)
	case reflect.Chan:
		s = formatChan(value)
	case reflect.Pointer:
		if value.IsNil() {
			//value.Type().Elem().Kind()
			s = CObjectEmpty
		} else {
			s = Format(value.Elem())
		}
	case reflect.Interface:
		if value.IsNil() {
			s = CNull
		} else {
			s = Format(value.Elem())
		}
	default:
		s = "__UNKNOWN_TYPE__"
	}

	return s
}
