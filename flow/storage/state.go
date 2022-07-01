package storage

import "github.com/auho/go-toolkit/time/timing"

type State struct {
	Concurrency int
	Amount      int64
	Title       string
	Status      string
	Duration    timing.Duration
}

type PageState struct {
	State
	Page      int64
	PageSize  int64
	TotalPage int64
	Total     int64
}
