package job

import (
	"log"
	"sync"
)

type Result struct {
	Name string
	Data map[string]int
}

type item struct {
	key string
	num int
}

type PatternCounters struct {
	names []string
	list  map[string]patternCounter

	wg sync.WaitGroup
}

func NewPatternCounters(names ...string) *PatternCounters {
	pcs := &PatternCounters{
		list: make(map[string]patternCounter),
	}

	pcs.prepare(names...)

	return pcs
}

func (pcs *PatternCounters) prepare(names ...string) {
	for _, name := range names {
		pcs.names = append(pcs.names, name)
		pc := newPatternCounter(name)

		pcs.wg.Add(1)
		go func() {
			pc.receive()
			pcs.wg.Done()
		}()

		pcs.list[name] = pc
	}
}

func (pcs *PatternCounters) Increment(name string, key string, num int) {
	pc, ok := pcs.list[name]
	if !ok {
		log.Fatalf("name[%s] not found in patternCounters", name)
	}

	pc.increment(key, num)
}

func (pcs *PatternCounters) Done() {
	for _, pc := range pcs.list {
		pc.done()
	}

	pcs.wg.Wait()
}

func (pcs *PatternCounters) Result() []Result {
	var rets []Result
	for _, name := range pcs.names {
		pc := pcs.list[name]
		rets = append(rets, Result{
			Name: name,
			Data: pc.statistics,
		})
	}

	return rets
}

type patternCounter struct {
	name       string
	statistics map[string]int
	channel    chan item
}

func newPatternCounter(name string) patternCounter {
	return patternCounter{
		name:       name,
		statistics: make(map[string]int),
		channel:    make(chan item, 10),
	}
}

func (pc *patternCounter) increment(key string, num int) {
	pc.channel <- item{key, num}
}

func (pc *patternCounter) receive() {
	for _item := range pc.channel {
		pc.statistics[_item.key] += _item.num
	}
}

func (pc *patternCounter) done() {
	close(pc.channel)
}
