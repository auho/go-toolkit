package redis

import (
	"github.com/auho/go-toolkit/redis/client"
)

type Rediser interface {
	GetClient() *client.Redis
}
