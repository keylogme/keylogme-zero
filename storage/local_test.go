package storage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/types"
)

// test to replicate the issue "concurrent map writes"
// when using Datafile.AddKeyLog
func TestMultipleSaves(t *testing.T) {
	var wg sync.WaitGroup
	// create temp file
	config := ConfigStorage{
		FileOutput:   "test.json",
		PeriodicSave: types.Duration{Duration: time.Duration(1 * time.Second)},
	}
	fs := MustGetNewFileStorage(context.TODO(), config)

	numGoRoutines := 10000
	// save keylog in multiple goroutines to replicate the issue
	for range numGoRoutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = fs.SaveKeylog("device1", 1)
		}()
	}

	wg.Wait() // wait for all goroutines to finish
	t.Log("done")

	// check response
	resultKeylogs := fs.dataFile.Keylogs["device1"][1]
	if resultKeylogs != int64(numGoRoutines) {
		t.Errorf("expected %d, got %d", numGoRoutines, resultKeylogs)
	}
}
