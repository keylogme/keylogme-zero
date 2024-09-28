package main

import (
	"encoding/json"
	"fmt"
	"gokeny/internal"
	"gokeny/internal/keylogger"
	"log/slog"
	"os"
	"slices"
	"time"
)

type KeyLog struct {
	Code uint16 `json:"code"`
}

type Config struct {
	Shortcuts []internal.Shortcut
}

// Use lsinput to see which input to be used
// apt install input-utils
// sudo lsinput
// If your keyboard name appeared multiple times,
// try with all of them

func main() {
	APIKEY := os.Args[1]
	ORIGIN_ENDPOINT := os.Args[2]
	// Get config
	config := Config{
		Shortcuts: []internal.Shortcut{
			{ID: 1, Values: []string{"J", "S"}, Type: internal.SequentialShortcutType},
			{ID: 2, Values: []string{"J", "F"}, Type: internal.SequentialShortcutType},
			{ID: 3, Values: []string{"J", "G"}, Type: internal.SequentialShortcutType},
			{ID: 4, Values: []string{"J", "S", "G"}, Type: internal.SequentialShortcutType},
		},
	}
	chEvt := make(chan keylogger.InputEvent)
	sender := internal.MustGetNewSender(ORIGIN_ENDPOINT, APIKEY)
	defer sender.Close()

	sd := internal.NewShortcutsDetector(config.Shortcuts)

	keylogger.MustGetDevice("foostan Corne", chEvt)
	keylogger.MustGetDevice("MOSART Semi. 2.4G INPUT DEVICE Mouse", chEvt)

	modifiers := []uint16{29, 97, 42, 54, 56, 100} // ctrl, shft, alt

	slog.Info("Listening...")

	modPress := []uint16{}
	for i := range chEvt {
		if i.KeyPress() && slices.Contains(modifiers, i.Code) {
			modPress = append(modPress, i.Code)
		}
		if i.Type == keylogger.EvKey && i.KeyRelease() {
			start := time.Now()

			detectedShortcutID := sd.Detect(i.KeyString())
			if detectedShortcutID != 0 {
				sendShortcut(sender, detectedShortcutID)
			}
			//
			// FIXME: mod+key is sent, but when mod is released , is sent again
			keylogs := []uint16{i.Code}
			keylogs = append(keylogs, modPress...)
			sendKeylog(sender, keylogs)
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
	Type TypePayload     `json:"type"`
	Data json.RawMessage `json:"data"`
}

func getPayload(typePayload TypePayload, data any) ([]byte, error) {
	db, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	p := Payload{Type: typePayload, Data: db}
	pb, err := json.Marshal(p)
	if err != nil {
		return []byte{}, err
	}
	return pb, nil
}

func sendKeylog(ws *internal.Sender, kls []uint16) error {
	payloadBytes, err := getPayload(KeyLogPayload, kls)
	if err != nil {
		return err
	}
	return ws.Send(payloadBytes)
}

func sendShortcut(ws *internal.Sender, scID int64) error {
	start := time.Now()
	defer func() {
		slog.Info(fmt.Sprintf("| %s | Shortcut %d\n", time.Since(start), scID))
	}()
	pb, err := getPayload(ShortcutPayload, scID)
	if err != nil {
		return err
	}
	return ws.Send(pb)
}

// TODO: add ws conn to send to backend
