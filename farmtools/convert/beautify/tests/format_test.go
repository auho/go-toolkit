package tests

import (
	"testing"
)

func Test_format(t *testing.T) {
	_testFormat(t, nil)
}

func Benchmark_format(b *testing.B) {
	_testFormat(nil, b)
}

func _testFormat(t *testing.T, b *testing.B) {
	type args struct {
		value interface{}
	}

	var structAnyEmpty interface{} = _structEmpty
	var structAny interface{} = _struct
	var structPtrAnyNil interface{} = (*structType)(nil)
	var structPtrAnyEmpty interface{} = &_structEmpty
	var structPtrAny interface{} = &_struct

	tests := []struct {
		name string
		args args
		want string
	}{
		{"struct interface empty", args{structAnyEmpty}, ""},

		{"bool true", args{true}, ""},
		{"uint", args{1}, ""},
		{"int", args{-1}, ""},
		{"float 32", args{float32(-1.1)}, ""},
		{"float 64", args{-1.1}, ""},
		{"string", args{"s1"}, ""},

		{"struct interface empty", args{structAnyEmpty}, ""},
		{"struct interface", args{structAny}, ""},
		{"struct pointer interface empty", args{structPtrAnyEmpty}, ""},
		{"struct pointer interface nil", args{structPtrAnyNil}, ""},
		{"struct pointer interface", args{structPtrAny}, ""},
	}

	for _, tt := range tests {
		if t != nil {
			t.Run(tt.name, func(t *testing.T) {
				if got := _decoder.Decode(tt.args.value); got != tt.want {
					_assert(t, got, tt.want)
				}
			})
		}

		if b != nil {
			b.Run(tt.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_decoder.Decode(tt.args.value)
				}
			})
		}
	}
}
