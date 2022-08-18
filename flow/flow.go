package flow

import (
	"fmt"
	"time"

	"github.com/auho/go-toolkit/console/output"
	"github.com/auho/go-toolkit/flow/action"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
	"github.com/auho/go-toolkit/time/timing"
)

func WithSource[E storage.Entry](sf storage.Sourceor[E]) func(i *Flow[E]) {
	return func(f *Flow[E]) {
		f.source = sf
	}
}

func WithTasker[E storage.Entry](t task.Tasker[E]) func(*Flow[E]) {
	return func(f *Flow[E]) {
		f.actioners = append(f.actioners, action.NewAction(action.WithTasker(t)))
	}
}

type Flow[E storage.Entry] struct {
	source        storage.Sourceor[E]
	stateTicker   *time.Ticker
	refreshOutput *output.Refresh
	actioners     []action.Actioner[E]
}

func RunFlow[E storage.Entry](options ...func(*Flow[E])) {
	d := timing.NewDuration()
	d.Start()

	i := &Flow[E]{}
	for _, o := range options {
		o(i)
	}

	i.run()

	d.StringStartToStop()
}

func (f *Flow[E]) run() {
	f.process()
	f.transport()
	f.Finish()
}

func (f *Flow[E]) process() {
	f.stateTicker = time.NewTicker(time.Millisecond * 100)
	f.refreshOutput = output.NewRefresh()

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

func (f *Flow[E]) transport() {
	needCopy := false
	if len(f.actioners) > 1 {
		needCopy = true
	}

	go func() {
		for {
			items, ok := <-f.source.ReceiveChan()
			if !ok {
				break
			}

			for _, a := range f.actioners {
				if needCopy {
					newItems := f.source.Duplicate(items)
					a.Receive(newItems)
				} else {
					a.Receive(items)
				}
			}
		}

		f.actionerDone()
	}()
}

func (f *Flow[E]) Finish() {
	f.actionerFinish()
	f.stateTicker.Stop()
	f.refreshOutput.Stop()
	f.output()
}

func (f *Flow[E]) state() {
	sss := f.source.State()

	for _, a := range f.actioners {
		sss = append(sss, a.State()...)
	}

	f.refreshOutput.CoverAll(sss)
}

func (f *Flow[E]) output() {
	for _, a := range f.actioners {
		for _, s := range a.Output() {
			fmt.Println(s)
		}

		fmt.Println()
	}
}

func (f *Flow[E]) actionerFinish() {
	for _, a := range f.actioners {
		a.Finish()
	}
}

func (f *Flow[E]) actionerDone() {
	for _, a := range f.actioners {
		a.Done()
	}
}
