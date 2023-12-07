package diff

import (
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Source struct {
	Name string
	DB   *simpleDb.SimpleDB
}
