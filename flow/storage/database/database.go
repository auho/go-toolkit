package database

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

type Source[E storage.Entry] interface {
	storage.Sourceor[E]
	GetDriver() simple.Driver
}
