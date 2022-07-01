package database

import (
	"fmt"
	"github.com/auho/go-toolkit/flow/storage"
)

type State struct {
	storage.PageState
}

func NewState() *State {
	return &State{}
}

func (s *State) State() string {
	return fmt.Sprintf("Concurrency: %d, Amount: %d/%d, Page: %d/%d(%d)", s.Concurrency, s.Amount, s.Total, s.Page, s.TotalPage, s.PageSize)
}
