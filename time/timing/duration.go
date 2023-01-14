package timing

import (
	"fmt"
	"time"
)

// DefaultDuration 默认 Duration
var DefaultDuration Duration

type Duration struct {
	startTime time.Time
	stopTime  time.Time
	beginTime time.Time
	endTime   time.Time
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
	t.stopTime = time.Now()
}

// Begin 开始一小段计时
func (t *Duration) Begin() {
	t.beginTime = time.Now()
}

// End 结束一小段计时
func (t *Duration) End() {
	t.endTime = time.Now()
}

// SubBegin 从上一个小段 Begin 到结束的时长
func (t *Duration) SubBegin() time.Duration {
	if t.endTime == (time.Time{}) {
		return time.Now().Sub(t.beginTime)
	} else {
		return t.endTime.Sub(t.beginTime)
	}
}

// SubStart 从启动到停止的时长
func (t *Duration) SubStart() time.Duration {
	if t.stopTime == (time.Time{}) {
		return time.Now().Sub(t.startTime)
	} else {
		return t.stopTime.Sub(t.startTime)
	}
}

// StringStartToStop 从启动到停止的时长
func (t *Duration) StringStartToStop() string {
	return t.stringPretty(t.SubStart())
}

// StringBeginToEnd 从开始到结束的时长
func (t *Duration) StringBeginToEnd() string {
	return t.stringPretty(t.SubBegin())
}

func (t *Duration) stringPretty(d time.Duration) string {
	seconds := d.Seconds()
	if seconds < 60 {
		return fmt.Sprintf("%f 秒", seconds)
	} else if seconds < 3600 {
		m := int64(seconds / 60)
		s := int64(seconds) % 60

		return fmt.Sprintf("%d 分 %d 秒", m, s)
	} else if seconds < 86400 {
		h := int64(seconds / 3600)
		m := (int64(seconds) % 3600) / 60
		s := int64(seconds) % 60

		return fmt.Sprintf("%d 小时 %d 分 %d 秒", h, m, s)
	} else {
		d := int64(seconds / 86400)
		h := (int64(seconds) % 86400) / 3600
		m := (int64(seconds) % 3600) / 60

		return fmt.Sprintf("%d 天 %d 小时 %d 分 ", d, h, m)
	}
}
