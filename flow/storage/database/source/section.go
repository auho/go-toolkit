package source

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/database"
)

var _ storage.Sourceor[storage.MapEntry] = (*Section[storage.MapEntry])(nil)
var _ database.Driver = (*Section[storage.MapEntry])(nil)

type sectionQuery[E storage.Entry] interface {
	query(se *Section[E], startId, size int64) ([]E, error)
	duplicate([]E) []E
}

// Section 分段查询
type Section[E storage.Entry] struct {
	storage.Storage
	db        *database.DB
	scanSw    sync.WaitGroup
	idRangeSw sync.WaitGroup
	conf      *QueryConfig

	total     int64
	totalPage int64
	startId   int64 // 开区间
	maxId     int64 // 闭区间

	failureLastId []int
	idRangeChan   chan []int64
	rowsChan      chan []E
	state         *storage.PageState
	sr            sectionQuery[E]
}

func newSection[E storage.Entry](config *QueryConfig, sr sectionQuery[E], b database.BuildDb) (*Section[E], error) {
	s := &Section[E]{}
	err := s.config(config, b)
	if err != nil {
		return nil, err
	}

	s.sr = sr

	return s, nil
}

func (s *Section[E]) DB() *database.DB {
	return s.db
}

func (s *Section[E]) State() []string {
	return []string{s.state.Overview()}
}

func (s *Section[E]) Summary() []string {
	return []string{fmt.Sprintf("%s: total: %d, total page: %d, page size: %d, start id: %d, end id: %d ",
		s.Title(),
		s.total,
		s.totalPage,
		s.conf.PageSize,
		s.startId,
		s.maxId)}
}

func (s *Section[E]) Duplicate(items []E) []E {
	return s.sr.duplicate(items)
}

func (s *Section[E]) Scan() error {
	s.state.StatusScan()
	s.state.DurationStart()
	s.idRangeChan = make(chan []int64, s.conf.Concurrency)
	s.rowsChan = make(chan []E, s.conf.Concurrency)

	err := s.idRange()
	if err != nil {
		return err
	}

	go s.idSection()

	for i := 0; i < s.conf.Concurrency; i++ {
		s.scanSw.Add(1)
		go func() {
			s.scanRows()
			s.scanSw.Done()
		}()
	}

	go func() {
		s.scanSw.Wait()
		close(s.rowsChan)

		s.state.DurationStop()
		s.state.StatusFinish()
	}()

	return nil
}

func (s *Section[E]) ReceiveChan() <-chan []E {
	return s.rowsChan
}

func (s *Section[E]) scanRows() {
	for idRange := range s.idRangeChan {
		atomic.AddInt64(&s.state.Page, 1)

		leftId := idRange[0]
		size := idRange[1]

		rows, err := s.sr.query(s, leftId, size)
		if err != nil {
			s.LogFatalWithTitle("left id:", leftId, err)
		}

		if len(rows) == 0 {
			continue
		}

		s.state.AddAmount(int64(len(rows)))

		s.rowsChan <- rows
	}
}

func (s *Section[E]) idSection() {
	s.queryPages()

	go func() {
		s.idRangeSw.Wait()
		close(s.idRangeChan)
	}()

	s.state.PageSize = s.conf.PageSize
	s.state.TotalPage = s.totalPage
	s.state.Total = s.total
}

// queryPages 分段
func (s *Section[E]) queryPages() {
	shard := int64(math.Ceil(float64(s.totalPage) / float64(s.conf.Concurrency)))
	shardSize := shard * s.conf.PageSize

	for i := 0; i < s.conf.Concurrency; i++ {
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
func (s *Section[E]) queryPage(startId, endId int64) {
	var rightId int64 = 0
	leftId := startId

	var row struct {
		Id int64
	}

	for {
		rightId = leftId + s.conf.PageSize
		if rightId > endId {
			rightId = endId
		}

		err := s.db.Table(s.conf.TableName).
			Select(fmt.Sprintf("MAX(%s) AS id", s.conf.IdName)).
			Where(fmt.Sprintf("%s > ? AND %s <= ?", s.conf.IdName, s.conf.IdName), leftId, rightId).
			Group(fmt.Sprintf("%s desc", s.conf.IdName)).
			Scan(&row).Error
		if err != nil {
			s.LogFatalWithTitle(fmt.Sprintf("query page: start id %d end id %d; left id: %d right id: %d; %v", startId, endId, leftId, rightId, err))
		}

		s.idRangeChan <- []int64{leftId, rightId - leftId}

		if rightId >= endId {
			break
		}

		leftId += s.conf.PageSize
	}
}

func (s *Section[E]) config(config *QueryConfig, b database.BuildDb) (err error) {
	s.conf = config

	s.total = config.Maximum
	s.startId = config.StartId
	s.maxId = config.EndId

	s.db, err = b()
	if err != nil {
		return
	}

	err = s.db.Ping()
	if err != nil {
		return
	}

	if s.conf.Concurrency <= 0 {
		err = fmt.Errorf("concurrency[%d] is error", s.conf.Concurrency)
		return
	}

	if s.conf.PageSize <= 0 {
		err = fmt.Errorf("page size[%d] is error", s.conf.PageSize)
		return
	}

	if s.total > 0 && s.conf.PageSize > s.total {
		s.conf.PageSize = s.total
	}

	s.state = storage.NewPageState()
	s.state.Concurrency = s.conf.Concurrency
	s.state.Title = s.Title()
	s.state.StatusConfig()

	return
}

func (s *Section[E]) idRange() error {
	var row struct {
		Max int64
		Min int64
	}

	query := fmt.Sprintf("MAX(%s) AS max, MIN(%s) AS min", s.conf.IdName, s.conf.IdName)
	err := s.db.Table(s.conf.TableName).Select(query).Scan(&row).Error
	if err != nil {
		return fmt.Errorf("id range %w", err)
	}

	if row.Min > s.startId {
		s.startId = row.Min - 1
	}

	if s.maxId <= 0 || s.maxId > row.Max {
		s.maxId = row.Max
	}

	if s.maxId < s.startId {
		return fmt.Errorf("mysql max id %d < start id %d", s.maxId, s.startId)
	}

	total := s.maxId - s.startId
	if s.total == 0 {
		s.total = total
	} else {
		if s.total < total {
			s.maxId = s.startId + s.total
		} else if s.total > total {
			s.total = total
		}
	}

	if s.conf.PageSize > s.total {
		s.conf.PageSize = s.total
	}

	s.totalPage = int64(math.Ceil(float64(s.total) / float64(s.conf.PageSize)))

	return nil
}

func (s *Section[E]) Title() string {
	return fmt.Sprintf("Sourceor db[%s]", s.db.Name())
}

func (s *Section[E]) Close() error {
	return s.db.Close()
}
