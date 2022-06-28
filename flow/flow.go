package flow

import (
	"fmt"
	"github.com/auho/go-toolkit/console/output"
	"github.com/auho/go-toolkit/flow/action"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
	"github.com/auho/go-toolkit/time/timing"
	"time"
)

func WithSourceInterface(sf storage.Source) func(i *Interface) {
	return func(f *Interface) {
		f.source = sf
	}
}

func WithTasker(t task.Tasker) func(i *Interface) {
	return func(i *Interface) {
		i.taskers = append(i.taskers, t)
	}
}

type Interface struct {
	source        storage.Source
	ticker        *time.Ticker
	refreshOutput *output.Refresh
	actionors     []action.Interfaceor
	taskers       []task.Tasker
}

func RunInterface(options ...func(*Interface)) {
	d := timing.NewDuration()
	d.Start()

	i := &Interface{}
	for _, o := range options {
		o(i)
	}

	i.run()

	d.StringStartToNowSeconds()
}

func (i *Interface) run() {
	i.process()
	i.transport()
	i.done()
}

func (i *Interface) process() {
	i.ticker = time.NewTicker(time.Millisecond * 100)
	i.refreshOutput = output.NewRefresh()
	for _, t := range i.taskers {
		t.SetSource(i.source)
		i.actionors = append(i.actionors, action.NewInterface(action.WithTasker(t)))
	}

	i.source.Scan()
	for _, a := range i.actionors {
		a.Prepare()
		a.Do()
	}

	i.refreshOutput.Start()

	go func() {
		for range i.ticker.C {
			i.taskerStatus()
		}
	}()
}

func (i *Interface) transport() {
	needCopy := false
	if len(i.actionors) > 1 {
		needCopy = true
	}

	go func() {
		for {
			items, ok := i.source.Next()
			if !ok {
				break
			}

			for _, a := range i.actionors {
				if needCopy {
					newItems := i.copySliceMapInterface(items)
					a.Receive(newItems)
				} else {
					a.Receive(items)
				}
			}
		}

		i.actionorFinish()
	}()
}

func (i *Interface) done() {
	i.actionorDone()
	i.taskerAfterDone()
	i.ticker.Stop()
	i.refreshOutput.Stop()
	i.taskerOutput()
}

func (i *Interface) taskerStatus() {
	lines := 1

	i.refreshOutput.Print(lines, i.source.State())
	lines += 1

	for _, t := range i.taskers {
		for _, s := range t.State() {
			i.refreshOutput.Print(lines, s)
			lines += 1
		}

		i.refreshOutput.Print(lines, "")
		lines += 1
	}
}

func (i *Interface) taskerOutput() {
	for _, t := range i.taskers {
		for _, s := range t.Output() {
			fmt.Println(s)
		}

		fmt.Println()
	}
}

func (i *Interface) actionorFinish() {
	for _, a := range i.actionors {
		a.Finish()
	}
}

func (i *Interface) actionorDone() {
	for _, a := range i.actionors {
		a.Done()
	}
}

func (i *Interface) taskerAfterDone() {
	for _, t := range i.taskers {
		t.AfterDone()
	}
}

func (i *Interface) copySliceMapInterface(sm []map[string]interface{}) []map[string]interface{} {
	newSM := make([]map[string]interface{}, len(sm))
	for k, m := range sm {
		newSM[k] = m
	}

	return newSM
}

func (i *Interface) copyMapInterface(m map[string]interface{}) map[string]interface{} {
	newM := make(map[string]interface{}, len(m))
	for k, v := range m {
		newM[k] = v
	}

	return newM
}
