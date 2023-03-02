package beautify

import "github.com/auho/go-toolkit/farmtools/convert/beautify/decode"

var _ decode.Beautifier = (*Json)(nil)

type Json struct {
}

func NewJsonDecoder() *decode.Beautify {
	return decode.NewDecode(&Json{})
}

func (j *Json) Key(s string) string {
	return keySymbol + s + keySymbol
}

func (j *Json) BoolValue(s string) string {
	return s
}

func (j *Json) IntValue(s string) string {
	return s
}

func (j *Json) UintValue(s string) string {
	return s
}

func (j *Json) Float32Value(s string) string {
	return s
}

func (j *Json) Float64Value(s string) string {
	return s
}

func (j *Json) ComplexValue(s string) string {
	return s
}

func (j *Json) StringValue(s string) string {
	return s
}

func (j *Json) SliceBegin() string {
	return "["
}

func (j *Json) SliceEnd() string {
	return "]"
}

func (j *Json) SliceValue(s string) string {
	return s
}

func (j *Json) SliceValueSeparator() string {
	return ", "
}

func (j *Json) SliceSeparator() string {
	return ""
}

func (j *Json) MapBegin() string {
	return "{"
}

func (j *Json) MapEnd() string {
	return "}"
}

func (j *Json) MapKey(s string) string {
	return j.Key(s)
}

func (j *Json) MapValue(s string) string {
	return s
}

func (j *Json) MapKeySeparator() string {
	return ": "
}

func (j *Json) MapValueSeparator() string {
	return ", "
}
func (j *Json) MapSeparator() string {
	return ""
}

func (j *Json) StructBegin() string {
	return "{"
}

func (j *Json) StructEnd() string {
	return "}"
}

func (j *Json) StructFiled(s string) string {
	return j.Key(s)
}

func (j *Json) StructValue(s string) string {
	return s
}

func (j *Json) StructFiledSeparator() string {
	return ": "
}

func (j *Json) StructValueSeparator() string {
	return ", "
}

func (j *Json) StructSeparator() string {
	return ""
}

func (j *Json) ChanValue(s string) string {
	return j.Key(s)
}
