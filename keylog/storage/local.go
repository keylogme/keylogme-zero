package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type Storage interface {
	SaveKeylog(deviceId int64, keycode uint16) error
	SaveShortcut(deviceId int64, shortcutId int64) error
}

type FileStorage struct {
	fname     string
	keylogs   map[int64]map[uint16]int64 // deviceId - keycode - counter
	shortcuts map[int64]map[int64]int64  // deviceId - shortcutId - counter
}

type DataFile struct {
	Keylogs   map[int64]map[uint16]int64 `json:"keylogs"`
	Shortcuts map[int64]map[int64]int64  `json:"shortcuts"`
}

func NewFileStorage(ctx context.Context, fname string) *FileStorage {
	ffs := &FileStorage{
		fname:     fname,
		keylogs:   map[int64]map[uint16]int64{},
		shortcuts: map[int64]map[int64]int64{},
	}
	// go func(ctx context.Context) {
	// 	ffs.savingInBackground(ctx)
	// }(ctx)
	go ffs.savingInBackground(ctx)
	return ffs
}

func (f *FileStorage) SaveKeylog(deviceId int64, keycode uint16) error {
	if _, ok := f.keylogs[deviceId]; !ok {
		f.keylogs[deviceId] = map[uint16]int64{}
	}
	if _, ok := f.keylogs[deviceId][keycode]; !ok {
		f.keylogs[deviceId][keycode] = 0
	}
	f.keylogs[deviceId][keycode] += 1
	return nil
}

func (f *FileStorage) SaveShortcut(deviceId int64, shortcutId int64) error {
	if _, ok := f.shortcuts[deviceId]; !ok {
		f.shortcuts[deviceId] = map[int64]int64{}
	}
	if _, ok := f.shortcuts[deviceId][shortcutId]; !ok {
		f.shortcuts[deviceId][shortcutId] = 0
	}
	f.shortcuts[deviceId][shortcutId] += 1
	return nil
}

func (f *FileStorage) prepareDataToSave() (DataFile, error) {
	content, err := os.ReadFile(f.fname)
	if err != nil {
		return DataFile{}, err
	}
	dataFile := new(DataFile)
	err = json.Unmarshal(content, dataFile)
	if err != nil {
		return DataFile{}, err
	}

	for kId := range f.keylogs {
		for keycode := range f.keylogs[kId] {
			if _, ok := dataFile.Keylogs[kId][keycode]; ok {
				dataFile.Keylogs[kId][keycode] += f.keylogs[kId][keycode]
				continue
			}
			if _, ok := dataFile.Keylogs[kId]; !ok {
				dataFile.Keylogs[kId] = map[uint16]int64{}
			}
			if _, ok := dataFile.Keylogs[kId][keycode]; !ok {
				dataFile.Keylogs[kId][keycode] = f.keylogs[kId][keycode]
			}
		}
	}
	for kId := range f.shortcuts {
		for scId := range f.shortcuts[kId] {
			if _, ok := dataFile.Shortcuts[kId][scId]; ok {
				dataFile.Shortcuts[kId][scId] += f.shortcuts[kId][scId]
				continue
			}
			if _, ok := dataFile.Shortcuts[kId]; !ok {
				dataFile.Shortcuts[kId] = map[int64]int64{}
			}
			if _, ok := dataFile.Shortcuts[kId][kId]; !ok {
				dataFile.Shortcuts[kId][scId] = f.shortcuts[kId][scId]
			}
		}
	}
	return *dataFile, nil
}

func (f *FileStorage) saveToFile() error {
	if len(f.keylogs) == 0 && len(f.shortcuts) == 0 {
		return nil
	}
	start := time.Now()
	data, err := f.prepareDataToSave()
	if err != nil {
		return err
	}
	pb, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(f.fname, pb, 0777)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("| %s | File updated.\n", time.Since(start)))
	f.keylogs = map[int64]map[uint16]int64{}
	f.shortcuts = map[int64]map[int64]int64{}
	return nil
}

func (f *FileStorage) savingInBackground(ctx context.Context) {
	for {
		select {
		case <-time.After(3 * time.Second):
			// TODO: And set time to save every 30 s
			f.saveToFile()
		case <-ctx.Done():
			// TODO: gracefull shutdown, make last save .
			slog.Info("Leaving goroutine...")
			return
		}
	}
}
