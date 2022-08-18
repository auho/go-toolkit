package storage

import (
	"errors"
)

var EOF = errors.New("EOF")

type Sourceor[E Entry] interface {
	Scan() error
	ReceiveChan() <-chan []E
	Summary() []string
	State() []string
	Duplicate([]E) []E
}
