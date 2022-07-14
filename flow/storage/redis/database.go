package redis

import (
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/redis/client"
)

type Source interface {
	storage.Sourceor
	GetClient() *client.Redis
}
