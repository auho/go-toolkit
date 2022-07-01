package mock

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"
)

func WithPageSize(i int64) func(m *MockSource) {
	return func(m *MockSource) {
		m.pageSize = i
	}
}

func WithTotal(t int64) func(m *MockSource) {
	return func(m *MockSource) {
		m.total = t
	}
}

func WithIdName(s string) func(m *MockSource) {
	return func(m *MockSource) {
		m.idName = s
	}
}

type MockSource struct {
	id        int64
	total     int64 // 最大数量(总数)
	page      int64
	pageSize  int64
	totalPage int64
	amount    int64
	idName    string
	itemChan  chan []map[string]interface{}
}

func NewMock(options ...func(*MockSource)) *MockSource {
	m := &MockSource{}

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

func (m *MockSource) Scan() {
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

func (m *MockSource) ReceiveChan() <-chan []map[string]interface{} {
	return m.itemChan
}

func (m *MockSource) Summary() string {
	return fmt.Sprintf("%s: max: %d, pageSize: %d", m.title(), m.total, m.pageSize)
}

func (m *MockSource) State() []string {
	return []string{fmt.Sprintf("amount: %d/%d, page: %d/%d(%d)", m.amount, m.total, m.page, m.totalPage, m.pageSize)}
}

func (m *MockSource) title() string {
	return "MockSource"
}
