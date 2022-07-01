package storage

import (
	"errors"
)

var EOF = errors.New("EOF")

type Source interface {
	Scan()
	ReceiveChan() <-chan []map[string]interface{}
	Summary() []string
	State() []string
}
