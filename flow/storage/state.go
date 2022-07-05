package storage

import (
	"fmt"
	"github.com/auho/go-toolkit/time/timing"
)

type Stater interface {
	GetStatus() string
	Overview() string
}

type State struct {
	Concurrency int
	Amount      int64
	Title       string
	Status      string
	Duration    timing.Duration
}

func (s *State) GetStatus() string {
	return s.Status
}

type PageState struct {
	State
	Page      int64
	PageSize  int64
	TotalPage int64
	Total     int64
}

func (p *PageState) Overview() string {
	return fmt.Sprintf("Concurrency: %d, Amount: %d/%d, Page: %d/%d(%d, duration: %s)",
		p.Concurrency,
		p.Amount,
		p.Total,
		p.Page,
		p.TotalPage,
		p.PageSize,
		p.Duration.StringStartToStop())
}
