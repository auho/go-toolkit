package source

import (
	"github.com/go-redis/redis/v8"
)

type Config struct {
	Concurrency int
	PageSize    int64
	Key         string
	Options     *redis.Options
}
