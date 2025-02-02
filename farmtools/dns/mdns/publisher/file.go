package publisher

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"dns/mdns/zone"
	"github.com/fsnotify/fsnotify"
	"github.com/go-ini/ini"
)

const defaultMdns = "mdns.conf"
const defaultHosts = "hosts.ini"

type File struct {
	mdnsConfFile  string
	hostsConfFile string

	conf         Config
	mdnsEntries  []string
	hostsEntries []string
	records      zone.Records
}

func newFile() (*File, error) {
	f := &File{
		records: zone.NewRecords(),
	}

	return f, nil
}

func RunFile() error {
	f, err := newFile()
	if err != nil {
		return err
	}

	f.flagParse()

	err = f.conf.config()
	if err != nil {
		return err
	}

	err = f.parseEntries()
	if err != nil {
		return err
	}

	go f.entriesChange()

	return runPublisher(&f.conf)
}

func (f *File) entriesChange() {
	f.parseRecords()
	f.recordsChange()

	err := runFsNotify(f.handleNotify,
		f.mdnsConfFile,
		f.hostsConfFile,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *File) recordsChange() {
	f.conf.recordsChangeChan <- f.records
}

func (f *File) handleNotify(event fsnotify.Event) {
	hasChange := false
	if event.Has(fsnotify.Write) {
		if strings.Contains(event.Name, f.mdnsConfFile) {
			err := f.parseMdnsEntries()
			if err != nil {
				log.Println(err)
			} else {
				hasChange = true
			}
		}

		if strings.Contains(event.Name, f.hostsConfFile) {
			err := f.parseHostsEntries()
			if err != nil {
				log.Println(err)
			} else {
				hasChange = true
			}
		}
	}

	if hasChange {
		f.parseRecords()
		f.recordsChange()
	}
}

func (f *File) parseRecords() {
	var err error
	var entries []string
	var records = zone.NewRecords()
	entries = append(entries, f.mdnsEntries...)
	entries = append(entries, f.hostsEntries...)

	for _, entry := range entries {
		if entry == "" {
			continue
		}

		err = records.New(entry)
		if err != nil {
			log.Fatalf("entry[%s]: %v\n", entry, err)
		}
	}

	f.records = records
}

func (f *File) parseEntries() error {
	err := f.parseMdnsEntries()
	if err != nil {
		return err
	}

	err = f.parseHostsEntries()
	if err != nil {
		return err
	}

	return nil
}

func (f *File) parseMdnsEntries() error {
	entries, err := f.readMdnsEntries()
	if err != nil {
		return err
	}

	f.mdnsEntries = entries

	return nil
}

func (f *File) parseHostsEntries() error {
	entries, err := f.readHostsEntries()
	if err != nil {
		return err
	}

	f.hostsEntries = entries

	return nil
}

func (f *File) readMdnsEntries() ([]string, error) {
	if f.mdnsConfFile == "" {
		//return nil, fmt.Errorf("no config file found")
		f.mdnsConfFile = defaultMdns
	}

	return f.readLinesFromFile(f.mdnsConfFile)
}

func (f *File) readHostsEntries() ([]string, error) {
	if f.hostsConfFile == "" {
		//return nil, fmt.Errorf("no local entries config file found")
		f.hostsConfFile = defaultHosts
	}

	iniFile, err := ini.Load(f.hostsConfFile)
	if err != nil {
		return nil, err
	}

	ifsIP, err := newInterfacesIp()
	if err != nil {
		return nil, err
	}

	var dnsEntry = "%s.local. %d IN %s %s"
	var entries []string

	sections := iniFile.Sections()
	for _, section := range sections {
		iIP, ok := ifsIP[section.Name()]
		if !ok {
			continue
		}

		var hosts = make(map[string]int)
		for _, key := range section.Keys() {
			ttl, err := key.Duration()
			if err != nil {
				log.Println(err)
			}
			if ttl < time.Minute*2 {
				ttl = time.Minute * 2
			}

			hosts[key.Name()] = int(ttl.Seconds())
		}

		for host, ttl := range hosts {
			entries = append(entries, fmt.Sprintf(dnsEntry, host, ttl, "A", iIP.ipv4Addr()))
			entries = append(entries, fmt.Sprintf(dnsEntry, host, ttl, "AAAA", iIP.ipv6Addr()))
		}
	}

	return entries, nil
}

func (f *File) readLinesFromFile(s string) ([]string, error) {
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

	flag.StringVar(&f.mdnsConfFile, "conf", "", "config file path")
	flag.StringVar(&f.mdnsConfFile, "c", "", "config file path")

	flag.Parse()
}
