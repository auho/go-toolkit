package strings

import (
	"github.com/auho/go-toolkit/farmtools/convert/types/ints"
	strings2 "github.com/auho/go-toolkit/farmtools/convert/types/strings"
	"strings"
)

func SpiltToSliceString(from string, sep string) ([]string, error) {
	return spiltToSliceAny(from, sep, func(v any) (string, error) {
		return strings2.FromAny(v)
	})
}

func SpiltToSliceInt(from string, sep string) ([]int, error) {
	return spiltToSliceAny(from, sep, func(v any) (int, error) {
		return ints.FromAny(v)
	})
}

func spiltToSliceAny[E string | int](from string, sep string, valueHandler func(v any) (E, error)) ([]E, error) {
	ss := strings.Split(from, sep)

	var ns []E
	for _, s := range ss {
		_s, err := valueHandler(s)
		if err != nil {
			return nil, err
		}

		ns = append(ns, _s)
	}

	return ns, nil
}
