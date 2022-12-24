package task

import (
	"fmt"
	"sync/atomic"

	"github.com/auho/go-toolkit/console/output"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/time/timing"
)

type Tasker[E storage.Entry] interface {
	// Title need to be implemented
	Title() string

	// Prepare need to be implemented
	Prepare()

	// Do need to be implemented
	Do([]E)

	// AfterDo need to be implemented
	AfterDo()

	Concurrency() int
	State() []string
	Output() []string
}

func WithConcurrency(c int) func(i *Task) {
	return func(i *Task) {
		i.concurrency = c
	}
}

type Task struct {
	concurrency   int
	amount        int64
	failureAmount int64
	duration      *timing.Duration
	state         *output.MultilineText
	output        *output.MultilineText
}

func (t *Task) Init(options ...func(*Task)) {
	for _, o := range options {
		o(t)
	}

	t.duration = timing.NewDuration()
	t.state = output.NewMultilineText()
	t.output = output.NewMultilineText()

	if t.concurrency <= 0 {
		t.concurrency = 2
	}
}

func (t *Task) Concurrency() int {
	return t.concurrency
}

func (t *Task) State() []string {
	sss := t.state.Content()
	if len(sss) <= 0 {
		sss = append(sss, fmt.Sprintf("Amount: %d", t.amount))
	}

	return sss
}

func (t *Task) Output() []string {
	return t.output.Content()
}

func (t *Task) SetState(line int, s string) {
	t.state.Print(line, s)
}

func (t *Task) Printf(format string, a ...interface{}) {
	t.output.PrintNext(fmt.Sprintf(format, a...))
}

func (t *Task) Println(a ...interface{}) {
	t.output.PrintNext(fmt.Sprint(a...))
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
