package keylogger

/*
// -I include current directory for headers
#cgo CFLAGS: -I.
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit

#include "keylogger_darwin.h"
*/
import "C"

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/keylogme/keylogme-zero/types"
)

// var hidManager = map[int]map[int]chan InputEvent{}

type hidManager struct {
	mapVendorProductChan map[int]map[int]chan InputEvent
	mu                   sync.Mutex
}

var hid = hidManager{
	mapVendorProductChan: make(map[int]map[int]chan InputEvent),
	mu:                   sync.Mutex{},
}

func (h *hidManager) exists(vendorID, productID int) (chan InputEvent, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.mapVendorProductChan[vendorID]; !ok {
		return nil, false
	}
	if _, ok := h.mapVendorProductChan[vendorID][productID]; !ok {
		return nil, false
	}
	return h.mapVendorProductChan[vendorID][productID], true
}

func (h *hidManager) setChannel(vendorID, productID int) chan InputEvent {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.mapVendorProductChan[vendorID]; !ok {
		h.mapVendorProductChan[vendorID] = map[int]chan InputEvent{}
	}
	if _, ok := h.mapVendorProductChan[vendorID][productID]; !ok {
		event := make(chan InputEvent)
		h.mapVendorProductChan[vendorID][productID] = event
	}
	return h.mapVendorProductChan[vendorID][productID]
}

func (h *hidManager) closeChannel(vendorID, productID int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.mapVendorProductChan[vendorID]; !ok {
		return
	}
	if _, ok := h.mapVendorProductChan[vendorID][productID]; !ok {
		return
	}
	close(h.mapVendorProductChan[vendorID][productID]) // close channel
	delete(h.mapVendorProductChan[vendorID], productID)

	// remove vendorID if no products left
	if len(h.mapVendorProductChan[vendorID]) == 0 {
		delete(h.mapVendorProductChan, vendorID)
	}
	// stop if no devices left
	if len(h.mapVendorProductChan) == 0 {
		C.Stop()
	}
}

type KeyLogger struct {
	vendorID  int
	productID int
}

func NewKeylogger(kInput types.KeyloggerInput) (*KeyLogger, error) {
	// C.ListConnectedHIDDevices()
	exists := C.checkDeviceIsConnected(C.int(kInput.VendorId), C.int(kInput.ProductId))
	if !exists {
		slog.Debug("Device not found")
		return nil, fmt.Errorf("device not available")
	}
	k := &KeyLogger{vendorID: int(kInput.VendorId), productID: int(kInput.ProductId)}

	go func() {
		// TODO: add mutex to prevent multiple calls
		C.Start()
	}()
	return k, nil
}

func (k *KeyLogger) Read() chan InputEvent {
	return hid.setChannel(k.vendorID, k.productID)
}

func (k *KeyLogger) Close() error {
	hid.closeChannel(k.vendorID, k.productID)
	return nil
}

//export GoHandleKeyEvent
func GoHandleKeyEvent(code, value, vendorID, productID C.int) {
	vID := int(vendorID)
	pID := int(productID)

	c, ok := hid.exists(vID, pID)
	if !ok {
		slog.Debug("Vendor id and product id not in HIDManager")
		return
	}
	pressed := int32(value)
	if pressed != 0 && pressed != 1 {
		return
	}
	c <- InputEvent{Time: time.Now(), Code: uint16(code), Type: KeyEvent(value)}
}

//export GoHandleDeviceEvent
func GoHandleDeviceEvent(vendorID, productID, connected C.int) {
	vID := int(vendorID)
	pID := int(productID)
	status := "disconnected"
	if connected != 0 {
		status = "connected"
		// INFO: Read will add device to hidManager
		// fmt.Printf("[Device] %s: VID=0x%04x, PID=0x%04x\n", status, vendorID, productID)
		return
	}
	// disconnect
	fmt.Printf("[Device] %s: VID=0x%04x, PID=0x%04x\n", status, vendorID, productID)
	hid.closeChannel(vID, pID)
}
