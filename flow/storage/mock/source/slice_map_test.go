package source

import (
	"testing"

	"github.com/auho/go-toolkit/flow/storage"
)

func TestSliceMap(t *testing.T) {
	_testMock[storage.MapEntry](t, func(config Config) *mock[storage.MapEntry] {
		return NewSliceMap(config)
	})
}
