package main

import (
	"encoding/json"
	"fmt"
	"gokeny/internal"
	"log"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/MarinX/keylogger"
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
		},
	}
	sender := internal.MustGetNewSender(ORIGIN_ENDPOINT, APIKEY)
	defer sender.Close()

	sd := internal.NewShortcutsDetector(config.Shortcuts)
	// keylogger
	// selectedKb := "/dev/input/event14" // name: foostan Corne , phys: usb-0000:00:14.0-4.3/input0
	// selectedKb := "/dev/input/event9" // mouse vertical
	selectedKb := "/dev/input/event18"
	// selectedKb := "/dev/input/event11"
	slog.Info(fmt.Sprintf("Device selected: %s\n", selectedKb))
	kl, err := keylogger.New(selectedKb)
	if err != nil {
		log.Fatal(err)
	}
	defer kl.Close()

	chIn := kl.Read()
	modifiers := []uint16{29, 97, 42, 54, 56, 100} // ctrl, shft, alt
	slog.Info("Listening...")
	modPress := []uint16{}
	for i := range chIn {
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
	Type TypePayload
	Data []byte
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
