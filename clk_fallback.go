// +build !linux

// Fallback fake monotonic clock

package clk

import (
	"time"
)

var timeRef = time.Now()

func lap() time.Duration {
	return time.Now().Sub(timeRef)
}
