package action

import (
	"github.com/auho/go-toolkit/flow/task"
	"sync"
)

type Interfaceor interface {
	Prepare()
	Do()
	Finish()
	Done()
	Receive([]map[string]interface{})
}

func WithTasker(t task.Tasker) func(i *Interface) {
	return func(i *Interface) {
		i.tasker = t
	}
}

type Interface struct {
	wg        sync.WaitGroup
	itemsChan chan []map[string]interface{}
	tasker    task.Tasker
	doFunc    func([]map[string]interface{})
}

func NewInterface(options ...func(i *Interface)) *Interface {
	i := &Interface{}

	for _, o := range options {
		o(i)
	}

	return i
}

func (i *Interface) Prepare() {
	i.itemsChan = make(chan []map[string]interface{}, 0)
}

func (i *Interface) Do() {
	concurrency := i.tasker.Concurrency()
	if concurrency <= 0 {
		concurrency = 1
	}

	for j := 0; j < concurrency; j++ {
		i.wg.Add(1)
		go func() {
			for items := range i.itemsChan {
				i.tasker.Do(items)
			}

			i.wg.Done()
		}()
	}
}

func (i *Interface) Receive(items []map[string]interface{}) {
	i.itemsChan <- items
}

func (i *Interface) Finish() {
	close(i.itemsChan)
}

func (i *Interface) Done() {
	i.wg.Wait()
}
