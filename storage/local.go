package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/keylogme/keylogme-zero/types"
	"github.com/keylogme/keylogme-zero/utils"
)

type ConfigStorage struct {
	FileOutput   string         `json:"file_output"`
	PeriodicSave types.Duration `json:"periodic_save"`
}

func (c *ConfigStorage) Validate() error {
	if c.FileOutput == "" {
		return errors.New("file_output is required")
	}
	// absPath, err := utils.ExpandUserHome(c.FileOutput)
	// if err != nil {
	// 	return err
	// }
	slog.Info(fmt.Sprintf("File will be saved at %s\n", c.FileOutput))
	// c.FileOutput = absPath
	if c.PeriodicSave.Duration == 0 {
		return errors.New("periodic_save_in_sec is required and non zero")
	}
	return nil
}

type FileStorage struct {
	config   ConfigStorage
	dataFile *DataFile
}

type DataFile struct {
	mu sync.Mutex
	// deviceId - layerId - keycode - counter
	Keylogs map[string]map[int64]map[uint16]int64 `json:"keylogs,omitempty"`
	// deviceId - shortcutId - counter
	Shortcuts map[string]map[string]int64 `json:"shortcuts,omitempty"`
	// deviceId - modifier - keycode - counter
	ShiftStates map[string]map[uint16]map[uint16]int64 `json:"shift_states,omitempty"`
	// deviceId - modifier - keycode - counter
	ShiftStatesAuto map[string]map[uint16]map[uint16]int64 `json:"shift_states_auto,omitempty"`

	// deviceId - layerId - counter
	LayerChanges map[string]map[int64]int64 `json:"layer_changes,omitempty"`
}

func (d *DataFile) AddKeylog(deviceId string, layerId int64, keycode uint16, addQty int64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Keylogs[deviceId]; !ok {
		d.Keylogs[deviceId] = map[int64]map[uint16]int64{}
	}
	if _, ok := d.Keylogs[deviceId][layerId]; !ok {
		d.Keylogs[deviceId][layerId] = map[uint16]int64{}
	}
	if _, ok := d.Keylogs[deviceId][layerId][keycode]; !ok {
		d.Keylogs[deviceId][layerId][keycode] = 0
	}
	d.Keylogs[deviceId][layerId][keycode] += addQty
}

func (d *DataFile) AddShortcut(deviceId string, shortcutId string, addQty int64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Shortcuts[deviceId]; !ok {
		d.Shortcuts[deviceId] = map[string]int64{}
	}
	if _, ok := d.Shortcuts[deviceId][shortcutId]; !ok {
		d.Shortcuts[deviceId][shortcutId] = 0
	}
	d.Shortcuts[deviceId][shortcutId] += addQty
}

func updateShiftState(
	shiftStates *map[string]map[uint16]map[uint16]int64,
	deviceId string,
	modifier uint16,
	keycode uint16,
	addQty int64,
) {
	ss := (*shiftStates)
	if _, ok := ss[deviceId]; !ok {
		ss[deviceId] = map[uint16]map[uint16]int64{}
	}
	if _, ok := ss[deviceId][modifier]; !ok {
		ss[deviceId][modifier] = map[uint16]int64{}
	}
	if _, ok := ss[deviceId][modifier][keycode]; !ok {
		ss[deviceId][modifier][keycode] = 0
	}
	ss[deviceId][modifier][keycode] += addQty
}

func (d *DataFile) AddShiftState(
	deviceId string,
	modifier uint16,
	keycode uint16,
	addQty int64,
) {
	d.mu.Lock()
	defer d.mu.Unlock()
	updateShiftState(&d.ShiftStates, deviceId, modifier, keycode, addQty)
}

func (d *DataFile) AddLayerChange(deviceId string, layerId int64, addQty int64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.LayerChanges[deviceId]; !ok {
		d.LayerChanges[deviceId] = map[int64]int64{}
	}
	if _, ok := d.LayerChanges[deviceId][layerId]; !ok {
		d.LayerChanges[deviceId][layerId] = 0
	}
	d.LayerChanges[deviceId][layerId] += addQty
}

func (d *DataFile) AddAutoShiftState(
	deviceId string,
	modifier uint16,
	keycode uint16,
	addQty int64,
) {
	d.mu.Lock()
	defer d.mu.Unlock()
	updateShiftState(&d.ShiftStatesAuto, deviceId, modifier, keycode, addQty)
}

func (d *DataFile) Merge(data *DataFile) {
	// Lock input datafile
	data.mu.Lock()
	defer data.mu.Unlock()
	//
	for kId := range data.Keylogs {
		for layerId := range data.Keylogs[kId] {
			for keycode := range data.Keylogs[kId][layerId] {
				d.AddKeylog(kId, layerId, keycode, data.Keylogs[kId][layerId][keycode])
			}
		}
	}
	for kId := range data.Shortcuts {
		for scId := range data.Shortcuts[kId] {
			d.AddShortcut(kId, scId, data.Shortcuts[kId][scId])
		}
	}
	for kId := range data.ShiftStates {
		for modifier := range data.ShiftStates[kId] {
			for keycode := range data.ShiftStates[kId][modifier] {
				d.AddShiftState(kId, modifier, keycode, data.ShiftStates[kId][modifier][keycode])
			}
		}
	}
	for kId := range data.ShiftStatesAuto {
		for modifier := range data.ShiftStatesAuto[kId] {
			for keycode := range data.ShiftStatesAuto[kId][modifier] {
				d.AddAutoShiftState(
					kId,
					modifier,
					keycode,
					data.ShiftStatesAuto[kId][modifier][keycode],
				)
			}
		}
	}
	for kId := range data.LayerChanges {
		for layerId := range data.LayerChanges[kId] {
			d.AddLayerChange(kId, layerId, data.LayerChanges[kId][layerId])
		}
	}
}

func (d *DataFile) Reset() {
	d.Keylogs = map[string]map[int64]map[uint16]int64{}
	d.Shortcuts = map[string]map[string]int64{}
	d.ShiftStates = map[string]map[uint16]map[uint16]int64{}
	d.ShiftStatesAuto = map[string]map[uint16]map[uint16]int64{}
	d.LayerChanges = map[string]map[int64]int64{}
}

func newDataFile() *DataFile {
	d := DataFile{}
	d.Reset()
	return &d
}

func MustGetNewFileStorage(ctx context.Context, config ConfigStorage) *FileStorage {
	err := config.Validate()
	if err != nil {
		log.Fatalf("Invalid config: %v", err.Error())
	}
	ffs := &FileStorage{
		config:   config,
		dataFile: newDataFile(),
	}
	go ffs.savingInBackground(ctx)
	return ffs
}

func (f *FileStorage) SaveKeylog(deviceId string, layerId int64, keycode uint16) error {
	f.dataFile.AddKeylog(deviceId, layerId, keycode, 1)
	return nil
}

func (f *FileStorage) SaveShortcut(deviceId string, shortcutId string) error {
	f.dataFile.AddShortcut(deviceId, shortcutId, 1)
	return nil
}

func (f *FileStorage) SaveShiftState(
	deviceId string,
	modifier uint16,
	keycode uint16,
	auto bool,
) error {
	if auto {
		updateShiftState(&f.dataFile.ShiftStatesAuto, deviceId, modifier, keycode, 1)
	} else {
		updateShiftState(&f.dataFile.ShiftStates, deviceId, modifier, keycode, 1)
	}
	return nil
}

func (f *FileStorage) SaveLayerChange(deviceId string, layerId int64) error {
	f.dataFile.AddLayerChange(deviceId, layerId, 1)
	return nil
}

func (f *FileStorage) prepareDataToSave() error {
	dataFile := newDataFile()
	_, err := os.Stat(f.config.FileOutput)
	if errors.Is(err, os.ErrNotExist) {
		slog.Info(fmt.Sprintf("File %s not exist, it will be created", f.config.FileOutput))
	} else {
		err := utils.ParseFromFile(f.config.FileOutput, dataFile)
		if err != nil {
			return err
		}
	}
	f.dataFile.Merge(dataFile)
	return nil
}

func (f *FileStorage) saveToFile() error {
	// if len(f.keylogs) == 0 && len(f.shortcuts) == 0 {
	// 	return nil
	// }
	start := time.Now()
	err := f.prepareDataToSave()
	if err != nil {
		return err
	}
	pb, err := json.Marshal(f.dataFile)
	if err != nil {
		return err
	}
	// Permissions for the new file (e.g., 0644 for read/write for owner, read-only for others)
	// You can use os.FileMode(0644) or a more readable form like 0o644 (octal literal)
	permissions := os.FileMode(0o644) // Or 0o644

	err = os.WriteFile(f.config.FileOutput, pb, permissions)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("| %s | File %s updated.\n", time.Since(start), f.config.FileOutput))
	// Reset data
	f.dataFile.Reset()
	return nil
}

func (f *FileStorage) savingInBackground(ctx context.Context) {
	for {
		select {
		case <-time.After(f.config.PeriodicSave.Duration):
			err := f.saveToFile()
			if err != nil {
				slog.Error(fmt.Sprintf("Error saving file: %v", err))
			}
		case <-ctx.Done():
			slog.Info("Closing file storage...")
			err := f.saveToFile()
			if err != nil {
				slog.Error(fmt.Sprintf("Error saving file: %v", err))
			}
			return
		}
	}
}

func (f *FileStorage) CloneInMemoryDatafile() (*DataFile, error) {
	copyDatafile := &DataFile{}
	data, err := json.Marshal(f.dataFile)
	if err != nil {
		return copyDatafile, err
	}
	err = json.Unmarshal(data, copyDatafile)
	if err != nil {
		return copyDatafile, err
	}
	return copyDatafile, nil
}
