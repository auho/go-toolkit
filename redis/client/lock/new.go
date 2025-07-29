package lock

import (
	"github.com/bsm/redislock"
)

func NewRedisLocker(client redislock.RedisClient) (*RedisLocker, error) {
	r := &RedisLocker{}
	r.locker = redislock.New(client)

	return r, nil
}
