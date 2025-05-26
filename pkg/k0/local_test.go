package k0

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/internal/types"
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
			_ = fs.SaveKeylog("device1", 1, 1)
		}()
	}

	wg.Wait() // wait for all goroutines to finish
	t.Log("done")

	// check response
	resultKeylogs := fs.dataFile.Keylogs["device1"][1][1]
	if resultKeylogs != int64(numGoRoutines) {
		t.Errorf("expected %d, got %d", numGoRoutines, resultKeylogs)
	}
}

func TestPeriodicSave(t *testing.T) {
	before := runtime.NumGoroutine()
	// create temp file
	tf, err := os.MkdirTemp("", "local_test")
	if err != nil {
		t.Fatal(err)
	}
	filename := fmt.Sprintf("device_%d", rand.Int())
	filepath := path.Join(tf, filename)
	periodicSave := 200 * time.Millisecond
	config := ConfigStorage{
		FileOutput:   filepath,
		PeriodicSave: types.Duration{Duration: periodicSave},
	}
	ctx, cancel := context.WithCancel(context.Background())
	fs := MustGetNewFileStorage(ctx, config)

	// save keylog
	err = fs.SaveKeylog("device1", 1, 1)
	if err != nil {
		t.Fatal(err)
	}

	// wait for periodic save
	time.Sleep(2 * periodicSave)

	// open content of fd.Name() file

	dataFile := newDataFile()
	err = ParseFromFile(filepath, dataFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dataFile: %+v\n", dataFile)
	if dataFile.Keylogs["device1"][1][1] != 1 {
		t.Fatal("expected 1 keylog")
	}
	cancel() // close file storage

	time.Sleep(2 * time.Second)
	after := runtime.NumGoroutine()
	if after > before {
		t.Fatalf("Goroutines leak. Before: %d, After: %d", before, after)
	}

	// cleanup file
	err = os.Remove(filepath)
	if err != nil {
		t.Fatal(err)
	}
}
