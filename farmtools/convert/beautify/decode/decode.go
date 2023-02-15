package decode

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const cNull = "null"
const cBoolFalse = "false"
const cBoolTrue = "true"
const cUnknownType = "__UNKNOWN_TYPE__"

type Beautifier interface {
	Key(string) string

	BoolValue(string) string
	IntValue(string) string
	UintValue(string) string
	Float32Value(string) string
	Float64Value(string) string
	ComplexValue(string) string
	StringValue(string) string

	SliceBegin() string
	SliceEnd() string
	SliceValue(string) string
	SliceValueSeparator() string
	SliceSeparator() string

	MapBegin() string
	MapEnd() string
	MapKey(string) string
	MapValue(string) string
	MapKeySeparator() string
	MapValueSeparator() string
	MapSeparator() string

	StructBegin() string
	StructEnd() string
	StructFiled(string) string
	StructValue(string) string
	StructFiledSeparator() string
	StructValueSeparator() string
	StructSeparator() string

	ChanValue(string) string
}

type Beautify struct {
	decoder Beautifier
}

func NewDecode(d Beautifier) *Beautify {
	b := &Beautify{
		decoder: d,
	}

	return b
}

func (b *Beautify) Decode(a any) string {
	aRef := reflect.ValueOf(a)

	return b.format(aRef)
}

func (b *Beautify) format(value reflect.Value) string {
	elemKind := value.Kind()

	s := ""
	switch elemKind {
	case reflect.Bool:
		if value.Bool() {
			s = cBoolTrue
		} else {
			s = cBoolFalse
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = b.decoder.IntValue(strconv.FormatInt(value.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s = b.decoder.UintValue(strconv.FormatUint(value.Uint(), 10))
	case reflect.Float32:
		s = b.decoder.Float32Value(strconv.FormatFloat(value.Float(), 'E', -1, 32))
	case reflect.Float64:
		s = b.decoder.Float32Value(strconv.FormatFloat(value.Float(), 'E', -1, 64))
	case reflect.Complex64, reflect.Complex128:
		s = b.decoder.ComplexValue(fmt.Sprintf("%v", value))
	case reflect.String:
		s = b.decoder.StringValue(value.String())
	case reflect.Slice, reflect.Array:
		s = b.slice(value)
	case reflect.Map:
		s = b.maps(value)
	case reflect.Struct:
		s = b.structs(value)
	case reflect.Chan:
		s = b.chans(value)
	case reflect.Pointer:
		s = b.pointer(value)
	case reflect.Interface:
		s = b.interfaces(value)
	default:
		s = cUnknownType
	}

	return s
}

func (b *Beautify) slice(value reflect.Value) string {
	sep := ""
	_len := value.Len()
	isLiteral := b.isLiteral(value.Type().Elem().Kind())
	if !isLiteral {
		sep = b.decoder.SliceSeparator()
	}

	var ss []string
	ss = append(ss, b.decoder.SliceBegin())

	for i := 0; i < _len; i++ {
		ss = append(ss, b.format(value.Index(i))+b.decoder.SliceValueSeparator())

	}

	ss = append(ss, b.decoder.SliceEnd())

	return strings.Join(ss, sep)
}

func (b *Beautify) maps(value reflect.Value) string {
	sep := ""
	isLiteral := b.isLiteral(value.Type().Elem().Kind())
	if !isLiteral {
		sep = b.decoder.MapSeparator()
	}

	var ss []string
	ss = append(ss, b.decoder.MapBegin())

	iterator := value.MapRange()
	for iterator.Next() {
		k := b.underlyingKindString(iterator.Key())
		v := b.format(iterator.Value())

		ss = append(ss, b.decoder.MapKey(k)+b.decoder.MapKeySeparator()+b.decoder.MapValue(v)+b.decoder.MapValueSeparator())
	}

	ss = append(ss, b.decoder.MapEnd())

	return strings.Join(ss, sep)
}

func (b *Beautify) structs(value reflect.Value) string {
	_filedNum := value.NumField()
	_valueType := value.Type()

	var ss []string
	ss = append(ss, b.decoder.StructBegin())

	for i := 0; i < _filedNum; i++ {
		fieldElem := value.Field(i)
		fieldName := _valueType.Field(i).Name

		ss = append(ss, b.decoder.StructFiled(fieldName)+b.decoder.StructFiledSeparator()+b.decoder.StructValue(b.format(fieldElem))+
			b.decoder.StructValueSeparator())
	}

	ss = append(ss, b.decoder.StructEnd())

	return strings.Join(ss, b.decoder.StructSeparator())
}

func (b *Beautify) chans(value reflect.Value) string {
	return b.decoder.Key(b.symbolType(value.Type().String()))
}

func (b *Beautify) pointer(value reflect.Value) string {
	if value.IsNil() {
		//value.Type().Elem().Kind()
		return cNull
	} else {
		return b.format(value.Elem())
	}
}

func (b *Beautify) interfaces(value reflect.Value) string {
	if value.IsNil() {
		return cNull
	} else {
		return b.format(value.Elem())
	}
}

func (b *Beautify) underlyingKindString(value reflect.Value) string {
	_isLiteral := false

	kind := value.Kind()
	s := ""
	switch {
	case b.isLiteral(kind):
		s = b.format(value)
		_isLiteral = true
	case kind == reflect.Pointer:
		s = b.underlyingKindString(value.Elem())
		s = b.symbolPointer(s)
	default:
		if value.Type().Kind() == reflect.Interface {
			if value.IsNil() {
				s = cNull
				_isLiteral = true
			} else {
				s = value.Kind().String()
			}
		} else {
			s = value.Type().String()
		}
	}

	if !_isLiteral {
		s = b.symbolType(s)
	}

	return s
}

func (b *Beautify) symbolType(s string) string {
	return "<" + s + ">"
}

func (b *Beautify) symbolPointer(s string) string {
	return "*" + s
}
func (b *Beautify) isFigure(kind reflect.Kind) bool {
	return (kind >= reflect.Int && kind <= reflect.Uint64) || (kind >= reflect.Float32 && kind <= reflect.Complex128)
}

func (b *Beautify) isIndirect(kind reflect.Kind) bool {
	return b.isFigure(kind) || kind == reflect.Bool
}

func (b *Beautify) isLiteral(kind reflect.Kind) bool {
	return b.isIndirect(kind) || kind == reflect.String
}
