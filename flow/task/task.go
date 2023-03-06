package task

import (
	"fmt"
	"runtime"
	"sync/atomic"

	"github.com/auho/go-toolkit/console/output"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/time/timing"
)

type Tasker[E storage.Entry] interface {
	// Title need to be implemented
	Title() string

	// Prepare need to be implemented
	Prepare() error

	// Do need to be implemented
	Do([]E)

	// AfterDo need to be implemented
	AfterDo()

	HasBeenInit() bool
	Concurrency() int
	State() []string
	Output() []string
}

type WithTaskOption func(*Task)

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
	log           *output.MultilineText
	hasBeenInit   bool
}

func (t *Task) Init(options ...WithTaskOption) {
	for _, o := range options {
		o(t)
	}

	if t.concurrency <= 0 {
		t.concurrency = runtime.NumCPU()
	}

	if !t.HasBeenInit() {
		t.duration = timing.NewDuration()
		t.state = output.NewMultilineText()
		t.output = output.NewMultilineText()
		t.log = output.NewMultilineText()
	}

	t.hasBeenInit = true
}

func (t *Task) HasBeenInit() bool {
	return t.hasBeenInit
}

func (t *Task) Concurrency() int {
	return t.concurrency
}

func (t *Task) SetState(line int, s string) {
	t.state.Print(line, s)
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

func (t *Task) Log() []string {
	return t.log.Content()
}

func (t *Task) Printf(format string, a ...interface{}) {
	t.output.PrintNext(fmt.Sprintf(format, a...))
}

func (t *Task) Println(a ...interface{}) {
	t.output.PrintNext(fmt.Sprint(a...))
}

func (t *Task) Logf(format string, a ...interface{}) {
	t.log.PrintNext(fmt.Sprintf(format, a...))
}

func (t *Task) Logln(a ...interface{}) {
	t.log.PrintNext(fmt.Sprint(a...))
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
