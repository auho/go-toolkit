package client

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	*redis.Client
	locker *redislock.Client
}

func NewRedisClient(opt *redis.Options) (*Redis, error) {
	r := &Redis{}
	r.Client = redis.NewClient(opt)
	r.locker = redislock.New(r.Client)

	_, err := r.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return r, nil
}
