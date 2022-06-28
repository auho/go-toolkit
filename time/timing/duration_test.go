package timing

import (
	"fmt"
	"testing"
	"time"
)

func TestOne(t *testing.T) {
	d := DefaultDuration

	d.Start()
	time.Sleep(time.Second)
	d.Begin()
	time.Sleep(time.Second)
	d.End()
	fmt.Println(d.SubBegin().String())
	fmt.Println(d.SubStart().String())
	fmt.Println(d.StringStartToNowSeconds())

	DefaultDuration.End()
}
