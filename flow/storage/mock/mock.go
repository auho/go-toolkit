package mock

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"
)

func WithPageSize(i int64) func(m *Mock) {
	return func(m *Mock) {
		m.pageSize = i
	}
}

func WithTotal(t int64) func(m *Mock) {
	return func(m *Mock) {
		m.total = t
	}
}

func WithIdName(s string) func(m *Mock) {
	return func(m *Mock) {
		m.idName = s
	}
}

type Mock struct {
	id        int64
	total     int64 // 最大数量(总数)
	page      int64
	pageSize  int64
	totalPage int64
	amount    int64
	idName    string
	itemChan  chan []map[string]interface{}
}

func NewMock(options ...func(*Mock)) *Mock {
	m := &Mock{}

	for _, o := range options {
		o(m)
	}

	if m.total <= 0 {
		m.total = 1e2
	}

	if m.pageSize <= 0 {
		m.total = 1e1
	}

	if m.idName == "" {
		m.idName = "id"
	}

	m.totalPage = int64(math.Ceil(float64(m.total) / float64(m.pageSize)))

	return m
}

func (m *Mock) Scan() {
	m.itemChan = make(chan []map[string]interface{})

	go func() {
		for i := int64(0); i < m.total; i += m.pageSize {
			size := m.pageSize
			if i+m.pageSize > m.total {
				size = m.total - i
			}

			items := make([]map[string]interface{}, size, size)
			for j := int64(0); j < size; j++ {
				item := make(map[string]interface{})
				item[m.idName] = time.Now().Unix()*1e8 + atomic.AddInt64(&m.id, 1)
				items[j] = item
			}

			m.itemChan <- items

			atomic.AddInt64(&m.page, 1)
			atomic.AddInt64(&m.amount, int64(len(items)))
		}

		close(m.itemChan)
	}()
}

func (m *Mock) ReceiveChan() <-chan []map[string]interface{} {
	return m.itemChan
}

func (m *Mock) Next() ([]map[string]interface{}, bool) {
	s, ok := <-m.itemChan

	return s, ok
}

func (m *Mock) Summary() string {
	return fmt.Sprintf("%s: max: %d, pageSize: %d", m.title(), m.total, m.pageSize)
}

func (m *Mock) State() string {
	return fmt.Sprintf("amount: %d/%d, page: %d/%d(%d)", m.amount, m.total, m.page, m.totalPage, m.pageSize)
}

func (m *Mock) title() string {
	return "Mock"
}
