package k0

import (
	"runtime"
	"testing"
	"time"
)

func CheckGoroutineLeak(t *testing.T, before int) {
	time.Sleep(2 * time.Second)
	after := runtime.NumGoroutine()
	if after > before {
		t.Fatalf("Goroutines leak. Before: %d, After: %d", before, after)
	}
}
