package task

import "github.com/auho/go-toolkit/flow/storage"

type Singleton[E storage.Entry] interface {
	Tasker[E]
	PostBatchDo([]E)
}
