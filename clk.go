package clk

import (
	"sync"
	"time"
)

// Cache type contains information for caching system clock access. It reduces syscalls overhead.
type Cache struct {
	lock    sync.Mutex
	last    time.Duration
	refresh time.Duration
	fresh   bool
}

// Lap returns the time elapsed since monotonic clock started.
func Lap() time.Duration {
	return lap()
}

// NewCache returns new Cache type.
//   clkid: (look at clock mode constants)
//   refresh: cache refresh period.
func NewCache(refresh time.Duration) *Cache {
	return &Cache{refresh: refresh}
}

// Lap returns the time elapsed since cached clock started (cached value).
func (p *Cache) Lap() time.Duration {
	p.lock.Lock()
	if p.fresh {
		last := p.last
		p.lock.Unlock()
		return last
	}
	last := Lap()
	p.last = last
	p.fresh = true
	time.AfterFunc(p.refresh, p.cancel)
	p.lock.Unlock()
	return last
}

func (p *Cache) cancel() {
	p.lock.Lock()
	p.fresh = false
	p.lock.Unlock()
}
