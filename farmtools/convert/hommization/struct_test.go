package hommization

import (
	"reflect"
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
		value reflect.Value
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"struct empty", args{reflect.ValueOf(_structEmpty)}, ""},
		{"struct", args{reflect.ValueOf(_struct)}, ""},
		{"struct pointer empty", args{reflect.ValueOf(&_structEmpty)}, ""},
		{"struct pointer", args{reflect.ValueOf(&_struct)}, ""},
	}

	for _, tt := range tests {
		if t != nil {
			t.Run(tt.name, func(t *testing.T) {
				if got := Format(tt.args.value); got != tt.want {
					_assert(t, got, tt.want)
				}
			})
		}

		if b != nil {
			b.Run(tt.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					Format(tt.args.value)
				}
			})
		}
	}
}
