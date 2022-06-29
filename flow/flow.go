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

func WithSource(sf storage.Source) func(i *Flow) {
	return func(f *Flow) {
		f.source = sf
	}
}

func WithTasker(t task.Tasker) func(i *Flow) {
	return func(i *Flow) {
		i.taskers = append(i.taskers, t)
	}
}

type Flow struct {
	source        storage.Source
	stateTicker   *time.Ticker
	refreshOutput *output.Refresh
	taskers       []task.Tasker
	actioners     []action.Actioner
}

func RunFlow(options ...func(*Flow)) {
	d := timing.NewDuration()
	d.Start()

	i := &Flow{}
	for _, o := range options {
		o(i)
	}

	i.run()

	d.StringStartToNowSeconds()
}

func (f *Flow) run() {
	f.process()
	f.transport()
	f.done()
}

func (f *Flow) process() {
	f.stateTicker = time.NewTicker(time.Millisecond * 100)
	f.refreshOutput = output.NewRefresh()
	for _, t := range f.taskers {
		f.actioners = append(f.actioners, action.NewAction(action.WithTasker(t)))
	}

	f.source.Scan()
	for _, a := range f.actioners {
		a.Do()
	}

	f.refreshOutput.Start()

	go func() {
		for range f.stateTicker.C {
			f.state()
		}
	}()
}

func (f *Flow) transport() {
	needCopy := false
	if len(f.actioners) > 1 {
		needCopy = true
	}

	go func() {
		for {
			items, ok := f.source.Next()
			if !ok {
				break
			}

			for _, a := range f.actioners {
				if needCopy {
					newItems := f.copySliceMapInterface(items)
					a.Receive(newItems)
				} else {
					a.Receive(items)
				}
			}
		}

		f.actionerFinish()
	}()
}

func (f *Flow) done() {
	f.actionerDone()
	f.stateTicker.Stop()
	f.refreshOutput.Stop()
	f.output()
}

func (f *Flow) state() {
	sss := []string{f.source.State()}

	for _, a := range f.actioners {
		sss = append(sss, a.State()...)
	}

	f.refreshOutput.CoverAll(sss)
}

func (f *Flow) output() {
	for _, a := range f.actioners {
		for _, s := range a.Output() {
			fmt.Println(s)
		}

		fmt.Println()
	}
}

func (f *Flow) actionerDone() {
	for _, a := range f.actioners {
		a.Done()
	}
}

func (f *Flow) actionerFinish() {
	for _, a := range f.actioners {
		a.Finish()
	}
}

func (f *Flow) copySliceMapInterface(sm []map[string]interface{}) []map[string]interface{} {
	newSM := make([]map[string]interface{}, len(sm))
	for k, m := range sm {
		newSM[k] = f.copyMapInterface(m)
	}

	return newSM
}

func (f *Flow) copyMapInterface(m map[string]interface{}) map[string]interface{} {
	newM := make(map[string]interface{}, len(m))
	for k, v := range m {
		newM[k] = v
	}

	return newM
}
