package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/keylogme/zero-trust-logger/keylog"
	"github.com/keylogme/zero-trust-logger/keylog/storage"
)

type KeyLog struct {
	Code uint16 `json:"code"`
}

// Use lsinput to see which input to be used
// apt install input-utils
// sudo lsinput
// If your keyboard name appeared multiple times,
// try with all of them

func main() {
	// Get config
	config := keylog.Config{
		Devices: []keylog.DeviceInput{
			{Id: 1, Name: "foostan Corne"},
			{Id: 2, Name: "MOSART Semi. 2.4G INPUT DEVICE Mouse"},
			{Id: 2, Name: "Logitech MX Master 2S"},
			// {Id: 2, Name: "Wacom Intuos BT M Pen"},
		},
		Shortcuts: []keylog.Shortcut{
			{Id: 1, Values: []string{"J", "S"}, Type: keylog.SequentialShortcutType},
			{Id: 2, Values: []string{"J", "F"}, Type: keylog.SequentialShortcutType},
			{Id: 3, Values: []string{"J", "G"}, Type: keylog.SequentialShortcutType},
			{Id: 4, Values: []string{"J", "S", "G"}, Type: keylog.SequentialShortcutType},
		},
	}
	chEvt := make(chan keylog.DeviceEvent)
	ffs := storage.NewFileStorage(context.Background(), "test.json")

	sd := keylog.NewShortcutsDetector(config.Shortcuts)

	keylog.GetDevice(config.Devices[0], chEvt)
	keylog.GetDevice(config.Devices[1], chEvt)
	keylog.GetDevice(config.Devices[2], chEvt)

	modifiers := []uint16{29, 97, 42, 54, 56, 100} // ctrl, shft, alt

	slog.Info("Listening...")

	modPress := []uint16{}
	for i := range chEvt {
		if i.KeyPress() && slices.Contains(modifiers, i.Code) {
			modPress = append(modPress, i.Code)
		}
		if i.Type == keylog.EvKey && i.KeyRelease() {
			start := time.Now()

			detectedShortcutID := sd.Detect(i.KeyString())
			if detectedShortcutID != 0 {
				// sendShortcut(sender, i.DeviceId, detectedShortcutID)
				ffs.SaveShortcut(i.DeviceId, detectedShortcutID)
			}
			//
			// FIXME: mod+key is sent, but when mod is released , is sent again
			// keylogs := []uint16{i.Code}
			// keylogs = append(keylogs, modPress...)
			// err := sendKeylog(sender, i.DeviceId, i.Code)
			err := ffs.SaveKeylog(i.DeviceId, i.Code)
			if err != nil {
				fmt.Printf("error %s\n", err.Error())
			}
			slog.Info(fmt.Sprintf("| %s | Key :%d %s\n", time.Since(start), i.Code, i.KeyString()))
			// Reset modPress
			modPress = []uint16{}
		}
	}
	fmt.Println("Closing...")
}

// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	log.Printf("%s took %s", name, elapsed)
// }

type TypePayload string

const (
	KeyLogPayload   TypePayload = "kl"
	ShortcutPayload TypePayload = "sc"
)

type Payload struct {
	Version int             `json:"version"`
	Type    TypePayload     `json:"type"`
	Data    json.RawMessage `json:"data"` // why not json.RawMessage?
}

type KeylogPayloadV1 struct {
	KeyboardDeviceId int64  `json:"kID"`
	Code             uint16 `json:"c"`
}

type ShortcutPayloadV1 struct {
	KeyboardDeviceId int64 `json:"kID"`
	ShortcutId       int64 `json:"scID"`
}

func getPayload(typePayload TypePayload, data any) ([]byte, error) {
	db, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	p := Payload{Version: 1, Type: typePayload, Data: db}
	pb, err := json.Marshal(p)
	if err != nil {
		return []byte{}, err
	}
	return pb, nil
}

func sendKeylog(ws *keylog.Sender, kId int64, code uint16) error {
	payloadBytes, err := getPayload(
		KeyLogPayload,
		KeylogPayloadV1{KeyboardDeviceId: kId, Code: code},
	)
	if err != nil {
		return err
	}
	return ws.Send(payloadBytes)
}

func sendShortcut(ws *keylog.Sender, kId, scID int64) error {
	start := time.Now()
	defer func() {
		slog.Info(fmt.Sprintf("| %s | Shortcut %d\n", time.Since(start), scID))
	}()
	pb, err := getPayload(
		ShortcutPayload,
		ShortcutPayloadV1{KeyboardDeviceId: kId, ShortcutId: scID},
	)
	if err != nil {
		return err
	}
	return ws.Send(pb)
}

// TODO: add ws conn to send to backend
