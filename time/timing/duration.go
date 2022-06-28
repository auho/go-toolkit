package timing

import (
	"fmt"
	"time"
)

// DefaultDuration 默认 Duration
var DefaultDuration Duration

type Duration struct {
	startTime  time.Time
	finishTime time.Time
	beginTime  time.Time
	endTime    time.Time
}

func NewDuration() *Duration {
	return &Duration{}
}

// Start 启动 Duration
func (t *Duration) Start() {
	t.startTime = time.Now()
}

// Stop 停止 Duration
func (t *Duration) Stop() {
	t.finishTime = time.Now()
}

// Begin 开始一小段计时
func (t *Duration) Begin() {
	t.beginTime = time.Now()
}

// End 结束一小段计时
func (t *Duration) End() {
	t.endTime = time.Now()
}

// SubBegin 从上一个小段 Begin 到现在的时长
func (t *Duration) SubBegin() time.Duration {
	return time.Now().Sub(t.beginTime)
}

// SubStart 从启动到现在的时长
func (t *Duration) SubStart() time.Duration {
	return time.Now().Sub(t.startTime)
}

// StringStartToNowSeconds 从启动到现在的秒数
func (t *Duration) StringStartToNowSeconds() string {
	return fmt.Sprintf("duration %f 秒", t.SubStart().Seconds())
}
