package task

import (
	"fmt"
	"github.com/auho/go-toolkit/flow/storage"
	"sync/atomic"
)

type Tasker interface {
	Concurrency() int
	SetSource(storage.Source)
	Do([]map[string]interface{})
	AfterDone()
	State() []string
	Output() []string
}

func WithConcurrency(c int) func(i *Task) {
	return func(i *Task) {
		i.concurrency = c
	}
}

type Task struct {
	Source        storage.Source
	concurrency   int
	amount        int64
	failureAmount int64
	state         []string
	output        []string
}

func NewTask(options ...func(*Task)) *Task {
	i := &Task{}

	for _, o := range options {
		o(i)
	}

	return i
}

func (t *Task) Concurrency() int {
	return t.concurrency
}

func (t *Task) SetSource(s storage.Source) {
	t.Source = s
}

func (t *Task) SetState(line int, s string) {
	stateLen := len(t.state)
	if line > stateLen {
		for j := 0; j < line-stateLen; j++ {
			t.state = append(t.state, "")
		}
	}

	t.state[line-1] = s
}

func (t *Task) State() []string {
	return t.state
}

func (t *Task) Output() []string {
	return t.output
}

func (t *Task) Printf(format string, a ...interface{}) {
	t.output = append(t.output, fmt.Sprintf(format, a...))
}

func (t *Task) Println(a ...interface{}) {
	t.output = append(t.output, fmt.Sprint(a...))
}

func (t *Task) AddAmount(a int64) {
	atomic.AddInt64(&t.amount, a)
}

func (t *Task) AddFailureAmount(a int64) {
	atomic.AddInt64(&t.failureAmount, a)
}

func (t *Task) Amount() int64 {
	return t.amount
}

func (t *Task) FailureAmount() int64 {
	return t.failureAmount
}
