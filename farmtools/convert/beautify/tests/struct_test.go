package tests

import (
	"testing"
)

func Test_struct(t *testing.T) {
	_testStruct(t, nil)
}

func Benchmark_struct(b *testing.B) {
	_testStruct(nil, b)
}

func _testStruct(t *testing.T, b *testing.B) {
	type args struct {
		value interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"struct empty", args{_structEmpty}, ""},
		{"struct", args{_struct}, ""},
		{"struct pointer empty", args{&_structEmpty}, ""},
		{"struct pointer", args{&_struct}, ""},
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
