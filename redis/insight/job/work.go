package job

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/auho/go-toolkit/farmtools/sort/maps"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
	"github.com/auho/go-toolkit/redis/client"
)

type Worker[E storage.Entry] interface {
	task.Work[E]
	WithRedisSource(*client.Redis)
	WithTimeMark(string)
}

type Work struct {
	task.Task

	Client   *client.Redis
	TimeMark string

	patternCounters *PatternCounters
}

func (w *Work) WithRedisSource(c *client.Redis) {
	w.Client = c
}

func (w *Work) WithTimeMark(s string) {
	w.TimeMark = s
}

func (w *Work) RegisterPatternCounter(names ...string) {
	w.patternCounters = NewPatternCounters(names...)
}

func (w *Work) DonePatternCounter() {
	w.patternCounters.Done()
}

func (w *Work) PrintlnPatternCounter(lessNum int) {
	rets := w.patternCounters.Result()
	for _, ret := range rets {
		w.Println(ret.Name + ": ")
		w.PrintlnCounter(ret.Data, lessNum)
		w.Println("")
	}
}

func (w *Work) IncrementPrintlnCounter(name, key string, lessNum int) {
	w.patternCounters.Increment(name, key, lessNum)
}

// PrintlnCounter 打印计数
// m 计数 map
// lessAmount 小于的次数不打印
func (w *Work) PrintlnCounter(m map[string]int, lessNum int) {
	lessAmount := 0
	for k, v := range m {
		if v <= lessNum {
			lessAmount += 1
		} else {
			w.Println(k, "\t", v)
		}
	}

	title := fmt.Sprintf("*key less %d amount*", lessNum)
	w.Println(title, "\t", lessAmount)
}

func (w *Work) SortedAndLogToFile(fileName string, m map[string]int, lessNum int) {
	_, err := os.Stat("data")
	if err != nil {
		if !os.IsExist(err) {
			err = os.Mkdir("data", 0744)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	filePath := fmt.Sprintf("data/%s_%s.log", w.TimeMark, fileName)
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	w.SortedAndLog(f, m, lessNum)

	defer func() {
		_ = f.Close()
	}()
}

func (w *Work) SortedAndLog(i io.WriteCloser, m map[string]int, lessNum int) {
	keys, values := maps.SortValueDesc(m)

	var err error

	lessAmount := 0
	for _i, _k := range keys {
		if values[_i] <= lessNum {
			lessAmount += 1
		} else {
			_, err = i.Write([]byte(fmt.Sprintln(_k, "\t", values[_i])))
			if err != nil {
				w.Println(err)
			}
		}
	}

	_, err = i.Write([]byte(fmt.Sprintf("-key:less:%d:amount-%s%d\n", lessNum, "\t", lessAmount)))
	if err != nil {
		w.Println(err)
	}
}
