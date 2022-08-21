package source

import (
	"context"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/tool"
	"github.com/auho/go-toolkit/redis/client"
)

var _ keyer[storage.MapOfStringsEntry] = (*sortedSetsKey)(nil)

type sortedSetsKey struct {
	storage.Storage
	amount    int64
	itemsChan chan storage.MapOfStringsEntries
}

func NewSortedSets(config Config) (*key[storage.MapOfStringsEntry], error) {
	return newKey[storage.MapOfStringsEntry](config, &sortedSetsKey{})
}

func (s *sortedSetsKey) keyType() keyType {
	return keyTypeSortedSets
}

func (s *sortedSetsKey) len(c *client.Redis, key string) (int64, error) {
	return c.ZCard(context.Background(), key).Result()
}

func (s *sortedSetsKey) scan(entriesChan chan<- storage.MapOfStringsEntries, c *client.Redis, key string, amount int64, count int64) {
	cursor := uint64(0)
	for {
		items, cursor, err := c.ZScan(context.Background(), key, cursor, "", count).Result()
		if err != nil {
			s.LogFatal(err)
		}

		entries := make(storage.MapOfStringsEntries, 0, len(items)/2)

		for i := 0; i < len(items)/2; i++ {
			entries = append(entries, storage.MapOfStringsEntry{items[i]: items[i+1]})
		}

		s.amount += int64(len(entries))
		entriesChan <- entries

		if cursor == 0 {
			break
		}

		if s.amount >= amount {
			break
		}
	}
}

func (s *sortedSetsKey) duplicate(items storage.MapOfStringsEntries) storage.MapOfStringsEntries {
	return tool.DuplicateSliceMap[tool.StringEntry](items)
}

func (s *sortedSetsKey) stateAmount() int64 {
	return s.amount
}
