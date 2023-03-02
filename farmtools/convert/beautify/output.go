package beautify

import "github.com/auho/go-toolkit/farmtools/convert/beautify/decode"

var _ decode.Beautifier = (*Output)(nil)

type Output struct {
}

func NewOutputDecoder() *decode.Beautify {
	return decode.NewDecode(&Output{})
}

func (o *Output) Key(s string) string {
	return keySymbol + s + keySymbol
}

func (o *Output) BoolValue(s string) string {
	return s
}

func (o *Output) IntValue(s string) string {
	return s
}

func (o *Output) UintValue(s string) string {
	return s
}

func (o *Output) Float32Value(s string) string {
	return s
}

func (o *Output) Float64Value(s string) string {
	return s
}

func (o *Output) ComplexValue(s string) string {
	return s
}

func (o *Output) StringValue(s string) string {
	return s
}

func (o *Output) SliceBegin() string {
	return "["
}

func (o *Output) SliceEnd() string {
	return "]"
}

func (o *Output) SliceValue(s string) string {
	return s
}

func (o *Output) SliceValueSeparator() string {
	return ", "
}

func (o *Output) SliceSeparator() string {
	return "\n"
}

func (o *Output) MapBegin() string {
	return "{"
}

func (o *Output) MapEnd() string {
	return "}"
}

func (o *Output) MapKey(s string) string {
	return o.Key(s)
}

func (o *Output) MapValue(s string) string {
	return s
}

func (o *Output) MapKeySeparator() string {
	return ": "
}

func (o *Output) MapValueSeparator() string {
	return ", "
}

func (o *Output) MapSeparator() string {
	return "\n"
}

func (o *Output) StructBegin() string {
	return "{"
}

func (o *Output) StructEnd() string {
	return "}"
}

func (o *Output) StructFiled(s string) string {
	return o.Key(s)
}

func (o *Output) StructValue(s string) string {
	return s
}

func (o *Output) StructFiledSeparator() string {
	return ": "
}

func (o *Output) StructValueSeparator() string {
	return ", "
}

func (o *Output) StructSeparator() string {
	return "\n"
}

func (o *Output) ChanValue(s string) string {
	return o.Key(s)
}
