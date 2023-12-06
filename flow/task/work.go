package task

import "github.com/auho/go-toolkit/flow/storage"

type Work[E storage.Entry] interface {
	Tasker[E]

	// Do need to be implemented
	// effected
	Do([]E) int
}
