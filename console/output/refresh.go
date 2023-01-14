package output

import (
	"fmt"
	"time"
)

type Refresh struct {
	MultilineText
	currentLine      int
	lastRefreshTime  time.Time
	ticker           *time.Ticker
	isInterval       bool
	intervalDuration time.Duration
}

func NewRefresh() *Refresh {
	r := &Refresh{}
	r.intervalDuration = time.Millisecond * 200

	r.content = make([]string, 0)
	r.ticker = time.NewTicker(r.intervalDuration)
	r.lastRefreshTime = time.Now()

	return r
}

// Start 开始刷新输出，定时刷新内容到输出
func (r *Refresh) Start() {
	r.interval()
}

// Stop 结束刷新输出
// 输出内容
// 清空内容
// 停止定时输出
func (r *Refresh) Stop() {
	r.refresh()
	r.flushContent()
	r.ticker.Stop()
}

func (r *Refresh) CleanAndStart() {
	r.Start()
	r.ticker.Reset(r.intervalDuration)
}

func (r *Refresh) refresh() {
	contentLen := len(r.content)
	content := make([]string, contentLen, contentLen)
	copy(content, r.content)

	if r.currentLine == 0 {
		for i := 0; i < contentLen; i++ {
			fmt.Printf("%c[1B\r\n", 0x1B)
			r.currentLine += 1
		}
	} else if contentLen > r.currentLine {
		for i := 0; i < contentLen-r.currentLine; i++ {
			fmt.Printf("%c[1B\r\n", 0x1B)
			r.currentLine += 1
		}
	}

	fmt.Printf("%c[%dA", 0x1B, r.currentLine+1)
	r.currentLine = 0

	for _, s := range content {
		fmt.Printf("%c[1B\r%c[K%c[1;40;32m%s%c[0m", 0x1B, 0x1B, 0x1B, s, 0x1B)
		r.currentLine += 1
	}

	fmt.Printf("%c[1B\r", 0x1B)

	r.lastRefreshTime = time.Now()
}

func (r *Refresh) interval() {
	if r.isInterval {
		return
	}

	r.isInterval = true

	go func() {
		for {
			if _, ok := <-r.ticker.C; ok {
				r.refresh()
			} else {
				break
			}
		}
	}()
}

func (r *Refresh) flushContent() {
	r.content = make([]string, 0)
	r.currentLine = 0
}

func (r *Refresh) saveCursorPosition() {
	fmt.Printf("%c[s", 0x1B)
}

func (r *Refresh) restoreCursorPosition() {
	fmt.Printf("%c[u", 0x1B)
}
