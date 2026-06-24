# Flow 包优化实现计划（阶段一：紧急修复 + 阶段二：错误处理改造）

> **面向 AI 代理的工作者：** 必需子技能：使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务实现此计划。步骤使用复选框（`- [ ]`）语法来跟踪进度。

**目标：** 修复 flow 包中的 BUG、拼写错误、错误处理反模式，使库代码不再使用 log.Fatal/panic 终止进程

**架构：** 分两阶段推进——阶段一修复 4 个 BUG + 拼写错误 + channel 缓冲问题（无 API 变更）；阶段二将 log.Fatal/panic 替换为 error 返回，并增加 goroutine 错误恢复机制

**技术栈：** Go 1.21, 标准库 testing

---

## 文件结构

| 文件 | 职责 | 变更类型 |
|------|------|---------|
| `storage/redis/sets.go` | Redis Sets KeyType 定义 | 修改 |
| `storage/redis/sorted_sets.go` | Redis SortedSets KeyType 定义 | 修改 |
| `storage/mock/source/source.go` | Mock 数据源 | 修改 |
| `flow/flow.go` | 核心编排 | 修改 |
| `storage/state.go` | 状态管理 | 修改 |
| `storage/redis/destination/key.go` | Redis 目标端 | 修改 |
| `storage/storage.go` | 存储基础结构 | 修改 |
| `storage/database/source/section.go` | 数据库分段查询源 | 修改 |
| `storage/database/destination/destination.go` | 数据库目标端 | 修改 |
| `storage/redis/destination/hashes.go` | Redis Hashes 目标端 | 修改 |
| `storage/redis/destination/lists.go` | Redis Lists 目标端 | 修改 |
| `storage/redis/destination/sets.go` | Redis Sets 目标端 | 修改 |
| `storage/redis/destination/sorted_sets.go` | Redis SortedSets 目标端 | 修改 |
| `storage/file/source/line.go` | 文件行源 | 修改 |
| `storage/file/destination/line.go` | 文件行目标端 | 修改 |
| `flow/flow_bug_test.go` | BUG 修复测试 | 新建 |
| `storage/redis/redis_type_test.go` | Redis 类型测试 | 新建 |
| `storage/mock/source/source_test.go` | Mock 源测试 | 新建 |
| `storage/state_test.go` | State 拼写测试 | 新建 |
| `storage/storage_test.go` | Storage 错误处理测试 | 新建 |

---

## 任务组 A：BUG 修复（P0 紧急）

### 任务 1：修复 Redis Sets Type() 返回值

**文件：**
- 修改：`flow/storage/redis/sets.go:18`
- 新建：`flow/storage/redis/redis_type_test.go`

- [ ] **步骤 1：编写失败的测试**

```go
// flow/storage/redis/redis_type_test.go
package redis

import "testing"

func TestSetsType(t *testing.T) {
	s := Sets{}
	if s.Type() != KeyTypeSet {
		t.Errorf("Sets.Type() = %q, want %q", s.Type(), KeyTypeSet)
	}
}

func TestSortedSetsType(t *testing.T) {
	s := SortedSets{}
	if s.Type() != KeyTypeSortedSets {
		t.Errorf("SortedSets.Type() = %q, want %q", s.Type(), KeyTypeSortedSets)
	}
}

func TestListsType(t *testing.T) {
	l := Lists{}
	if l.Type() != KeyTypeList {
		t.Errorf("Lists.Type() = %q, want %q", l.Type(), KeyTypeList)
	}
}

func TestHashesType(t *testing.T) {
	h := Hashes{}
	if h.Type() != KeyTypeHash {
		t.Errorf("Hashes.Type() = %q, want %q", h.Type(), KeyTypeHash)
	}
}
```

- [ ] **步骤 2：运行测试验证失败**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./storage/redis/ -run TestSetsType -v`
预期：FAIL — `Sets.Type() = "lists", want "sets"`

- [ ] **步骤 3：修复 Sets.Type()**

将 `flow/storage/redis/sets.go` 中：

```go
func (l *Sets) Type() KeyType {
	return KeyTypeList
}
```

改为：

```go
func (l *Sets) Type() KeyType {
	return KeyTypeSet
}
```

- [ ] **步骤 4：运行测试验证通过**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./storage/redis/ -run TestSetsType -v`
预期：PASS

- [ ] **步骤 5：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add storage/redis/sets.go storage/redis/redis_type_test.go
git commit -m "fix: Sets.Type() returns KeyTypeSet instead of KeyTypeList"
```

---

### 任务 2：修复 Redis SortedSets Type() 返回值

**文件：**
- 修改：`flow/storage/redis/sorted_sets.go:18`

- [ ] **步骤 1：运行已有测试验证失败**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./storage/redis/ -run TestSortedSetsType -v`
预期：FAIL — `SortedSets.Type() = "lists", want "sortedSets"`

- [ ] **步骤 2：修复 SortedSets.Type()**

将 `flow/storage/redis/sorted_sets.go` 中：

```go
func (l *SortedSets) Type() KeyType {
	return KeyTypeList
}
```

改为：

```go
func (l *SortedSets) Type() KeyType {
	return KeyTypeSortedSets
}
```

- [ ] **步骤 3：运行测试验证通过**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./storage/redis/ -run TestSortedSetsType -v`
预期：PASS

- [ ] **步骤 4：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add storage/redis/sorted_sets.go
git commit -m "fix: SortedSets.Type() returns KeyTypeSortedSets instead of KeyTypeList"
```

---

### 任务 3：修复 Mock Source pageSize 配置校验

**文件：**
- 修改：`flow/storage/mock/source/source.go:56`
- 新建：`flow/storage/mock/source/source_test.go`

- [ ] **步骤 1：编写失败的测试**

```go
// flow/storage/mock/source/source_test.go
package source

import "testing"

func TestMockDefaultPageSize(t *testing.T) {
	m := newMock[map[string]any](Config{
		Total:    100,
		PageSize: 0,
	}, &SliceMap{})

	if m.pageSize <= 0 {
		t.Errorf("pageSize should be set to default when <= 0, got %d", m.pageSize)
	}
}

func TestMockDefaultTotal(t *testing.T) {
	m := newMock[map[string]any](Config{
		Total:    0,
		PageSize: 10,
	}, &SliceMap{})

	if m.total <= 0 {
		t.Errorf("total should be set to default when <= 0, got %d", m.total)
	}
}
```

- [ ] **步骤 2：运行测试验证失败**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./storage/mock/source/ -run TestMockDefaultPageSize -v`
预期：FAIL — `pageSize should be set to default when <= 0, got 0`

- [ ] **步骤 3：修复 pageSize 校验逻辑**

将 `flow/storage/mock/source/source.go` 中：

```go
if m.pageSize <= 0 {
    m.total = 1e1
}
```

改为：

```go
if m.pageSize <= 0 {
    m.pageSize = 1e1
}
```

- [ ] **步骤 4：运行测试验证通过**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./storage/mock/source/ -run TestMockDefault -v`
预期：PASS

- [ ] **步骤 5：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add storage/mock/source/source.go storage/mock/source/source_test.go
git commit -m "fix: mock source sets pageSize (not total) when pageSize <= 0"
```

---

### 任务 4：Flow.check() 增加 Source nil 校验

**文件：**
- 修改：`flow/flow/flow.go:63-67`
- 新建：`flow/flow/flow_bug_test.go`

- [ ] **步骤 1：编写失败的测试**

```go
// flow/flow/flow_bug_test.go
package flow

import "testing"

func TestRunFlowNoSource(t *testing.T) {
	err := RunFlow[map[string]any]()
	if err == nil {
		t.Error("expected error when source is nil")
	}
}

func TestRunFlowNoAction(t *testing.T) {
	err := RunFlow[map[string]any]()
	if err == nil {
		t.Error("expected error when actions are empty")
	}
}
```

- [ ] **步骤 2：运行测试验证失败**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./flow/ -run TestRunFlowNoSource -v`
预期：FAIL 或 panic（因为 source 为 nil 时未校验直接使用）

- [ ] **步骤 3：修复 Flow.check()**

将 `flow/flow/flow.go` 中：

```go
func (f *Flow[E]) check() error {
	if len(f.actions) <= 0 {
		return errors.New("action not found")
	}

	return nil
}
```

改为：

```go
func (f *Flow[E]) check() error {
	if f.source == nil {
		return errors.New("source not found")
	}

	if len(f.actions) <= 0 {
		return errors.New("action not found")
	}

	return nil
}
```

- [ ] **步骤 4：运行测试验证通过**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./flow/ -run TestRunFlowNo -v`
预期：PASS

- [ ] **步骤 5：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add flow/flow.go flow/flow_bug_test.go
git commit -m "fix: Flow.check() validates source is not nil"
```

---

## 任务组 B：拼写错误修正（P2）

### 任务 5：修正所有拼写错误

**文件：**
- 修改：`flow/flow/flow.go:101,109`
- 修改：`flow/storage/state.go:104`
- 修改：`flow/storage/redis/destination/key.go:133`

- [ ] **步骤 1：编写拼写检查测试**

```go
// flow/storage/state_test.go
package storage

import "testing"

func TestTotalStateOverview(t *testing.T) {
	ts := NewTotalState()
	ts.status = StatusConfig
	ts.Concurrency = 4
	ts.amount = 100
	ts.Total = 200
	overview := ts.Overview()
	if overview == "" {
		t.Error("overview should not be empty")
	}
}
```

- [ ] **步骤 2：运行测试确认当前行为**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./storage/ -run TestTotalStateOverview -v`
预期：PASS（但输出包含拼写错误 "Concurrentcy"）

- [ ] **步骤 3：修复 flow.go 拼写错误**

将 `flow/flow/flow.go` 第 101 行：

```go
return fmt.Errorf("actioins prepare error;%w", err)
```

改为：

```go
return fmt.Errorf("actions prepare error; %w", err)
```

将第 109 行：

```go
return fmt.Errorf("actioins run error;%w", err)
```

改为：

```go
return fmt.Errorf("actions run error; %w", err)
```

- [ ] **步骤 4：修复 state.go 拼写错误**

将 `flow/storage/state.go` 第 104 行：

```go
return fmt.Sprintf("Status: %s, Concurrentcy:%d, Amount: %d/%d, Duration: %s",
```

改为：

```go
return fmt.Sprintf("Status: %s, Concurrency: %d, Amount: %d/%d, Duration: %s",
```

- [ ] **步骤 5：修复 redis/destination/key.go 拼写错误**

将 `flow/storage/redis/destination/key.go` 第 133 行：

```go
return fmt.Sprintf("Destiantion redis[%s] %s", k.keyer.Type(), k.keyName)
```

改为：

```go
return fmt.Sprintf("Destination redis[%s] %s", k.keyer.Type(), k.keyName)
```

- [ ] **步骤 6：运行全量编译验证**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go build ./...`
预期：编译成功，无错误

- [ ] **步骤 7：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add flow/flow.go storage/state.go storage/redis/destination/key.go storage/state_test.go
git commit -m "fix: correct typos (actioins→actions, Concurrentcy→Concurrency, Destiantion→Destination)"
```

---

## 任务组 C：错误处理改造（P1 重要）

### 任务 6：移除 Storage.LogFatal，改为返回 error

**文件：**
- 修改：`flow/storage/storage.go`
- 修改：`flow/storage/database/source/section.go`
- 修改：`flow/storage/redis/source/key.go`
- 修改：`flow/storage/redis/source/scan.go`
- 新建：`flow/storage/storage_test.go`

- [ ] **步骤 1：编写失败的测试**

```go
// flow/storage/storage_test.go
package storage

import "testing"

func TestStorageFatalError(t *testing.T) {
	s := &Storage{}
	err := s.FatalError("test title", nil)
	if err != nil {
		t.Errorf("FatalError with nil err should return nil, got %v", err)
	}

	fakeErr := fmt.Errorf("some error")
	err = s.FatalError("test title", fakeErr)
	if err == nil {
		t.Error("FatalError should return error")
	}
	if !strings.Contains(err.Error(), "test title") {
		t.Error("error should contain title")
	}
	if !strings.Contains(err.Error(), "some error") {
		t.Error("error should contain original error message")
	}
}
```

注意：需要在 `storage_test.go` 顶部添加 import：

```go
import (
	"fmt"
	"strings"
	"testing"
)
```

- [ ] **步骤 2：在 Storage 上添加 FatalError 方法**

在 `flow/storage/storage.go` 中，保留原有 `LogFatal` 和 `LogFatalWithTitle`（标记废弃），新增 `FatalError` 方法：

```go
func (s *Storage) FatalError(title string, err error) error {
	if err == nil {
		return nil
	}
	if title != "" {
		return fmt.Errorf("%s; %w", title, err)
	}
	return err
}
```

- [ ] **步骤 3：运行测试验证通过**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go test ./storage/ -run TestStorageFatalError -v`
预期：PASS

- [ ] **步骤 4：改造 database/source/section.go — scanRows()**

将 `flow/storage/database/source/section.go` 中 `scanRows()` 方法，从：

```go
func (s *Section[E]) scanRows() {
	for i := 0; i < s.conf.Concurrency; i++ {
		s.scanSw.Add(1)
		go func() {
			for idRange := range s.idRangeChan {
				atomic.AddInt64(&s.state.Page, 1)

				leftId := idRange[0]
				rightId := idRange[1]

				rows, err := s.sq.Query(s, leftId, rightId)
				if err != nil {
					s.LogFatalWithTitle(fmt.Sprintf("left id[%d] - right id[%d]", leftId, rightId), err)
				}

				if len(rows) == 0 {
					continue
				}

				s.state.AddAmount(int64(len(rows)))

				s.rowsChan <- rows
			}

			s.scanSw.Done()
		}()
	}
}
```

改为：

```go
func (s *Section[E]) scanRows() {
	for i := 0; i < s.conf.Concurrency; i++ {
		s.scanSw.Add(1)
		go func() {
			defer s.scanSw.Done()

			for idRange := range s.idRangeChan {
				atomic.AddInt64(&s.state.Page, 1)

				leftId := idRange[0]
				rightId := idRange[1]

				rows, err := s.sq.Query(s, leftId, rightId)
				if err != nil {
					fmt.Fprintf(os.Stderr, "left id[%d] - right id[%d]; %v\n", leftId, rightId, err)
					continue
				}

				if len(rows) == 0 {
					continue
				}

				s.state.AddAmount(int64(len(rows)))

				s.rowsChan <- rows
			}
		}()
	}
}
```

注意：需要在 `section.go` 顶部添加 `"fmt"` 和 `"os"` 到 import。

- [ ] **步骤 5：改造 redis/source/key.go — config()**

将 `flow/storage/redis/source/key.go` 中 `config()` 方法，从：

```go
if k.keyName == "" {
    k.LogFatalWithTitle("key name is empty")
}

if config.Options == nil {
    k.LogFatalWithTitle("config options is nil")
}
```

改为：

```go
if k.keyName == "" {
    return fmt.Errorf("key name is empty")
}

if config.Options == nil {
    return fmt.Errorf("config options is nil")
}
```

- [ ] **步骤 6：改造 redis/source/scan.go — config()**

将 `flow/storage/redis/source/scan.go` 中 `config()` 方法，从：

```go
if config.Options == nil {
    s.LogFatalWithTitle("config options is nil")
}
```

改为：

```go
if config.Options == nil {
    return fmt.Errorf("config options is nil")
}
```

- [ ] **步骤 7：运行全量编译验证**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go build ./...`
预期：编译成功

- [ ] **步骤 8：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add storage/storage.go storage/storage_test.go storage/database/source/section.go storage/redis/source/key.go storage/redis/source/scan.go
git commit -m "refactor: replace Storage.LogFatal with error returns in source implementations"
```

---

### 任务 7：替换 panic 为 error 返回 — database/destination

**文件：**
- 修改：`flow/storage/database/destination/destination.go`

- [ ] **步骤 1：为 Destination 添加 error channel**

将 `flow/storage/database/destination/destination.go` 中 `Destination` 结构体，从：

```go
type Destination[E storage.Entry] struct {
	storage.Storage
	db     *database.DB
	isDone bool

	isTruncate  bool
	concurrency int
	table       string
	pageSize    int64

	state     *storage.State
	doWg      sync.WaitGroup
	dst       Destinationer[E]
	itemsChan chan []E
}
```

改为：

```go
type Destination[E storage.Entry] struct {
	storage.Storage
	db     *database.DB
	isDone bool

	isTruncate  bool
	concurrency int
	table       string
	pageSize    int64

	state     *storage.State
	doWg      sync.WaitGroup
	dst       Destinationer[E]
	itemsChan chan []E
	errChan   chan error
	firstErr  error
}
```

- [ ] **步骤 2：改造 do() 方法，用 errChan 替代 panic**

将 `do()` 方法从：

```go
func (d *Destination[E]) do() {
	duration := timing.NewDuration()
	duration.Start()
	var descItems []E

	duration.Begin()
	for items := range d.itemsChan {
		if len(items) <= 0 {
			continue
		}

		descItems = append(descItems, items...)

		_len := len(descItems)
		_start := 0
		_end := 0
		_size := int(d.pageSize)
		for {
			_end = _start + _size
			if _end <= _len {
				err := d.dst.Exec(d, descItems[_start:_end])
				if err != nil {
					panic(err)
				}

				d.state.AddAmount(int64(_size))

				_start += _size
			} else {
				descItems = slices.Clone(descItems[_start:])
				descItems = slices.Clip(descItems)

				break
			}
		}
	}

	if len(descItems) > 0 {
		err := d.dst.Exec(d, descItems)
		if err != nil {
			panic(err)
		}

		d.state.AddAmount(int64(len(descItems)))
	}

	duration.End()
	duration.Stop()
}
```

改为：

```go
func (d *Destination[E]) do() {
	duration := timing.NewDuration()
	duration.Start()
	var descItems []E

	duration.Begin()
	for items := range d.itemsChan {
		if len(items) <= 0 {
			continue
		}

		descItems = append(descItems, items...)

		_len := len(descItems)
		_start := 0
		_end := 0
		_size := int(d.pageSize)
		for {
			_end = _start + _size
			if _end <= _len {
				err := d.dst.Exec(d, descItems[_start:_end])
				if err != nil {
					d.errChan <- fmt.Errorf("exec batch error; %w", err)
					return
				}

				d.state.AddAmount(int64(_size))

				_start += _size
			} else {
				descItems = slices.Clone(descItems[_start:])
				descItems = slices.Clip(descItems)

				break
			}
		}
	}

	if len(descItems) > 0 {
		err := d.dst.Exec(d, descItems)
		if err != nil {
			d.errChan <- fmt.Errorf("exec remaining error; %w", err)
			return
		}

		d.state.AddAmount(int64(len(descItems)))
	}

	duration.End()
	duration.Stop()
}
```

- [ ] **步骤 3：改造 Accept() 初始化 errChan**

将 `Accept()` 方法从：

```go
func (d *Destination[E]) Accept() (err error) {
	d.state.StatusAccept()
	d.state.DurationStart()

	if d.isTruncate {
		err = d.db.Truncate(d.table)
		if err != nil {
			return
		}
	}

	d.itemsChan = make(chan []E, d.concurrency)

	for i := 0; i < d.concurrency; i++ {
		d.doWg.Add(1)
		go func() {
			d.do()

			d.doWg.Done()
		}()
	}

	return nil
}
```

改为：

```go
func (d *Destination[E]) Accept() (err error) {
	d.state.StatusAccept()
	d.state.DurationStart()

	if d.isTruncate {
		err = d.db.Truncate(d.table)
		if err != nil {
			return
		}
	}

	d.itemsChan = make(chan []E, d.concurrency)
	d.errChan = make(chan error, d.concurrency)

	for i := 0; i < d.concurrency; i++ {
		d.doWg.Add(1)
		go func() {
			d.do()

			d.doWg.Done()
		}()
	}

	return nil
}
```

- [ ] **步骤 4：改造 Finish() 收集错误**

将 `Finish()` 方法从：

```go
func (d *Destination[E]) Finish() {
	d.doWg.Wait()

	d.state.StatusFinish()
	d.state.DurationStop()
}
```

改为：

```go
func (d *Destination[E]) Finish() {
	d.doWg.Wait()
	close(d.errChan)

	for err := range d.errChan {
		if d.firstErr == nil {
			d.firstErr = err
		}
	}

	d.state.StatusFinish()
	d.state.DurationStop()
}
```

- [ ] **步骤 5：添加 Err() 方法暴露错误**

在 `Destination` 上新增：

```go
func (d *Destination[E]) Err() error {
	return d.firstErr
}
```

- [ ] **步骤 6：运行全量编译验证**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go build ./...`
预期：编译成功

- [ ] **步骤 7：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add storage/database/destination/destination.go
git commit -m "refactor: database destination replaces panic with error channel"
```

---

### 任务 8：替换 panic 为 error 返回 — redis/destination

**文件：**
- 修改：`flow/storage/redis/destination/key.go`
- 修改：`flow/storage/redis/destination/hashes.go`
- 修改：`flow/storage/redis/destination/lists.go`
- 修改：`flow/storage/redis/destination/sets.go`
- 修改：`flow/storage/redis/destination/sorted_sets.go`

- [ ] **步骤 1：为 redis/destination key 添加 error channel**

将 `flow/storage/redis/destination/key.go` 中 `key` 结构体，从：

```go
type key[E storage.Entry] struct {
	storage.Storage
	concurrency int
	isTruncate  bool
	pageSize    int64
	keyName     string
	isDone      bool
	itemsChan   chan []E
	doWg        sync.WaitGroup
	client      *client.Redis
	keyer       keyer[E]
	state       *storage.State
}
```

改为：

```go
type key[E storage.Entry] struct {
	storage.Storage
	concurrency int
	isTruncate  bool
	pageSize    int64
	keyName     string
	isDone      bool
	itemsChan   chan []E
	doWg        sync.WaitGroup
	client      *client.Redis
	keyer       keyer[E]
	state       *storage.State
	errChan     chan error
	firstErr    error
}
```

- [ ] **步骤 2：改造 keyer 接口，accept 返回 error**

将 `keyer` 接口从：

```go
type keyer[E storage.Entry] interface {
	redis.Keyer
	accept(itemsChan <-chan []E, c *client.Redis, key string, pageSize int64)
	stateAmount() int64
}
```

改为：

```go
type keyer[E storage.Entry] interface {
	redis.Keyer
	accept(itemsChan <-chan []E, c *client.Redis, key string, pageSize int64) error
	stateAmount() int64
}
```

- [ ] **步骤 3：改造 key.Accept() 初始化 errChan**

将 `Accept()` 方法从：

```go
func (k *key[E]) Accept() error {
	k.state.StatusAccept()
	k.state.DurationStart()

	if k.isTruncate {
		_, err := k.keyer.Truncate(context.Background(), k.client, k.keyName)
		if err != nil {
			return err
		}
	}

	k.itemsChan = make(chan []E, k.concurrency)

	for i := 0; i < k.concurrency; i++ {
		k.doWg.Add(1)
		go func() {
			k.keyer.accept(k.itemsChan, k.client, k.keyName, k.pageSize)

			k.doWg.Done()
		}()
	}

	return nil
}
```

改为：

```go
func (k *key[E]) Accept() error {
	k.state.StatusAccept()
	k.state.DurationStart()

	if k.isTruncate {
		_, err := k.keyer.Truncate(context.Background(), k.client, k.keyName)
		if err != nil {
			return err
		}
	}

	k.itemsChan = make(chan []E, k.concurrency)
	k.errChan = make(chan error, k.concurrency)

	for i := 0; i < k.concurrency; i++ {
		k.doWg.Add(1)
		go func() {
			if err := k.keyer.accept(k.itemsChan, k.client, k.keyName, k.pageSize); err != nil {
				k.errChan <- err
			}

			k.doWg.Done()
		}()
	}

	return nil
}
```

- [ ] **步骤 4：改造 key.Finish() 收集错误**

将 `Finish()` 方法从：

```go
func (k *key[E]) Finish() {
	k.doWg.Wait()

	k.state.StatusFinish()
	k.state.DurationStop()
}
```

改为：

```go
func (k *key[E]) Finish() {
	k.doWg.Wait()
	close(k.errChan)

	for err := range k.errChan {
		if k.firstErr == nil {
			k.firstErr = err
		}
	}

	k.state.StatusFinish()
	k.state.DurationStop()
}
```

- [ ] **步骤 5：添加 key.Err() 方法**

```go
func (k *key[E]) Err() error {
	return k.firstErr
}
```

- [ ] **步骤 6：改造 hashes.accept() — 返回 error**

将 `flow/storage/redis/destination/hashes.go` 中 `accept` 方法，从：

```go
func (h *hashes) accept(itemsChan <-chan []storage.MapEntry, c *client.Redis, key string, pageSize int64) {
	ctx := context.Background()
	pipe := c.Pipeline()

	for items := range itemsChan {
		l := len(items)
		for i := 0; i < l; i += int(pageSize) {
			end := i + int(pageSize)
			if end > l {
				end = l
			}

			entries := items[i:end]
			for _, entry := range entries {
				for k, v := range entry {
					pipe.HMSet(ctx, key, k, v)
				}
			}

			_, err := pipe.Exec(ctx)
			if err != nil {
				panic(err)
			}
		}

		atomic.AddInt64(&h.amount, int64(l))
	}

	_ = pipe.Close()
}
```

改为：

```go
func (h *hashes) accept(itemsChan <-chan []storage.MapEntry, c *client.Redis, key string, pageSize int64) error {
	ctx := context.Background()
	pipe := c.Pipeline()

	for items := range itemsChan {
		l := len(items)
		for i := 0; i < l; i += int(pageSize) {
			end := i + int(pageSize)
			if end > l {
				end = l
			}

			entries := items[i:end]
			for _, entry := range entries {
				for k, v := range entry {
					pipe.HMSet(ctx, key, k, v)
				}
			}

			_, err := pipe.Exec(ctx)
			if err != nil {
				_ = pipe.Close()
				return fmt.Errorf("hashes accept exec error; %w", err)
			}
		}

		atomic.AddInt64(&h.amount, int64(l))
	}

	_ = pipe.Close()
	return nil
}
```

注意：需要在 `hashes.go` 顶部添加 `"fmt"` 到 import。

- [ ] **步骤 7：改造 lists.accept() — 返回 error**

将 `flow/storage/redis/destination/lists.go` 中 `accept` 方法，从：

```go
func (h *lists) accept(itemsChan <-chan []string, c *client.Redis, key string, pageSize int64) {
	ctx := context.Background()
	for items := range itemsChan {
		l := len(items)
		for i := 0; i < l; i += int(pageSize) {
			end := i + int(pageSize)
			if end > l {
				end = l
			}

			entries := items[i:end]

			entriesAny := make([]any, 0, end-i)
			for _, entry := range entries {
				entriesAny = append(entriesAny, entry)
			}

			_, err := c.LPush(ctx, key, entriesAny...).Result()
			if err != nil {
				panic(err)
			}
		}

		atomic.AddInt64(&h.amount, int64(l))
	}
}
```

改为：

```go
func (h *lists) accept(itemsChan <-chan []string, c *client.Redis, key string, pageSize int64) error {
	ctx := context.Background()
	for items := range itemsChan {
		l := len(items)
		for i := 0; i < l; i += int(pageSize) {
			end := i + int(pageSize)
			if end > l {
				end = l
			}

			entries := items[i:end]

			entriesAny := make([]any, 0, end-i)
			for _, entry := range entries {
				entriesAny = append(entriesAny, entry)
			}

			_, err := c.LPush(ctx, key, entriesAny...).Result()
			if err != nil {
				return fmt.Errorf("lists accept lpush error; %w", err)
			}
		}

		atomic.AddInt64(&h.amount, int64(l))
	}
	return nil
}
```

注意：需要在 `lists.go` 顶部添加 `"fmt"` 到 import。

- [ ] **步骤 8：改造 sets.accept() — 返回 error**

将 `flow/storage/redis/destination/sets.go` 中 `accept` 方法，从：

```go
func (h *sets) accept(itemsChan <-chan []string, c *client.Redis, key string, pageSize int64) {
	ctx := context.Background()
	for items := range itemsChan {
		l := len(items)
		for i := 0; i < l; i += int(pageSize) {
			end := i + int(pageSize)
			if end > l {
				end = l
			}

			entries := items[i:end]

			entriesAny := make([]any, 0, end-i)
			for _, entry := range entries {
				entriesAny = append(entriesAny, entry)
			}

			_, err := c.SAdd(ctx, key, entriesAny...).Result()
			if err != nil {
				panic(err)
			}
		}

		atomic.AddInt64(&h.amount, int64(l))
	}
}
```

改为：

```go
func (h *sets) accept(itemsChan <-chan []string, c *client.Redis, key string, pageSize int64) error {
	ctx := context.Background()
	for items := range itemsChan {
		l := len(items)
		for i := 0; i < l; i += int(pageSize) {
			end := i + int(pageSize)
			if end > l {
				end = l
			}

			entries := items[i:end]

			entriesAny := make([]any, 0, end-i)
			for _, entry := range entries {
				entriesAny = append(entriesAny, entry)
			}

			_, err := c.SAdd(ctx, key, entriesAny...).Result()
			if err != nil {
				return fmt.Errorf("sets accept sadd error; %w", err)
			}
		}

		atomic.AddInt64(&h.amount, int64(l))
	}
	return nil
}
```

注意：需要在 `sets.go` 顶部添加 `"fmt"` 到 import。

- [ ] **步骤 9：改造 sorted_sets.accept() — 返回 error**

将 `flow/storage/redis/destination/sorted_sets.go` 中 `accept` 方法，从：

```go
func (h *sortedSets) accept(itemsChan <-chan []storage.ScoreMap, c *client.Redis, key string, pageSize int64) {
	ctx := context.Background()
	pipe := c.Pipeline()

	for items := range itemsChan {
		l := len(items)
		for i := 0; i < l; i += int(pageSize) {
			end := i + int(pageSize)
			if end > l {
				end = l
			}

			entries := items[i:end]
			for _, entry := range entries {
				for k, v := range entry {
					pipe.ZAdd(ctx, key, &redis2.Z{
						Score:  v,
						Member: k,
					})
				}
			}

			_, err := pipe.Exec(ctx)
			if err != nil {
				panic(err)
			}
		}

		atomic.AddInt64(&h.amount, int64(l))
	}

	_ = pipe.Close()
}
```

改为：

```go
func (h *sortedSets) accept(itemsChan <-chan []storage.ScoreMap, c *client.Redis, key string, pageSize int64) error {
	ctx := context.Background()
	pipe := c.Pipeline()

	for items := range itemsChan {
		l := len(items)
		for i := 0; i < l; i += int(pageSize) {
			end := i + int(pageSize)
			if end > l {
				end = l
			}

			entries := items[i:end]
			for _, entry := range entries {
				for k, v := range entry {
					pipe.ZAdd(ctx, key, &redis2.Z{
						Score:  v,
						Member: k,
					})
				}
			}

			_, err := pipe.Exec(ctx)
			if err != nil {
				_ = pipe.Close()
				return fmt.Errorf("sorted sets accept exec error; %w", err)
			}
		}

		atomic.AddInt64(&h.amount, int64(l))
	}

	_ = pipe.Close()
	return nil
}
```

注意：需要在 `sorted_sets.go` 顶部添加 `"fmt"` 到 import。

- [ ] **步骤 10：运行全量编译验证**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go build ./...`
预期：编译成功

- [ ] **步骤 11：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add storage/redis/destination/
git commit -m "refactor: redis destination replaces panic with error returns"
```

---

### 任务 9：替换 panic 为 error 返回 — file source/destination

**文件：**
- 修改：`flow/storage/file/source/line.go`
- 修改：`flow/storage/file/destination/line.go`

- [ ] **步骤 1：改造 file/source/line.go — Scan() 中的 panic**

将 `Scan()` 方法中：

```go
err := l.b.Err()
if err != nil {
    panic(l.Title() + err.Error())
}
```

改为：

```go
err := l.b.Err()
if err != nil {
    close(l.itemsChan)
    l.state.DurationStop()
    return fmt.Errorf("file source scan error; %w", err)
}
```

注意：需要在 `line.go` 顶部添加 `"fmt"` 到 import（如尚未存在）。

- [ ] **步骤 2：改造 file/destination/line.go — Receive() 中的 panic**

将 `Receive()` 方法从：

```go
func (l *Line) Receive(items []string) {
	for k := range items {
		l.state.AddAmount(1)
		_, err := l.b.WriteString(items[k] + "\n")
		if err != nil {
			panic(err)
		}
	}

	err := l.b.Flush()
	if err != nil {
		panic(err)
	}
}
```

改为：

```go
func (l *Line) Receive(items []string) error {
	for k := range items {
		l.state.AddAmount(1)
		_, err := l.b.WriteString(items[k] + "\n")
		if err != nil {
			return fmt.Errorf("file destination write error; %w", err)
		}
	}

	err := l.b.Flush()
	if err != nil {
		return fmt.Errorf("file destination flush error; %w", err)
	}

	return nil
}
```

注意：这会改变 `Destinationer` 接口签名。需要同步更新 `storage/destination.go`：

将：

```go
type Destinationer[E Entry] interface {
	Accept() error
	Receive([]E)
	Done()
	Finish()
	Close() error
	Summary() []string
	State() []string
}
```

改为：

```go
type Destinationer[E Entry] interface {
	Accept() error
	Receive([]E) error
	Done()
	Finish()
	Close() error
	Summary() []string
	State() []string
}
```

同时需要更新所有实现了 `Destinationer` 的类型的 `Receive` 方法签名：

- `mock/destination/destination.go` 的 `Destination.Receive` 改为 `func (d *Destination[E]) Receive(items []E) error { d.itemsChan <- items; return nil }`
- `database/destination/destination.go` 的 `Destination.Receive` 改为 `func (d *Destination[E]) Receive(items []E) error { d.itemsChan <- items; return nil }`
- `redis/destination/key.go` 的 `key.Receive` 改为 `func (k *key[E]) Receive(items []E) error { k.itemsChan <- items; return nil }`

- [ ] **步骤 3：运行全量编译验证**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go build ./...`
预期：编译成功

- [ ] **步骤 4：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add storage/file/ storage/destination.go storage/mock/destination/destination.go storage/database/destination/destination.go storage/redis/destination/key.go
git commit -m "refactor: file source/destination replaces panic with error returns; update Destinationer interface"
```

---

### 任务 10：Flow.transport() 增加 goroutine 错误恢复

**文件：**
- 修改：`flow/flow/flow.go`

- [ ] **步骤 1：为 Flow 添加 error 收集字段**

在 `Flow` 结构体中添加：

```go
type Flow[E storage.Entry] struct {
	source        storage.Sourceor[E]
	refreshOutput *output.Refresh
	actions       []action.Actor[E]
	stateInterval time.Duration
	firstErr      error
	errOnce       sync.Once
}
```

注意：需要在 import 中添加 `"sync"`。

- [ ] **步骤 2：改造 transport() 方法**

将 `transport()` 方法从：

```go
func (f *Flow[E]) transport() {
	needCopy := false
	if len(f.actions) > 1 {
		needCopy = true
	}

	go func() {
		for {
			items, ok := <-f.source.ReceiveChan()
			if !ok {
				break
			}

			for _, a := range f.actions {
				if needCopy {
					newItems := f.source.Copy(items)
					a.Receive(newItems)
				} else {
					a.Receive(items)
				}
			}
		}

		f.actionsDone()
	}()
}
```

改为：

```go
func (f *Flow[E]) transport() {
	needCopy := false
	if len(f.actions) > 1 {
		needCopy = true
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.errOnce.Do(func() {
					f.firstErr = fmt.Errorf("transport panic: %v", r)
				})
			}
		}()

		for {
			items, ok := <-f.source.ReceiveChan()
			if !ok {
				break
			}

			for _, a := range f.actions {
				if needCopy {
					newItems := f.source.Copy(items)
					a.Receive(newItems)
				} else {
					a.Receive(items)
				}
			}
		}

		f.actionsDone()
	}()
}
```

- [ ] **步骤 3：改造 finish() 方法检查错误**

将 `finish()` 方法从：

```go
func (f *Flow[E]) finish() error {
	err := f.actionsFinish()
	f.refreshOutput.Stop()
	f.actionsOutput()

	if err != nil {
		return fmt.Errorf("actions finish error; %w", err)
	} else {
		return nil
	}
}
```

改为：

```go
func (f *Flow[E]) finish() error {
	err := f.actionsFinish()
	f.refreshOutput.Stop()
	f.actionsOutput()

	if f.firstErr != nil {
		return f.firstErr
	}

	if err != nil {
		return fmt.Errorf("actions finish error; %w", err)
	}

	return nil
}
```

- [ ] **步骤 4：运行全量编译验证**

运行：`cd /opt/coding/github.com/go-toolkit/flow && go build ./...`
预期：编译成功

- [ ] **步骤 5：Commit**

```bash
cd /opt/coding/github.com/go-toolkit/flow
git add flow/flow.go
git commit -m "refactor: Flow.transport() adds panic recovery and error propagation"
```

---

## 自检

### 1. 规格覆盖度

| 规格需求 | 对应任务 |
|---------|---------|
| BUG-1: Sets.Type() 错误 | 任务 1 |
| BUG-2: SortedSets.Type() 错误 | 任务 2 |
| BUG-3: Mock pageSize 校验错误 | 任务 3 |
| BUG-4: Flow 未校验 Source nil | 任务 4 |
| IMPL-1: 拼写错误 | 任务 5 |
| ERR-1: log.Fatal 在库中 | 任务 6 |
| ERR-2: panic 替换 | 任务 7, 8, 9 |
| ERR-3: transport goroutine 错误处理 | 任务 10 |

**未覆盖（属于阶段三/四，需要更大范围 API 变更）：**
- ARCH-1: context.Context 支持
- ARCH-2: 空包清理
- ARCH-3: tool/copy.go 位置
- ARCH-4: Moder 重命名
- ARCH-5: Destinationer 未被 Flow 使用
- ARCH-6: 优雅停机
- IMPL-2~8: 输出统一、Copy 一致性、批处理优化等

### 2. 占位符扫描

无占位符。所有步骤包含完整代码。

### 3. 类型一致性

- 任务 9 修改了 `Destinationer.Receive` 签名为 `Receive([]E) error`，所有实现类在任务 9 中同步更新
- 任务 8 修改了 `keyer.accept` 签名为 `accept(...) error`，所有实现类在任务 8 中同步更新
- 任务 7 新增 `errChan` 和 `firstErr` 字段，与任务 8 的 `key` 结构体保持一致的模式
