package strings

import (
	"encoding/json"
	"github.com/auho/go-toolkit/farmtools/convert/types/ints"
	"github.com/auho/go-toolkit/farmtools/convert/types/strings"
)

func JsonToSlicesString(from string) ([]string, error) {
	return jsonToSliceAny(from, func(v any) (string, error) {
		return strings.FromAny(v)
	})
}

func JsonToSliceInt(from string) ([]int, error) {
	return jsonToSliceAny(from, func(v any) (int, error) {
		return ints.FromAny(v)
	})
}

func jsonToSliceAny[E string | int](from string, valueHandler func(v any) (E, error)) ([]E, error) {
	var m []any
	err := json.Unmarshal([]byte(from), &m)
	if err != nil {
		return nil, err
	}

	var nm []E
	for _, _m := range m {
		_s, err1 := valueHandler(_m)
		if err1 != nil {
			return nm, err1
		}

		nm = append(nm, _s)
	}

	return nm, nil
}
