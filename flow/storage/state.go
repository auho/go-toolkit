package storage

import (
	"fmt"
	"sync/atomic"

	"github.com/auho/go-toolkit/time/timing"
)

const StatusConfig = "config"
const StatusAccept = "accept"
const StatusScan = "scan"
const StatusDone = "done"
const StatusFinish = "finish"

type Stater interface {
	GetStatus() string
	Overview() string
}

type baseState struct {
	Concurrency int
	Title       string
	amount      int64
	duration    timing.Duration
	status      string
}

func (s *baseState) Status() string {
	return s.status
}

func (s *baseState) Amount() int64 {
	return s.amount
}

func (s *baseState) SetAmount(n int64) {
	s.amount = n
}

func (s *baseState) AddAmount(n int64) {
	atomic.AddInt64(&s.amount, n)
}

func (s *baseState) StatusCustomize(status string) {
	s.status = status
}

func (s *baseState) StatusConfig() {
	s.status = StatusConfig
}

func (s *baseState) StatusAccept() {
	s.status = StatusAccept
}

func (s *baseState) StatusScan() {
	s.status = StatusScan
}

func (s *baseState) StatusDone() {
	s.status = StatusDone
}

func (s *baseState) StatusFinish() {
	s.status = StatusFinish
}

func (s *baseState) DurationStart() {
	s.duration.Start()
}

func (s *baseState) DurationStop() {
	s.duration.Stop()
}

type State struct {
	baseState
}

func NewState() *State {
	return &State{}
}

func (s *State) Overview() string {
	return fmt.Sprintf("Concurrency: %d, amount: %d, duration: %s",
		s.Concurrency,
		s.amount,
		s.duration.StringStartToStop())
}

type TotalState struct {
	baseState
	Total int64
}

func NewTotalState() *TotalState {
	return &TotalState{}
}

func (s *TotalState) Overview() string {
	return fmt.Sprintf("Concurrentcy:%d, amount: %d/%d, duration: %s",
		s.Concurrency,
		s.amount,
		s.Total,
		s.duration.StringStartToStop())
}

type PageState struct {
	baseState
	Page      int64
	PageSize  int64
	TotalPage int64
	Total     int64
}

func NewPageState() *PageState {
	return &PageState{}
}

func (p *PageState) Overview() string {
	return fmt.Sprintf("Concurrency: %d, amount: %d/%d, Page: %d/%d(%d), duration: %s",
		p.Concurrency,
		p.amount,
		p.Total,
		p.Page,
		p.TotalPage,
		p.PageSize,
		p.duration.StringStartToStop())
}
