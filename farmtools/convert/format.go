package convert

import (
	"fmt"
	"reflect"
)

func format(value reflect.Value) string {
	elemKind := value.Kind()

	s := ""
	switch {
	case isIndirect(elemKind):
		s = fmt.Sprintf("%v", value)
	case elemKind == reflect.String:
		s = fmt.Sprintf(`"%v"`, value)
	case isArrayOrSlice(elemKind):
		s = formatSlice(value)
	case elemKind == reflect.Map:
		s = formatMap(value)
	case elemKind == reflect.Struct:
		s = formatStruct(value)
	case elemKind == reflect.Chan:
		s = formatChan(value)
	case elemKind == reflect.Pointer:
		if value.IsNil() {
			//value.Type().Elem().Kind()
			s = CObjectEmpty
		} else {
			elemKind = value.Elem().Kind()
			switch elemKind {
			case reflect.Struct:
				s = format(value.Elem())
			case reflect.Pointer:
				s = format(value.Elem())
			default:
				s = format(value.Elem())
			}
		}
	case elemKind == reflect.Interface:
		if value.IsNil() {
			s = CNull
		} else {
			s = format(value.Elem())
		}
	default:
		s = "__UNKNOWN_TYPE__"
	}

	return s
}
