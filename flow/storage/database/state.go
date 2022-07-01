package database

import (
	"fmt"
	"github.com/auho/go-toolkit/time/timing"
)

type State struct {
	Concurrency int
	Page        int64
	PageSize    int64
	TotalPage   int64
	Total       int64
	Amount      int64
	Title       string
	Status      string
	Duration    timing.Duration
}

func NewState() *State {
	return &State{}
}

func (s *State) State() string {
	return fmt.Sprintf("Concurrency: %d, Amount: %d/%d, Page: %d/%d(%d)", s.Concurrency, s.Amount, s.Total, s.Page, s.TotalPage, s.PageSize)
}
