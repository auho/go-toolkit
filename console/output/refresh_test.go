package output

import (
	"fmt"
	"testing"
	"time"
)

func TestOne(t *testing.T) {
	F := NewRefresh()

	fmt.Println("aaa")

	F.Print(1, "0:0")
	F.Print(2, "1:0")
	F.Print(3, "2:0")
	go func() {
		F.Start()
	}()

	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
			F.Print(4, time.Now().String())
		}
	}()

	time.Sleep(time.Second * 5)

	ticker.Stop()
}
