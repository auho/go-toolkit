package tests

import (
	"testing"
)

func Test_slice(t *testing.T) {
	_testSlice(t, nil)
}

func Benchmark_slice(b *testing.B) {
	_testSlice(nil, b)
}

func _testSlice(t *testing.T, b *testing.B) {
	type args struct {
		value interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"bool", args{[]bool{false, false, true, true}}, ""},
		{"int", args{[]int{1, 2, 3, 4, 5}}, ""},
		{"string", args{[]string{"s1", "s2", "s3", "s4", "s5"}}, ""},
		{"struct", args{[]structType{_struct}}, ""},
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
