package database

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

type Source interface {
	storage.Source
	GetDriver() simple.Driver
}
