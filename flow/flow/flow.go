package flow

import (
	"errors"
	"fmt"
	"time"

	"github.com/auho/go-toolkit/console/output"
	"github.com/auho/go-toolkit/flow/action"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/time/timing"
)

type Option[E storage.Entry] func(*Flow[E])

func WithSource[E storage.Entry](se storage.Sourceor[E]) Option[E] {
	return func(s *Flow[E]) {
		s.source = se
	}
}

func WithActor[E storage.Entry](actor action.Actor[E]) Option[E] {
	return func(s *Flow[E]) {
		s.actions = append(s.actions, actor)
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
		output.WithContent(func() ([]string, error) {
			return f.state(), nil
		}),
	)

	err := f.source.Scan()
	if err != nil {
		return err
	}

	err = f.actionsPrepare()
	if err != nil {
		return fmt.Errorf("actioins prepare error;%w", err)
	}

	f.summary()

	err = f.actionsRun()
	if err != nil {
		return fmt.Errorf("actioins run error;%w", err)
	}

	f.refreshOutput.Start()

	f.transport()

	return f.finish()
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
					newItems := f.source.Copy(items)
					a.Receive(newItems)
				} else {
					a.Receive(items)
				}
			}
		}

		f.actionsDone()
	}()
}

func (f *Flow[E]) finish() error {
	err := f.actionsFinish()
	f.refreshOutput.Stop()
	f.actionsOutput()

	if err != nil {
		return fmt.Errorf("actions finish error; %w", err)
	} else {
		return nil
	}
}

func (f *Flow[E]) summary() {
	sss := f.source.Summary()
	sss = append(sss, "Tasks: ")
	for _, a := range f.actions {
		sss = append(sss, "  "+a.Summary())
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
			sss = append(sss, "  "+_s)
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

func (f *Flow[E]) actionsRun() error {
	var err error
	for _, a := range f.actions {
		aErr := a.Run()
		if aErr != nil {
			if err != nil {
				err = fmt.Errorf("%w, %w", err, aErr)
			} else {
				err = fmt.Errorf("run error; %w", aErr)
			}
		}
	}

	return err
}

func (f *Flow[E]) actionsPrepare() error {
	var err error
	for _, a := range f.actions {
		aErr := a.Prepare()
		if aErr != nil {
			if err != nil {
				err = fmt.Errorf("%w, %w", err, aErr)
			} else {
				err = fmt.Errorf("prepare error; %w", aErr)
			}
		}
	}

	return err
}

func (f *Flow[E]) actionsFinish() error {
	var err error
	for _, a := range f.actions {
		aErr := a.Finish()
		if aErr != nil {
			if err != nil {
				err = fmt.Errorf("%w, %w", err, aErr)
			} else {
				err = fmt.Errorf("finish error; %w", aErr)
			}
		}
	}

	return err
}

func (f *Flow[E]) actionsDone() {
	for _, a := range f.actions {
		a.Done()
	}
}
