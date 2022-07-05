package source

import (
	"fmt"
	"github.com/auho/go-toolkit/flow/storage"
	"math"
	"sync/atomic"
	"time"
)

var _ storage.Source = (*Source)(nil)

func WithPageSize(i int64) func(m *Source) {
	return func(m *Source) {
		m.pageSize = i
	}
}

func WithTotal(t int64) func(m *Source) {
	return func(m *Source) {
		m.total = t
	}
}

func WithIdName(s string) func(m *Source) {
	return func(m *Source) {
		m.idName = s
	}
}

type Source struct {
	storage.Storage
	id        int64
	total     int64 // 最大数量(总数)
	page      int64
	pageSize  int64
	totalPage int64
	amount    int64
	idName    string
	itemChan  chan []map[string]interface{}
}

func NewSource(options ...func(*Source)) *Source {
	s := &Source{}

	for _, o := range options {
		o(s)
	}

	if s.total <= 0 {
		s.total = 1e2
	}

	if s.pageSize <= 0 {
		s.total = 1e1
	}

	if s.idName == "" {
		s.idName = "id"
	}

	s.totalPage = int64(math.Ceil(float64(s.total) / float64(s.pageSize)))

	return s
}

func (s *Source) Scan() {
	s.itemChan = make(chan []map[string]interface{})

	go func() {
		for i := int64(0); i < s.total; i += s.pageSize {
			size := s.pageSize
			if i+s.pageSize > s.total {
				size = s.total - i
			}

			items := make([]map[string]interface{}, size, size)
			for j := int64(0); j < size; j++ {
				item := make(map[string]interface{})
				item[s.idName] = time.Now().Unix()*1e8 + atomic.AddInt64(&s.id, 1)
				items[j] = item
			}

			s.itemChan <- items

			atomic.AddInt64(&s.page, 1)
			atomic.AddInt64(&s.amount, int64(len(items)))
		}

		close(s.itemChan)
	}()
}

func (s *Source) ReceiveChan() <-chan []map[string]interface{} {
	return s.itemChan
}

func (s *Source) Summary() []string {
	return []string{fmt.Sprintf("%s: max: %d, pageSize: %d", s.Title(), s.total, s.pageSize)}
}

func (s *Source) State() []string {
	return []string{fmt.Sprintf("amount: %d/%d, page: %d/%d(%d)", s.amount, s.total, s.page, s.totalPage, s.pageSize)}
}

func (s *Source) Title() string {
	return "Source:source"
}
