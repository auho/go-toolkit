package database

import (
	"fmt"
	"github.com/auho/go-toolkit/time/timing"
)

type State struct {
	title     string
	page      int64
	pageSize  int64
	totalPage int64
	total     int64
	amount    int64
	duration  timing.Duration
}

func newState() *State {
	return &State{}
}

func (s *State) State() string {
	return fmt.Sprintf("amount: %d/%d, page: %d/%d(%d)", s.amount, s.total, s.page, s.totalPage, s.pageSize)
}
