package publisher

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type fsNotify struct {
	watcher *fsnotify.Watcher
}

func newFsNotify() (*fsNotify, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &fsNotify{watcher: w}, nil
}

func runFsNotify(we func(fsnotify.Event), ss ...string) error {
	fs, err := newFsNotify()
	if err != nil {
		return err
	}

	err = fs.add(ss...)
	if err != nil {
		return err
	}

	go fs.watch(we)

	return nil
}

func (f *fsNotify) add(ss ...string) error {
	for _, s := range ss {
		err := f.watcher.Add(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *fsNotify) watch(we func(event fsnotify.Event)) {
	defer func() {
		err := f.watcher.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}

			we(event)
		case err, ok := <-f.watcher.Errors:
			if !ok {
				return
			}

			log.Printf("fs notify: %v", err)
		}
	}
}
