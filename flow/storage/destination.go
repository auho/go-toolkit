package storage

type Destination[E Entries] interface {
	Accept() error
	Receive(E)
	Done()
	Finish()
	Summary() []string
	State() []string
}
