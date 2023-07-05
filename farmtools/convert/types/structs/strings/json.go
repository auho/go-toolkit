package strings

import (
	"encoding/json"
	"github.com/auho/go-toolkit/farmtools/convert/types/structs/maps"
)

// JsonToStruct convert string to any struct
// s any must be a pointer
func JsonToStruct(s any, from string) (err error) {
	var _m map[string]any
	err = json.Unmarshal([]byte(from), &_m)
	if err != nil {
		return err
	}

	return maps.MapStringAnyToStruct(s, _m)
}
