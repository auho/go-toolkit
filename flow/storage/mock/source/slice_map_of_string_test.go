package source

import (
	"testing"

	"github.com/auho/go-toolkit/flow/storage"
)

func TestSliceMapString(t *testing.T) {
	_testMock[storage.MapOfStringsEntry](t, func(config Config) *mock[storage.MapOfStringsEntry] {
		return NewSliceMapOfString(config)
	})
}
