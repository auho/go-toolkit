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

func NewState() *State {
	return &State{}
}

func (s *State) GetStatus() string {
	return s.Status
}

func (s *State) Overview() string {
	return fmt.Sprintf("Concurrency: %d, Amount: %d, duration: %s",
		s.Concurrency,
		s.Amount,
		s.Duration.StringStartToStop())
}

type PageState struct {
	State
	Page      int64
	PageSize  int64
	TotalPage int64
	Total     int64
}

func NewPageState() *PageState {
	return &PageState{}
}

func (p *PageState) Overview() string {
	return fmt.Sprintf("Concurrency: %d, Amount: %d/%d, Page: %d/%d(%d), duration: %s",
		p.Concurrency,
		p.Amount,
		p.Total,
		p.Page,
		p.TotalPage,
		p.PageSize,
		p.Duration.StringStartToStop())
}
