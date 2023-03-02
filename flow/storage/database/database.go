package database

import (
	"gorm.io/gorm"
)

type BuildDb func() (*DB, error)

type DB struct {
	*gorm.DB
}

func (d *DB) Ping() error {
	sqlDb, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sqlDb.Ping()
}

func NewDB(d gorm.Dialector, c *gorm.Config) (*DB, error) {
	if c == nil {
		c = &gorm.Config{}
	}

	db, err := gorm.Open(d, c)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil

}

func (d *DB) Close() error {
	return nil
}

type Databaseor interface {
	DB() *DB
}
