package source

import (
	"context"
	"testing"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/redis/client"
	"github.com/go-redis/redis/v8"
)

var _sortedSetsKey = "test:sortedSets"

func _buildSortedSetsData(t *testing.T) {
	ctx := context.Background()
	c := redis.NewClient(&_redisOptions)
	c.Del(ctx, _sortedSetsKey)

	amount := _randAmount()
	pipe := c.Pipeline()
	for i := 0; i < amount; i++ {
		pipe.ZAdd(ctx, _sortedSetsKey, &redis.Z{
			Score:  float64(i),
			Member: i,
		})
		if i%99 == 0 {
			_, err := pipe.Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			pipe = c.Pipeline()
		}
	}
}

func TestNewSortedSets(t *testing.T) {
	_buildSortedSetsData(t)

	_testKey[storage.MapOfStringsEntry](
		t,
		_sortedSetsKey,
		NewSortedSets,
		func(ctx context.Context, c *client.Redis) (int64, error) {
			return c.ZCard(ctx, _sortedSetsKey).Result()
		})
}
