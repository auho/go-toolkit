package redis

import (
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/redis/client"
)

type Source[E storage.Entry] interface {
	storage.Sourceor[E]
	GetClient() *client.Redis
}
