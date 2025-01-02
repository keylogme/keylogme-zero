package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/keylogme/keylogme-zero/keylog/utils"
)

type ConfigStorage struct {
	FileOutput        string        `json:"file_output"`
	PeriodicSaveInSec time.Duration `json:"periodic_save_in_sec"`
}

func (c *ConfigStorage) Validate() error {
	if c.FileOutput == "" {
		return errors.New("file_output is required")
	}
	absPath, err := filepath.Abs(c.FileOutput)
	if err != nil {
		return err
	}
	fmt.Printf("File will be saved at %s\n", absPath)
	c.FileOutput = absPath
	if c.PeriodicSaveInSec == 0 {
		return errors.New("periodic_save_in_sec is required")
	}
	return nil
}

type Storage interface {
	SaveKeylog(deviceId string, keycode uint16) error
	SaveShortcut(deviceId string, shortcutId string) error
}

type FileStorage struct {
	config    ConfigStorage
	keylogs   map[string]map[uint16]int64 // deviceId - keycode - counter
	shortcuts map[string]map[string]int64 // deviceId - shortcutId - counter
}

type DataFile struct {
	Keylogs   map[string]map[uint16]int64 `json:"keylogs,omitempty"`
	Shortcuts map[string]map[string]int64 `json:"shortcuts,omitempty"`
}

func newDataFile() DataFile {
	return DataFile{
		Keylogs:   map[string]map[uint16]int64{},
		Shortcuts: map[string]map[string]int64{},
	}
}

func MustGetNewFileStorage(ctx context.Context, config ConfigStorage) *FileStorage {
	err := config.Validate()
	if err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
	ffs := &FileStorage{
		config:    config,
		keylogs:   map[string]map[uint16]int64{},
		shortcuts: map[string]map[string]int64{},
	}
	go ffs.savingInBackground(ctx)
	return ffs
}

func (f *FileStorage) SaveKeylog(deviceId string, keycode uint16) error {
	if _, ok := f.keylogs[deviceId]; !ok {
		f.keylogs[deviceId] = map[uint16]int64{}
	}
	if _, ok := f.keylogs[deviceId][keycode]; !ok {
		f.keylogs[deviceId][keycode] = 0
	}
	f.keylogs[deviceId][keycode] += 1
	return nil
}

func (f *FileStorage) SaveShortcut(deviceId string, shortcutId string) error {
	if _, ok := f.shortcuts[deviceId]; !ok {
		f.shortcuts[deviceId] = map[string]int64{}
	}
	if _, ok := f.shortcuts[deviceId][shortcutId]; !ok {
		f.shortcuts[deviceId][shortcutId] = 0
	}
	f.shortcuts[deviceId][shortcutId] += 1
	return nil
}

func (f *FileStorage) prepareDataToSave() (DataFile, error) {
	dataFile := newDataFile()
	_, err := os.Stat(f.config.FileOutput)
	if errors.Is(err, os.ErrNotExist) {
		slog.Info(fmt.Sprintf("File %s not exist, it will be created", f.config.FileOutput))
	} else {
		err := utils.ParseFromFile(f.config.FileOutput, &dataFile)
		if err != nil {
			return dataFile, err
		}
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
				dataFile.Shortcuts[kId] = map[string]int64{}
			}
			if _, ok := dataFile.Shortcuts[kId][scId]; !ok {
				dataFile.Shortcuts[kId][scId] = f.shortcuts[kId][scId]
			}
		}
	}
	return dataFile, nil
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
	err = os.WriteFile(f.config.FileOutput, pb, 0777)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("| %s | File %s updated.\n", time.Since(start), f.config.FileOutput))
	// Reset data
	f.keylogs = map[string]map[uint16]int64{}
	f.shortcuts = map[string]map[string]int64{}
	return nil
}

func (f *FileStorage) savingInBackground(ctx context.Context) {
	for {
		select {
		case <-time.After(f.config.PeriodicSaveInSec * time.Second):
			f.saveToFile()
		case <-ctx.Done():
			slog.Info("Closing file storage...")
			f.saveToFile()
			return
		}
	}
}
