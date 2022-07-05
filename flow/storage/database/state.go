package database

import (
	"github.com/auho/go-toolkit/flow/storage"
)

type State struct {
	storage.PageState
}

func NewState() *State {
	return &State{}
}
