package destination

import (
	"context"
	"sync/atomic"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/redis"
	"github.com/auho/go-toolkit/redis/client"
)

var _ keyer[storage.MapEntry] = (*hashes)(nil)

type hashes struct {
	redis.Hashes
	amount int64
}

func (h *hashes) stateAmount() int64 {
	return h.amount
}

func NewHashes(config Config) (*key[storage.MapEntry], error) {
	return newKey[storage.MapEntry](config, &hashes{})
}

func (h *hashes) accept(itemsChan <-chan []storage.MapEntry, c *client.Redis, key string, pageSize int64) {
	ctx := context.Background()
	for items := range itemsChan {
		l := len(items)
		for i := 0; i < l; i += int(pageSize) {
			pipe := c.Pipeline()

			end := i + int(pageSize)
			entries := items[i:end]
			for _, entry := range entries {
				for k, v := range entry {
					pipe.HMSet(ctx, key, k, v)
				}
			}

			_, err := pipe.Exec(ctx)
			if err != nil {
				panic(err)
			}
			_ = pipe.Close()
		}

		atomic.AddInt64(&h.amount, int64(l))
	}
}
