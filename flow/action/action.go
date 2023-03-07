package action

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
)

var _ Actor[string] = (*Action[string])(nil)

type Moder[E storage.Entry] interface {
	Concurrency() int
	Tasker() task.Tasker[E]
	Run([]E) int // Process data
}

type Actor[E storage.Entry] interface {
	Prepare() error // preparation before processing data
	Receive([]E)    // receive data asynchronously
	Run()           // Process data
	Done()          // triggered after upstream data processing
	Finish()        // data processing completed
	Summary() string
	State() []string
	Output() []string
}

type Action[E storage.Entry] struct {
	total     int64
	amount    int64
	itemsChan chan []E
	mode      Moder[E]
	task      task.Tasker[E]
	taskWg    sync.WaitGroup
}

func NewAction[E storage.Entry](mode Moder[E]) *Action[E] {
	a := &Action[E]{}
	a.mode = mode
	a.task = a.mode.Tasker()
	a.itemsChan = make(chan []E, a.mode.Concurrency())

	return a
}

func (a *Action[E]) Prepare() error {
	if !a.task.HasBeenInit() {
		return errors.New("task of action has not been init")
	}

	err := a.task.Prepare()
	if err != nil {
		return err
	}

	return nil
}

func (a *Action[E]) Receive(items []E) {
	a.itemsChan <- items
}

func (a *Action[E]) Run() {
	for i := 0; i < a.task.Concurrency(); i++ {
		a.taskWg.Add(1)

		go func() {
			for items := range a.itemsChan {
				atomic.AddInt64(&a.total, int64(len(items)))
				amount := a.mode.Run(items)
				atomic.AddInt64(&a.amount, int64(amount))
			}

			a.taskWg.Done()
		}()
	}
}

func (a *Action[E]) Done() {
	close(a.itemsChan)
}

func (a *Action[E]) Finish() {
	a.taskWg.Wait()

	a.task.PostDo()
}

func (a *Action[E]) Summary() string {
	return a.task.Title()
}

func (a *Action[E]) State() []string {
	return append([]string{fmt.Sprintf("Total: %d, Amount %d", a.total, a.amount)}, a.task.State()...)
}

func (a *Action[E]) Output() []string {
	return a.task.Output()
}
