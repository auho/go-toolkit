package flow

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/auho/go-toolkit/console/output"
	"github.com/auho/go-toolkit/flow/action"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
	"github.com/auho/go-toolkit/time/timing"
)

type Option[E storage.Entry] func(flow *Flow[E])

type Options[E storage.Entry] []Option[E]

func WithSource[E storage.Entry](sf storage.Sourceor[E]) Option[E] {
	return func(f *Flow[E]) {
		f.source = sf
	}
}

func WithTasker[E storage.Entry](t task.Tasker[E]) Option[E] {
	return func(f *Flow[E]) {
		f.actioners = append(f.actioners, action.NewAction(action.WithTasker(t)))
	}
}

func WithStateTickerDuration[E storage.Entry](d time.Duration) Option[E] {
	return func(f *Flow[E]) {
		f.stateTicker = time.NewTicker(d)
	}
}

type Flow[E storage.Entry] struct {
	source        storage.Sourceor[E]
	stateTicker   *time.Ticker
	refreshOutput *output.Refresh
	actioners     []action.Actioner[E]
}

func RunFlow[E storage.Entry](opts ...Option[E]) error {
	d := timing.NewDuration()
	d.Start()

	f := &Flow[E]{}
	for _, o := range opts {
		o(f)
	}

	err := f.check()
	if err != nil {
		return err
	}

	err = f.run()
	if err != nil {
		return err
	}

	d.StringStartToStop()

	return nil
}

func (f *Flow[E]) check() error {
	if len(f.actioners) <= 0 {
		return errors.New("action not found")
	}

	return nil
}

func (f *Flow[E]) run() error {
	if f.stateTicker == nil {
		f.stateTicker = time.NewTicker(time.Millisecond * 200)
	}

	f.refreshOutput = output.NewRefresh()

	f.summary()

	err := f.source.Scan()
	if err != nil {
		return err
	}

	f.actionerPerpare()
	f.actionerDo()

	f.refreshOutput.Start()

	go func() {
		for range f.stateTicker.C {
			f.state()
		}
	}()

	f.transport()
	f.finish()

	return nil
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

func (f *Flow[E]) finish() {
	f.actionerFinish()
	f.stateTicker.Stop()
	f.refreshOutput.Stop()
	f.actionerOutput()
}

func (f *Flow[E]) summary() {
	sss := f.source.Summary()
	sss = append(sss, "Tasks: ")
	for _, a := range f.actioners {
		sss = append(sss, a.Summary()...)
	}

	for _, s := range sss {
		fmt.Println(s)
	}
}

func (f *Flow[E]) state() {
	sss := f.source.State()

	for _, a := range f.actioners {
		sss = append(sss, a.State()...)
	}

	f.refreshOutput.CoverAll(sss)
}

func (f *Flow[E]) actionerOutput() {
	fmt.Println("\nOutput: ")

	for _, a := range f.actioners {
		for _, s := range a.Output() {
			fmt.Println(s)
		}

		fmt.Println()
	}
}

func (f *Flow[E]) actionerDo() {
	for _, a := range f.actioners {
		a.Do()
	}
}

func (f *Flow[E]) actionerPerpare() {
	for _, a := range f.actioners {
		err := a.Prepare()
		if err != nil {
			log.Fatalln(err)
		}
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
