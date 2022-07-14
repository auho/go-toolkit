package storage

import (
	"errors"
)

var EOF = errors.New("EOF")

type Sourceor interface {
	Scan() error
	ReceiveChan() <-chan []map[string]interface{}
	Summary() []string
	State() []string
}
