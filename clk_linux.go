package clk

import (
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// Clock mode constants
const (
	REALTIME           = 0
	MONOTONIC          = 1
	PROCESS_CPUTIME_ID = 2
	THREAD_CPUTIME_ID  = 3
	MONOTONIC_RAW      = 4
	REALTIME_COARSE    = 5
	MONOTONIC_COARSE   = 6
	BOOTTIME           = 7
	REALTIME_ALARM     = 8
	BOOTTIME_ALARM     = 9
)

// Lap returns the time elapsed since clockid started (look at clock mode constants).
func Lap(clockid uintptr) time.Duration {
	var ts syscall.Timespec
	syscall.Syscall(syscall.SYS_CLOCK_GETTIME, clockid, uintptr(unsafe.Pointer(&ts)), 0)
	return time.Duration(ts.Sec)*1e9 + time.Duration(ts.Nsec)
}

// Cache type contains information for caching system clock access. It reduces syscalls overhead.
type Cache struct {
	lock    sync.Mutex
	clkid   uintptr
	last    time.Duration
	refresh time.Duration
	fresh   bool
}

// NewCache returns new Cache type.
//   clkid: (look at clock mode constants)
//   refresh: cache refresh period.
func NewCache(clkid uintptr, refresh time.Duration) *Cache {
	return &Cache{clkid: clkid, refresh: refresh}
}

// Lap returns the time elapsed since cached clock started (cached value).
func (p *Cache) Lap() time.Duration {
	p.lock.Lock()
	if p.fresh {
		last := p.last
		p.lock.Unlock()
		return last
	}
	last := Lap(p.clkid)
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
