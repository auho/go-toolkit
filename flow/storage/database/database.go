package database

import (
	"github.com/auho/go-simple-db/simple"
)

type Databaseor interface {
	GetDriver() simple.Driver
}
