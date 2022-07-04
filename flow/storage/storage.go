package storage

import "log"

type SliceEntry []interface{}
type MapEntry map[string]interface{}

type SliceEntries []SliceEntry
type MapEntries []MapEntry

type Entries interface {
	SliceEntries | MapEntries
}

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
