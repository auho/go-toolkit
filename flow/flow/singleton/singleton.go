package singleton

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/auho/go-toolkit/console/output"
	"github.com/auho/go-toolkit/flow/action"
	"github.com/auho/go-toolkit/flow/action/singleton"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
	"github.com/auho/go-toolkit/time/timing"
)

type Option[E storage.Entry] func(*Flow[E])

func WithSource[E storage.Entry](se storage.Sourceor[E]) Option[E] {
	return func(s *Flow[E]) {
		s.source = se
	}
}

func WithTasker[E storage.Entry](t task.Singleton[E]) Option[E] {
	return func(s *Flow[E]) {
		s.actions = append(s.actions, singleton.NewAction(singleton.WithSingleton(t)))
	}
}

func WithStateInterval[E storage.Entry](d time.Duration) Option[E] {
	return func(f *Flow[E]) {
		f.stateInterval = d
	}
}

type Flow[E storage.Entry] struct {
	source        storage.Sourceor[E]
	refreshOutput *output.Refresh
	actions       []action.Actor[E]
	stateInterval time.Duration
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
	if len(f.actions) <= 0 {
		return errors.New("action not found")
	}

	return nil
}

func (f *Flow[E]) run() error {
	f.refreshOutput = output.NewRefresh(
		output.WithInterval(f.stateInterval),
		output.WithContentGetter(func() []string {
			return f.state()
		}),
	)

	f.summary()

	err := f.source.Scan()
	if err != nil {
		return err
	}

	f.actionsPrepare()
	f.actionsRun()

	f.refreshOutput.Start()

	f.transport()
	f.finish()

	return nil
}

func (f *Flow[E]) transport() {
	needCopy := false
	if len(f.actions) > 1 {
		needCopy = true
	}

	go func() {
		for {
			items, ok := <-f.source.ReceiveChan()
			if !ok {
				break
			}

			for _, a := range f.actions {
				if needCopy {
					newItems := f.source.Duplicate(items)
					a.Receive(newItems)
				} else {
					a.Receive(items)
				}
			}
		}

		f.actionsDone()
	}()
}

func (f *Flow[E]) finish() {
	f.actionsFinish()
	f.refreshOutput.Stop()
	f.actionsOutput()
}

func (f *Flow[E]) summary() {
	sss := f.source.Summary()
	sss = append(sss, "Tasks: ")
	for _, a := range f.actions {
		sss = append(sss, "\t"+a.Summary())
	}

	for _, s := range sss {
		fmt.Println(s)
	}

	fmt.Println("")
}

func (f *Flow[E]) state() []string {
	sss := f.source.State()

	for _, a := range f.actions {
		sss = append(sss, a.Summary())
		for _, _s := range a.State() {
			sss = append(sss, "\t"+_s)
		}
	}

	return sss
}

func (f *Flow[E]) actionsOutput() {
	fmt.Println("\nOutput: ")

	for _, a := range f.actions {
		for _, s := range a.Output() {
			fmt.Println(s)
		}

		fmt.Println()
	}
}

func (f *Flow[E]) actionsRun() {
	for _, a := range f.actions {
		a.Run()
	}
}

func (f *Flow[E]) actionsPrepare() {
	for _, a := range f.actions {
		err := a.Prepare()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (f *Flow[E]) actionsFinish() {
	for _, a := range f.actions {
		a.Finish()
	}
}

func (f *Flow[E]) actionsDone() {
	for _, a := range f.actions {
		a.Done()
	}
}
