package source

import (
	"errors"
	"fmt"
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/database"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// Section 分段查询
type Section struct {
	storage.Storage
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
	state         *database.State
}

func NewSectionFromQuery(config FromQueryConfig) (*Section, error) {
	s, err := newSource(config.Config)
	if err != nil {
		return nil, err
	}

	s.query = config.Query

	return s, nil
}

func NewSectionFromTable(config FromTableConfig) (*Section, error) {
	s, err := newSource(config.Config)
	if err != nil {
		return nil, err
	}

	s.fields = config.Fields

	fieldsSting := fmt.Sprintf("`%s`", strings.Join(s.fields, "`,`"))
	s.query = fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` > ? ORDER BY `%s` ASC limit ?", fieldsSting, s.tableName, s.idName, s.idName)

	return s, nil
}

func newSource(config Config) (*Section, error) {
	s := &Section{}
	err := s.config(config)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Section) GetDriver() simple.Driver {
	return s.Driver
}

func (s *Section) State() []string {
	return []string{s.state.State()}
}

func (s *Section) Summary() []string {
	return []string{fmt.Sprintf("%s: total: %d, total page: %d, page size: %d, start id: %d, end id: %d ",
		s.Title(),
		s.total,
		s.totalPage,
		s.pageSize,
		s.startId,
		s.maxId)}
}

func (s *Section) Scan() {
	s.state.Status = "scan"
	s.state.Duration.Start()

	go s.prepare()

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

		s.state.Duration.Stop()
		s.state.Status = "finish"
	}()
}

func (s *Section) ReceiveChan() <-chan []map[string]interface{} {
	return s.rowsChan
}

func (s *Section) scanRows() {
	for {
		idRange, ok := <-s.idRangeChan
		if !ok {
			break
		}

		atomic.AddInt64(&s.state.Page, 1)

		leftId := idRange[0]
		size := idRange[1]

		rows, err := s.QueryInterface(s.query, leftId, size)
		if err != nil {
			s.LogFatalWithTitle("left id:", leftId, err)
		}

		if len(rows) == 0 {
			continue
		}

		atomic.AddInt64(&s.state.Amount, int64(len(rows)))

		s.rowsChan <- rows
	}
}

func (s *Section) prepare() {
	s.idRangeChan = make(chan []int64, s.concurrency)

	s.idRange()
	s.queryPages()

	go func() {
		s.idRangeSw.Wait()
		close(s.idRangeChan)
	}()

	s.state.PageSize = s.pageSize
	s.state.TotalPage = s.totalPage
	s.state.Total = s.total
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
			s.LogFatalWithTitle(fmt.Sprintf("source[] last startid %d endId %d id: %d left id: %d right id: %d", startId, endId, leftId, rightId, err))
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

func (s *Section) config(config Config) (err error) {
	s.concurrency = config.Concurrency
	s.total = config.Maximum
	s.startId = config.StartId
	s.maxId = config.EndId
	s.pageSize = config.PageSize
	s.tableName = config.TableName
	s.idName = config.IdName

	s.Driver, err = simple.NewDriver(config.Driver, config.Dsn)
	if err != nil {
		return
	}

	if s.concurrency <= 0 {
		err = errors.New(fmt.Sprintf("concurrency[%d] is error", s.concurrency))
		return
	}

	if s.pageSize <= 0 {
		err = errors.New(fmt.Sprintf("page size[%d] is error", s.pageSize))
		return
	}

	if s.total > 0 && s.pageSize > s.total {
		s.pageSize = s.total
	}

	s.state = database.NewState()
	s.state.Concurrency = s.concurrency
	s.state.Title = s.Title()
	s.state.Status = "config"

	return
}

func (s *Section) idRange() {
	query := fmt.Sprintf("SELECT MAX(`%s`) AS `maxId`, MIN(`%s`) AS `minId` FROM `%s`", s.idName, s.idName, s.tableName)
	res, err := s.QueryInterfaceRow(query)
	if err != nil {
		s.LogFatalWithTitle("mysql id:", err)
	}

	maxId, err := strconv.ParseInt(string(res["maxId"].([]uint8)), 10, 64)
	if err != nil {
		s.LogFatalWithTitle("mysql max:", err)
	}

	minId, err := strconv.ParseInt(string(res["minId"].([]uint8)), 10, 64)
	if err != nil {
		s.LogFatalWithTitle("mysql min:", err)
	}

	if minId > s.startId {
		s.startId = minId - 1
	}

	if s.maxId <= 0 || s.maxId > maxId {
		s.maxId = maxId
	}

	if s.maxId < s.startId {
		s.LogFatalWithTitle(fmt.Sprintf("mysql max id %d < start id %d", s.maxId, s.startId))
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

func (s *Section) Title() string {
	return fmt.Sprintf("Source driver[%s]", s.DriverName())
}
