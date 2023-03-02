package database

import (
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

type BuildDb func() (*DB, error)

type DB struct {
	*goSimpleDb.SimpleDB
}

func NewDB(fn func() (*goSimpleDb.SimpleDB, error)) (*DB, error) {
	sd, err := fn()
	if err != nil {
		return nil, err
	}

	return &DB{SimpleDB: sd}, nil
}

type Driver interface {
	DB() *DB
}
