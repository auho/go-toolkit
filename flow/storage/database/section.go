package database

import (
	"fmt"
	"github.com/auho/go-simple-db/simple"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// Section 分段查询
type Section struct {
	simple.Driver
	scanSw        sync.WaitGroup
	idRangeSw     sync.WaitGroup
	concurrency   int
	pageSize      int64
	totalPage     int64
	startId       int64
	maxId         int64
	total         int64
	tableName     string
	idName        string
	query         string
	fields        []string
	failureLastId []int
	idRangeChan   chan []int64
	rowsChan      chan []map[string]interface{}
	state         *State
}

func NewSectionFromQuery(config FromQueryConfig) *Section {
	m := newSource(config.Config)
	m.query = config.Query

	m.prepare()

	return m
}

func NewSectionFromTable(config FromTableConfig) *Section {
	m := newSource(config.Config)
	m.fields = config.Fields

	fieldsSting := fmt.Sprintf("`%s`", strings.Join(m.fields, "`,`"))
	m.query = fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` > ? ORDER BY `%s` ASC limit ?", fieldsSting, m.tableName, m.idName, m.idName)

	m.prepare()

	return m
}

func newSource(config Config) *Section {
	var err error
	m := &Section{}
	m.state = newState()
	m.config(config)

	m.Driver, err = simple.NewDriver(config.Driver, config.Dsn)
	if err != nil {
		m.logFatalWithTitle(err)
	}

	m.idRangeChan = make(chan []int64, m.concurrency)

	return m
}

func (s *Section) GetDriver() simple.Driver {
	return s.Driver
}

func (s *Section) State() string {
	return s.state.State()
}

func (s *Section) Summary() string {
	return fmt.Sprintf("%s: total: %d, total page: %d, page size: %d, start id: %d, end id: %d ",
		s.title(),
		s.total,
		s.totalPage,
		s.pageSize,
		s.startId,
		s.maxId)
}

func (s *Section) Scan() {
	s.rowsChan = make(chan []map[string]interface{}, s.concurrency)
	for i := 0; i < s.concurrency; i++ {
		s.scanSw.Add(1)
		go func() {
			s.scanRows()
			s.scanSw.Done()
		}()
	}

	go func() {
		s.scanSw.Wait()
		close(s.rowsChan)
	}()
}

func (s *Section) ReceiveChan() <-chan []map[string]interface{} {
	return s.rowsChan
}

func (s *Section) Next() ([]map[string]interface{}, bool) {
	ims, ok := <-s.rowsChan

	return ims, ok
}

func (s *Section) scanRows() {
	for {
		idRange, ok := <-s.idRangeChan
		if !ok {
			break
		}

		atomic.AddInt64(&s.state.page, 1)

		leftId := idRange[0]
		size := idRange[1]

		rows, err := s.QueryInterface(s.query, leftId, size)
		if err != nil {
			s.logFatalWithTitle("left id:", leftId, err)
		}

		if len(rows) == 0 {
			continue
		}

		atomic.AddInt64(&s.state.amount, int64(len(rows)))

		s.rowsChan <- rows
	}
}

func (s *Section) prepare() {
	s.idRange()
	s.queryPages()

	go func() {
		s.idRangeSw.Wait()
		close(s.idRangeChan)
	}()

	s.state.title = s.title()
	s.state.pageSize = s.pageSize
	s.state.totalPage = s.totalPage
	s.state.total = s.total
}

// queryPages 分段查询
func (s *Section) queryPages() {
	shard := int64(math.Ceil(float64(s.totalPage) / float64(s.concurrency)))
	shardSize := shard * s.pageSize

	for i := 0; i < s.concurrency; i++ {
		i64 := int64(i)
		startId := s.startId + i64*shardSize
		endId := startId + shardSize

		if endId > s.maxId {
			endId = s.maxId
		}

		s.idRangeSw.Add(1)
		go func() {
			s.queryPage(startId, endId)
			s.idRangeSw.Done()
		}()

		if endId >= s.maxId {
			break
		}
	}
}

// queryPage 查询分段
func (s *Section) queryPage(startId, endId int64) {
	var rightId int64 = 0
	leftId := startId

	for {
		rightId = leftId + s.pageSize
		if rightId > endId {
			rightId = endId
		}

		query := fmt.Sprintf("SELECT MAX(`%s`) AS `id` FROM `%s` WHERE `%s` > ? AND `%s` <= ? ORDER BY `%s` DESC LIMIT 1", s.idName, s.tableName, s.idName, s.idName, s.idName)
		res, err := s.QueryFieldInterface("id", query, leftId, rightId)
		if err != nil {
			s.logFatalWithTitle(fmt.Sprintf("source[] last startid %d endId %d id: %d left id: %d right id: %d", startId, endId, leftId, rightId, err))
		}

		if res != nil {
			s.idRangeChan <- []int64{leftId, rightId - leftId}
		}

		if rightId >= endId {
			break
		}

		leftId += s.pageSize
	}
}

func (s *Section) config(config Config) {
	s.concurrency = config.Concurrency
	s.total = config.Maximum
	s.startId = config.StartId
	s.maxId = config.EndId
	s.pageSize = config.PageSize
	s.tableName = config.TableName
	s.idName = config.IdName

	if s.concurrency <= 0 {
		s.logFatal(fmt.Sprintf("driver[%s] concurrency[%d] is error", config.Driver, s.concurrency))
	}

	if s.pageSize <= 0 {
		s.logFatal(fmt.Sprintf("driver[%s] concurrency[%d] is error", config.Driver, s.pageSize))
	}

	if s.total > 0 && s.pageSize > s.total {
		s.pageSize = s.total
	}

	if s.concurrency < 1 {
		s.concurrency = 1
	}
}

func (s *Section) idRange() {
	query := fmt.Sprintf("SELECT MAX(`%s`) AS `maxId`, MIN(`%s`) AS `minId` FROM `%s`", s.idName, s.idName, s.tableName)
	res, err := s.QueryInterfaceRow(query)
	if err != nil {
		s.logFatalWithTitle("mysql id:", err)
	}

	maxId, err := strconv.ParseInt(string(res["maxId"].([]uint8)), 10, 64)
	if err != nil {
		s.logFatalWithTitle("mysql max:", err)
	}

	minId, err := strconv.ParseInt(string(res["minId"].([]uint8)), 10, 64)
	if err != nil {
		s.logFatalWithTitle("mysql min:", err)
	}

	if minId > s.startId {
		s.startId = minId - 1
	}

	if s.maxId <= 0 || s.maxId > maxId {
		s.maxId = maxId
	}

	if s.maxId < s.startId {
		s.logFatalWithTitle(fmt.Sprintf("mysql max id %d < start id %d", s.maxId, s.startId))
	}

	total := s.maxId - s.startId + 1
	if s.total == 0 {
		s.total = total
	} else {
		if s.total < total {
			s.maxId = s.startId + s.total
		} else if s.total > total {
			s.total = total
		}
	}

	if s.pageSize > s.total {
		s.pageSize = s.total
	}

	s.totalPage = int64(math.Ceil(float64(s.total) / float64(s.pageSize)))
}

func (s *Section) title() string {
	return fmt.Sprintf("Source driver[%s]", s.DriverName())
}

func (s *Section) logFatalWithTitle(v ...any) {
	log.Fatal(append([]interface{}{s.title()}, v...)...)
}

func (s *Section) logFatal(v ...any) {
	log.Fatal(v...)
}
