package storage

type Destinationer[E Entry] interface {
	Accept() error
	Receive([]E)
	Done()
	Finish()
	Summary() []string
	State() []string
}
