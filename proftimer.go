package proftimer

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

type timer struct {
	start time.Time
	count uint64
	total time.Duration
}

var mu sync.Mutex
var timers map[string]*timer = map[string]*timer{}

// Resume all the named the timers. All of them record the same starting time.
func Resume(names ...string) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	for _, name := range names {
		t := timers[name]
		if t == nil {
			t = &timer{}
			timers[name] = t
		}
		t.start = now
	}
}

// Accum accumulates the current elapsed time from the start of each timer up
// to the moment Accum is called. Inexistant timers are ignored.
func Accum(names ...string) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	for _, name := range names {
		t := timers[name]
		if t == nil {
			continue
		}

		t.total = t.total + now.Sub(t.start)
	}
}

// Reset removes the named timers, effectively cleaning up after them.
func Reset(names ...string) {
	mu.Lock()
	defer mu.Unlock()

	for _, name := range names {
		delete(timers, name)
	}
}

// Event increases the event count on all named timers. Inexistant timers are
// ignored.
func Event(names ...string) {
	mu.Lock()
	defer mu.Unlock()

	for _, name := range names {
		t := timers[name]
		if t == nil {
			continue
		}

		t.count += 1
	}

}

// report the given timer to the given IO writer. The lock MUST be held.
func report(name string, w io.Writer) {
	t := timers[name]
	if t == nil {
		t = new(timer)
	}
	fmt.Fprintf(w, "%20s: %s\n", name, t.total)
}

// Report all the given timers to the given io writer. If no names are specified,
// all timers are reported.
func Report(w io.Writer, names ...string) {
	mu.Lock()
	defer mu.Unlock()

	if len(names) == 0 {
		names = make([]string, 0, len(timers))
		for n := range timers {
			names = append(names, n)
		}
		sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })
	}

	for _, name := range names {
		report(name, w)
	}
}
