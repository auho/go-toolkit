package output

import (
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	r := NewRefresh()

	r.Print(1, "test print")
	r.Print(2, "2:0")
	r.Print(3, "3:0")
	go func() {
		r.Start()
	}()

	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
			r.Print(4, time.Now().String())
		}
	}()

	time.Sleep(time.Second * 5)

	ticker.Stop()
}

func TestCoverAll(t *testing.T) {
	r := NewRefresh()

	r.Print(1, "test cover all")
	r.Print(2, "2:0")
	r.Print(3, "3:0")
	go func() {
		r.Start()
	}()

	ticker := time.NewTicker(time.Second)
	sss := []string{
		"test cover all new",
		"2:0 new",
		"3:0 new",
		"",
	}
	go func() {
		for range ticker.C {
			sss[3] = time.Now().String()
			r.CoverAll(sss)
		}
	}()

	time.Sleep(time.Second * 5)

	ticker.Stop()
}
