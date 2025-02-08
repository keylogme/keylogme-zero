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
	SaveShiftState(deviceId string, modifier uint16, keycode uint16, auto bool) error
}

type FileStorage struct {
	config   ConfigStorage
	dataFile DataFileV1
}

type DataFileV1 struct {
	// deviceId - keycode - counter
	Keylogs map[string]map[uint16]int64 `json:"keylogs,omitempty"`
	// deviceId - shortcutId - counter
	Shortcuts map[string]map[string]int64 `json:"shortcuts,omitempty"`
	// deviceId - modifier - keycode - counter
	ShiftStates map[string]map[uint16]map[uint16]int64 `json:"shift_states,omitempty"`
	// deviceId - modifier - keycode - counter
	ShiftStatesAuto map[string]map[uint16]map[uint16]int64 `json:"shift_states_auto,omitempty"`
}

func (d *DataFileV1) AddKeylog(deviceId string, keycode uint16, addQty int64) {
	if _, ok := d.Keylogs[deviceId]; !ok {
		d.Keylogs[deviceId] = map[uint16]int64{}
	}
	if _, ok := d.Keylogs[deviceId][keycode]; !ok {
		d.Keylogs[deviceId][keycode] = 0
	}
	d.Keylogs[deviceId][keycode] += addQty
}

func (d *DataFileV1) AddShortcut(deviceId string, shortcutId string, addQty int64) {
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

func (d *DataFileV1) AddShiftState(
	deviceId string,
	modifier uint16,
	keycode uint16,
	addQty int64,
) {
	updateShiftState(&d.ShiftStates, deviceId, modifier, keycode, addQty)
}

func (d *DataFileV1) AddAutoShiftState(
	deviceId string,
	modifier uint16,
	keycode uint16,
	addQty int64,
) {
	updateShiftState(&d.ShiftStatesAuto, deviceId, modifier, keycode, addQty)
}

func (d *DataFileV1) Merge(data DataFileV1) {
	for kId := range data.Keylogs {
		for keycode := range data.Keylogs[kId] {
			d.AddKeylog(kId, keycode, data.Keylogs[kId][keycode])
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
}

func (d *DataFileV1) Reset() {
	d.Keylogs = map[string]map[uint16]int64{}
	d.Shortcuts = map[string]map[string]int64{}
	d.ShiftStates = map[string]map[uint16]map[uint16]int64{}
	d.ShiftStatesAuto = map[string]map[uint16]map[uint16]int64{}
}

func newDataFile() DataFileV1 {
	d := DataFileV1{}
	d.Reset()
	return d
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

func (f *FileStorage) SaveKeylog(deviceId string, keycode uint16) error {
	f.dataFile.AddKeylog(deviceId, keycode, 1)
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

func (f *FileStorage) prepareDataToSave() error {
	dataFile := newDataFile()
	_, err := os.Stat(f.config.FileOutput)
	if errors.Is(err, os.ErrNotExist) {
		slog.Info(fmt.Sprintf("File %s not exist, it will be created", f.config.FileOutput))
	} else {
		err := utils.ParseFromFile(f.config.FileOutput, &dataFile)
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
	err = os.WriteFile(f.config.FileOutput, pb, 0777)
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
		case <-time.After(f.config.PeriodicSaveInSec * time.Second):
			f.saveToFile()
		case <-ctx.Done():
			slog.Info("Closing file storage...")
			f.saveToFile()
			return
		}
	}
}
