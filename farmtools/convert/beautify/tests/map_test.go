package tests

import (
	"testing"
)

func Test_map(t *testing.T) {
	_testMap(t, nil)
}

func Benchmark_map(b *testing.B) {
	_testMap(nil, b)
}

func _testMap(t *testing.T, b *testing.B) {
	type args struct {
		value interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"map bool empty", args{_mapBoolEmpty}, ""},
		{"map bool", args{_mapBool}, ""},
		{"map int", args{_mapInt}, ""},
		{"map string", args{_mapString}, ""},

		{"map pointer 1", args{_mapMultiPtr1}, ""},
		{"map pointer 2", args{_mapMultiPtr2}, ""},
		{"map pointer 3", args{_mapMultiPtr3}, ""},
		{"map pointer 4", args{_mapMultiPtr4}, ""},
		{"map pointer 5", args{_mapMultiPtr5}, ""},

		{"map interface nil", args{_mapAnyNil}, ""},
		{"map interface", args{_mapAny}, ""},
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
