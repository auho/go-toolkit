package timing

import (
	"fmt"
	"testing"
	"time"
)

func TestOne(t *testing.T) {
	d := DefaultDuration

	d.Start()
	d.Begin()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 100)
		fmt.Println(d.StringBeginToEnd())
		time.Sleep(time.Millisecond * 100)
	}

	d.End()
	fmt.Println(d.StringBeginToEnd())

	fmt.Println(d.StringStartToStop())

	DefaultDuration.Stop()

	fmt.Println(d.StringStartToStop())
}

func TestStringPretty(t *testing.T) {
	d := DefaultDuration
	d.Start()

	fmt.Println(d.stringPretty(time.Second * 59))
	fmt.Println(d.stringPretty(time.Second * 3599))
	fmt.Println(d.stringPretty(time.Second * 86399))
	fmt.Println(d.stringPretty(time.Second * 9999999))
}
