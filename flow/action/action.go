package action

import (
	"sync"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
)

type Actioner[E storage.Entry] interface {
	Receive([]E)
	Do()
	Finish()
	Done()
	State() []string
	Output() []string
}

func WithTasker[E storage.Entry](t task.Tasker[E]) func(*Action[E]) {
	return func(a *Action[E]) {
		a.tasker = t
	}
}

type Action[E storage.Entry] struct {
	tasker    task.Tasker[E]
	taskerWg  sync.WaitGroup
	itemsChan chan []E
}

func NewAction[E storage.Entry](options ...func(i *Action[E])) *Action[E] {
	a := &Action[E]{}

	for _, o := range options {
		o(a)
	}

	a.itemsChan = make(chan []E, a.tasker.Concurrency())

	return a
}

func (a *Action[E]) Receive(msi []E) {
	a.itemsChan <- msi
}

func (a *Action[E]) Do() {
	for i := 0; i < a.tasker.Concurrency(); i++ {
		a.taskerWg.Add(1)

		go func() {
			for items := range a.itemsChan {
				a.tasker.Do(items)
			}

			a.taskerWg.Done()
		}()
	}
}

func (a *Action[E]) Done() {
	close(a.itemsChan)
}

func (a *Action[E]) Finish() {
	a.taskerWg.Wait()

	a.tasker.AfterDo()
}

func (a *Action[E]) State() []string {
	return a.tasker.State()
}

func (a *Action[E]) Output() []string {
	return a.tasker.Output()
}
