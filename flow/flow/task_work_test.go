package flow

import (
	"github.com/auho/go-toolkit/flow/task"
)

var _ task.Work[map[string]any] = (*work)(nil)

type work struct {
	task.Task
}

func (w *work) Title() string {
	return "test work"
}

func (w *work) Prepare() error {
	w.SetState(0, "prepare")
	return nil
}

func (w *work) Do(items []map[string]any) {
	for _, item := range items {
		_ = item
	}
}

func (w *work) PostDo() {
	w.SetState(0, "post do")
	w.Println("post do")
}
