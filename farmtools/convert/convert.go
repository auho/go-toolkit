package convert

import (
	"reflect"
)

const CNull = "null"
const CObjectEmpty = "{}"

func _isFigure(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func _isIndirect(kind reflect.Kind) bool {
	return _isFigure(kind) || kind == reflect.Bool
}

func _isLiteral(kind reflect.Kind) bool {
	return _isIndirect(kind) || kind == reflect.String
}

func _isArrayOrSlice(kind reflect.Kind) bool {
	return kind == reflect.Slice || kind == reflect.Array
}
