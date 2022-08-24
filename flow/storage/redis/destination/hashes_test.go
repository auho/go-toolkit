package destination

import (
	"strconv"
	"testing"

	"github.com/auho/go-toolkit/flow/storage"
)

var _hashesKey = "test:hashes"

func _buildHashesData(k *key[storage.MapEntry]) int64 {
	amount := _randAmount()
	for i := 0; i < amount; i += 100 {
		items := make(storage.MapEntries, 0, 100)
		for j := 0; j < 100; j++ {
			a := i*100 + j
			items = append(items, map[string]interface{}{strconv.Itoa(a): a})
		}

		k.Receive(items)
	}

	return int64(amount)
}

func TestHashes(t *testing.T) {
	_testKey[storage.MapEntry](t, _hashesKey, NewHashes, _buildHashesData)
}
