package storage

import "log"

type Storage struct {
}

func (s *Storage) Title() string {
	return ""
}

func (s *Storage) LogFatalWithTitle(v ...any) {
	log.Fatal(append([]interface{}{s.Title()}, v...)...)
}

func (s *Storage) LogFatal(v ...any) {
	log.Fatal(v...)
}
