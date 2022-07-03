package action

import (
	"sync"

	"github.com/auho/go-toolkit/flow/task"
)

type Actioner interface {
	Receive([]map[string]interface{})
	Do()
	Finish()
	Done()
	State() []string
	Output() []string
}

func WithTasker(t task.Tasker) func(*Action) {
	return func(a *Action) {
		a.tasker = t
	}
}

type Action struct {
	tasker    task.Tasker
	taskerWg  sync.WaitGroup
	itemsChan chan []map[string]interface{}
}

func NewAction(options ...func(i *Action)) *Action {
	a := &Action{}

	for _, o := range options {
		o(a)
	}

	a.itemsChan = make(chan []map[string]interface{}, a.tasker.Concurrency())

	return a
}

func (a *Action) Receive(msi []map[string]interface{}) {
	a.itemsChan <- msi
}

func (a *Action) Do() {
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

func (a *Action) Done() {
	close(a.itemsChan)
}

func (a *Action) Finish() {
	a.taskerWg.Wait()

	a.tasker.AfterDo()
}

func (a *Action) State() []string {
	return a.tasker.State()
}

func (a *Action) Output() []string {
	return a.tasker.Output()
}
