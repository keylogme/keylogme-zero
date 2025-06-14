package k0

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/storage"
	"github.com/keylogme/keylogme-zero/types"
)

func TestStart(t *testing.T) {
	chEvt := make(chan DeviceEvent)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tf, err := os.MkdirTemp("", "local_test_start")
	if err != nil {
		t.Fatal(err)
	}
	filename := fmt.Sprintf("device_%d", rand.Int())
	filepath := path.Join(tf, filename)
	periodicSave := 20 * time.Second
	config := storage.ConfigStorage{
		FileOutput:   filepath,
		PeriodicSave: types.Duration{Duration: periodicSave},
	}
	ffs := storage.MustGetNewFileStorage(ctx, config)

	inputDevice := DeviceInput{
		DeviceId: "device1",
		Name:     "device1",
		Layers:   []LayerInput{},
	}

	si := SecurityInput{
		BaggageSize:   10,
		GhostingCodes: []uint16{25, 26, 27},
	}
	s := NewSecurity(si)

	seqShortcut := ShortcutCodes{
		Id:    "1",
		Name:  "dummy1",
		Codes: []uint16{4, 5},
		Type:  SequentialShortcutType,
	}
	shiftKey := keylogger.GetShiftCodes()[0]
	holdShortcutCode := uint16(11)
	holdShortcut := ShortcutCodes{
		Id:    "2",
		Name:  "dummy2",
		Codes: []uint16{shiftKey, holdShortcutCode},
		Type:  HoldShortcutType,
	}

	sgs := []ShortcutGroupInput{
		{
			Id: "1", Name: "group1", Shortcuts: []ShortcutCodes{
				seqShortcut, holdShortcut,
			},
		},
	}
	sd := MustGetNewShortcutsDetector(sgs)

	configShiftState := getTestShiftStateConfig()
	ss := NewShiftStateDetector(configShiftState)

	ld := NewLayersDetector([]DeviceInput{inputDevice}, configShiftState)

	Start(chEvt, s, sd, ss, ld, ffs)

	deviceId := "1"
	defaultLayer := int64(0)

	dePress := DeviceEvent{
		InputEvent: keylogger.InputEvent{Time: time.Now(), Code: 16, Type: keylogger.KeyPress},
		DeviceId:   deviceId,
	}
	deRelease := DeviceEvent{
		InputEvent: keylogger.InputEvent{Time: time.Now(), Code: 16, Type: keylogger.KeyRelease},
		DeviceId:   deviceId,
	}

	// security check

	// baggage
	for i := range si.BaggageSize {
		fmt.Println(i)
		chEvt <- dePress
		chEvt <- deRelease
	}
	data, err := ffs.CloneInMemoryDatafile()
	if err != nil {
		t.Fatal("error getting data")
	}

	if _, ok := data.Keylogs[deviceId]; ok {
		t.Fatal("log should be blocked by security layer (baggage not full)")
	}
	// first key to be saved
	chEvt <- dePress
	chEvt <- deRelease

	time.Sleep(20 * time.Millisecond)

	data, err = ffs.CloneInMemoryDatafile()
	if err != nil {
		t.Fatal("error getting data")
	}

	if _, ok := data.Keylogs[deviceId]; !ok {
		t.Fatal("first log should exist")
	}
	count, ok := data.Keylogs[deviceId][defaultLayer][16]
	if !ok {
		t.Fatal("code is not in data saved")
	}
	if count != 1 {
		t.Fatal("unexpected code count")
	}

	// ghost codes
	ghostCode := si.GhostingCodes[0]
	dePress.Code = ghostCode
	deRelease.Code = ghostCode

	chEvt <- dePress
	chEvt <- deRelease

	time.Sleep(20 * time.Millisecond)

	data, err = ffs.CloneInMemoryDatafile()
	if err != nil {
		t.Fatal("error getting data")
	}

	if _, ok := data.Keylogs[deviceId][defaultLayer][ghostCode]; ok {
		t.Fatal("ghost code should not have been saved")
	}

	// shortcut

	// sequential shortcut
	for _, c := range seqShortcut.Codes {
		dePress.Code = c
		deRelease.Code = c

		chEvt <- dePress
		chEvt <- deRelease
	}
	time.Sleep(20 * time.Millisecond) // wait to process channel

	data, err = ffs.CloneInMemoryDatafile()
	if err != nil {
		t.Fatal("error getting data")
	}

	if _, ok := data.Shortcuts[deviceId]; !ok {
		t.Fatal("device should have seq shortcuts saved")
	}
	if _, ok := data.Shortcuts[deviceId][seqShortcut.Id]; !ok {
		t.Fatal("seq shortcut code id should have been saved")
	}
	if data.Shortcuts[deviceId][seqShortcut.Id] != 1 {
		t.Fatal("unexpected seq shortcut count")
	}

	// hold shortcut
	chEvt <- getFakeEvent(deviceId, shiftKey, keylogger.KeyPress)
	chEvt <- getFakeEvent(deviceId, holdShortcutCode, keylogger.KeyPress)

	chEvt <- getFakeEvent(deviceId, holdShortcutCode, keylogger.KeyRelease)
	chEvt <- getFakeEvent(deviceId, shiftKey, keylogger.KeyRelease)

	time.Sleep(20 * time.Millisecond) // wait to process channel

	data, err = ffs.CloneInMemoryDatafile()
	if err != nil {
		t.Fatal("error getting data")
	}

	if _, ok := data.Shortcuts[deviceId][holdShortcut.Id]; !ok {
		t.Fatal("hold shortcut code id should have been saved")
	}
	if data.Shortcuts[deviceId][holdShortcut.Id] != 1 {
		t.Fatal("unexpected hold shortcut count")
	}

	// shift state

	// auto=true
	// the hold shortcut was also a shifted key (auto=true because it was pressed  fast)
	if _, ok := data.ShiftStatesAuto[deviceId]; !ok {
		t.Fatal("device should have shift state auto")
	}
	if data.ShiftStatesAuto[deviceId][shiftKey][holdShortcutCode] != 1 {
		t.Fatal("unexpected auto shifted state  count")
	}
	if _, ok := data.ShiftStates[deviceId]; ok {
		t.Fatal("device should not have shift state")
	}

	// auto=false
	delay := configShiftState.ThresholdAuto.Duration + 5*time.Millisecond
	chEvt <- getFakeEvent(deviceId, shiftKey, keylogger.KeyPress)
	time.Sleep(delay) // wait so  detector thinks it is human triggered
	chEvt <- getFakeEvent(deviceId, holdShortcutCode, keylogger.KeyPress)
	chEvt <- getFakeEvent(deviceId, holdShortcutCode, keylogger.KeyRelease)
	time.Sleep(delay) // wait so  detector thinks it is human triggered
	chEvt <- getFakeEvent(deviceId, shiftKey, keylogger.KeyRelease)

	time.Sleep(20 * time.Millisecond) // wait to process channel

	data, err = ffs.CloneInMemoryDatafile()
	if err != nil {
		t.Fatal("error getting data")
	}

	if _, ok := data.ShiftStates[deviceId]; !ok {
		t.Fatal("device should have shift state")
	}
	if data.ShiftStates[deviceId][shiftKey][holdShortcutCode] != 1 {
		t.Fatal("unexpected shifted state count")
	}
}
