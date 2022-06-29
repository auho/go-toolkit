package storage

import (
	"errors"
)

var EOF = errors.New("EOF")

type Source interface {
	Scan()
	ReceiveChan() <-chan []map[string]interface{}
	Next() ([]map[string]interface{}, bool)
	Summary() string
	State() string
}
