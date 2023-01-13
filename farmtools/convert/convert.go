package convert

import (
	"reflect"
)

const CNull = "null"
const CObjectEmpty = "{}"

func isFigure(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func isIndirect(kind reflect.Kind) bool {
	return isFigure(kind) || kind == reflect.Bool
}

func isLiteral(kind reflect.Kind) bool {
	return isIndirect(kind) || kind == reflect.String
}

func isArrayOrSlice(kind reflect.Kind) bool {
	return kind == reflect.Slice || kind == reflect.Array
}

func typeStringToSymbolString(s string) string {
	return addDoubleQuotationMark(addTypeSymbol(s))
}

func pointerStringToSymbolString(s string) string {
	return addDoubleQuotationMark(addPointerSymbol(s))
}

func addDoubleQuotationMark(s string) string {
	return `"` + s + `"`
}

func addTypeSymbol(s string) string {
	return "<" + s + ">"
}

func addPointerSymbol(s string) string {
	return "*" + s
}

func addBraces(s string) string {
	return "{" + s + "}"
}

func addBrackets(s string) string {
	return "[" + s + "]"
}
