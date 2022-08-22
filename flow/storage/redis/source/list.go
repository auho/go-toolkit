package source

import (
	"context"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/redis/client"
)

var _ keyer[string] = (*listKey)(nil)

type listKey struct {
	storage.Storage
	amount int64
}

func (l listKey) keyType() keyType {
	return keyTypeList
}

func (l listKey) len(c *client.Redis, key string) (int64, error) {
	return c.LLen(context.Background(), key).Result()
}

func (l listKey) scan(entriesChan chan<- []string, c *client.Redis, key string, amount int64, pageSize int64) {
	start := int64(0)
	stop := start + pageSize - 1

	for {
		items, err := c.LRange(context.Background(), key, start, stop).Result()
		if err != nil {
			l.LogFatal(err)
		}

		if len(items) <= 0 {
			break
		}

		entriesChan <- items

		start = stop + 1
		stop = start + pageSize - 1

		if start >= amount {
			break
		}
	}
}

func (l listKey) duplicate(items []string) []string {
	newItems := make([]string, 0, len(items))
	_ = copy(newItems, items)

	return newItems
}

func (l listKey) stateAmount() int64 {
	return l.amount
}
