package singleton

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
)

var _ Actor[string] = (*Action[string])(nil)

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

type Option[E storage.Entry] func(singleton *Action[E])

func WithTasker[E storage.Entry](s task.Singleton[E]) Option[E] {
	return func(a *Action[E]) {
		a.task = s
	}
}

type Action[E storage.Entry] struct {
	total     int64
	amount    int64
	task      task.Singleton[E]
	taskWg    sync.WaitGroup
	itemsChan chan []E
}

func NewAction[E storage.Entry](opts ...Option[E]) *Action[E] {
	a := &Action[E]{}

	for _, o := range opts {
		o(a)
	}

	a.itemsChan = make(chan []E, a.task.Concurrency())

	return a
}

func (a *Action[E]) Receive(msi []E) {
	a.itemsChan <- msi
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

func (a *Action[E]) Run() {
	for i := 0; i < a.task.Concurrency(); i++ {
		a.taskWg.Add(1)

		go func() {
			for items := range a.itemsChan {
				atomic.AddInt64(&a.total, int64(len(items)))

				newItems := make([]E, 0, len(items))
				for k := range items {
					if v, ok := a.task.Do(items[k]); ok {
						atomic.AddInt64(&a.amount, 1)
						newItems = append(newItems, v...)
					}
				}

				a.task.PostBatchDo(newItems)
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
