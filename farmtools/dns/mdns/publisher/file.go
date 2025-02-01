package publisher

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"dns/mdns/zone"
	"github.com/fsnotify/fsnotify"
)

type File struct {
	configFile string
	conf       Config
	records    zone.Records
}

func newFile() (*File, error) {
	f := &File{
		records: make(zone.Records, 16),
	}

	return f, nil
}

func RunFile() error {
	f, err := newFile()
	if err != nil {
		return err
	}

	err = f.config()
	if err != nil {
		return err
	}

	go f.entriesChange()

	return runPublisher(&f.conf)
}

func (f *File) recordsChange() {
	f.conf.recordsChangeChan <- f.records
}

func (f *File) entriesChange() {
	f.recordsChange()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = watcher.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = watcher.Add(f.configFile)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				err = f.readRecords()
				if err != nil {
					log.Printf("watch config file: %v\n", err)
				}

				f.recordsChange()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}

			log.Printf("wathcer: %v\n", err)
		}
	}
}

func (f *File) config() error {
	f.flagParse()

	err := f.conf.config()
	if err != nil {
		return err
	}

	err = f.readRecords()
	if err != nil {
		return err
	}

	return nil
}

func (f *File) readRecords() error {
	entries, err := f.readEntries(f.configFile)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry == "" {
			continue
		}

		err = f.records.New(entry)
		if err != nil {
			log.Fatalf("entry[%s]: %v\n", entry, err)
		}
	}

	return nil
}

func (f *File) readEntries(s string) ([]string, error) {
	if s == "" {
		//return nil, fmt.Errorf("no config file found")
		s = "default.entries"
	}

	cf, err := os.Open(s)
	if err != nil {
		return nil, fmt.Errorf("read config[%s]: %v", s, err)
	}

	defer func() {
		err = cf.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var ss []string
	scanner := bufio.NewScanner(cf)
	for scanner.Scan() {
		ss = append(ss, scanner.Text())
	}

	return ss, scanner.Err()
}

func (f *File) flagParse() {
	flag.BoolVar(&f.conf.enableIpv4, "enable-ipv4", false, "enable IPv4 address")
	flag.BoolVar(&f.conf.enableIpv6, "enable-ipv6", false, "enable IPv6 address")
	flag.DurationVar(&f.conf.broadcastInterval, "t", 0, "broadcast interval duration")

	flag.StringVar(&f.configFile, "conf", "", "config file path")
	flag.StringVar(&f.configFile, "c", "", "config file path")

	flag.Parse()
}
