package source

import (
	"context"
	"strconv"
	"testing"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/redis/client"
	"github.com/go-redis/redis/v8"
)

var _hashesKey = "test:hashes"

func _buildHashesData(t *testing.T) {
	ctx := context.Background()
	c := redis.NewClient(&_redisOptions)
	c.Del(ctx, _hashesKey)

	amount := _randAmount()
	pipe := c.Pipeline()
	for i := 0; i < amount; i++ {
		pipe.HSet(ctx, _hashesKey, strconv.Itoa(i), i)
		if i%99 == 0 {
			_, err := pipe.Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			pipe = c.Pipeline()
		}
	}
}

func TestNewHashes(t *testing.T) {
	_buildHashesData(t)

	_testKey[storage.MapOfStringsEntry](
		t,
		_hashesKey,
		NewHashes,
		func(ctx context.Context, c *client.Redis) (int64, error) {
			return c.HLen(ctx, _hashesKey).Result()
		})
}
