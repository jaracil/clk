// +build linux

package clk

import (
	"syscall"
	"time"
	"unsafe"
)

// Linux clock mode constants
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

func lap() time.Duration {
	var ts syscall.Timespec
	syscall.Syscall(syscall.SYS_CLOCK_GETTIME, MONOTONIC, uintptr(unsafe.Pointer(&ts)), 0)
	return time.Duration(ts.Sec)*time.Second + time.Duration(ts.Nsec)
}
