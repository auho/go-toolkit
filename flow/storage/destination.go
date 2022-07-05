package storage

type Destination[E Entry] interface {
	Accept() error
	Receive([]E)
	Done()
	Finish()
	Summary() []string
	State() []string
}
