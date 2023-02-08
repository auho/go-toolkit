package hommization

import (
	"reflect"
)

const CNull = "null"
const CObjectEmpty = "{}"
const CBoolFalse = "false"
const cBoolTrue = "true"

func isFigure(kind reflect.Kind) bool {
	return (kind >= reflect.Int && kind <= reflect.Uint64) || (kind >= reflect.Float32 && kind <= reflect.Complex128)
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
