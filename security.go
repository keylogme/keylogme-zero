package k0

import (
	"fmt"
	"log/slog"
	"math/rand"
	"slices"
)

// baggage holds a map of device IDs to their codes,
// allowing for a minimum number of keypresses to start keylogging
type baggage struct {
	size    int
	devices map[string][]uint16
}

func newBaggage(size int, d map[string][]uint16) *baggage {
	return &baggage{
		size:    size,
		devices: d,
	}
}

func getRandInt(max int) int {
	return rand.Intn(max)
}

func (b *baggage) isAuthorized(ke *DeviceEvent) bool {
	if _, exists := b.devices[ke.DeviceId]; !exists {
		b.devices[ke.DeviceId] = make([]uint16, 0, b.size)
	}

	sizeBaggageDevice := len(b.devices[ke.DeviceId])
	if sizeBaggageDevice < b.size {
		b.devices[ke.DeviceId] = append(b.devices[ke.DeviceId], ke.Code)
		// INFO: baggage has to fill up first before authorization
		return false
	}
	randomIndex := getRandInt(b.size)

	// INFO: swap codes
	currentCode := ke.Code
	ke.Code = b.devices[ke.DeviceId][randomIndex]
	b.devices[ke.DeviceId][randomIndex] = currentCode
	return true
}

// ghostingCodes holds a list of ghost codes that are not authorized to keylog
type ghostingCodes struct {
	ghostCodes []uint16
}

func newGhostingCodes(c []uint16) *ghostingCodes {
	return &ghostingCodes{
		ghostCodes: c,
	}
}

func (gh *ghostingCodes) isAuthorized(ke *DeviceEvent) bool {
	if slices.Contains(gh.ghostCodes, ke.Code) {
		// INFO: not authorized if code is in ghost codes
		return false
	}
	return true
}

type security struct {
	baggage       *baggage
	ghostingCodes *ghostingCodes
}

type SecurityInput struct {
	BaggageSize   int                 `json:"baggage_size"`
	Baggage       map[string][]uint16 `json:"baggage_per_device"`
	GhostingCodes []uint16            `json:"ghosting_codes"`
}

func NewSecurity(secInput SecurityInput) *security {
	return &security{
		baggage:       newBaggage(secInput.BaggageSize, secInput.Baggage),
		ghostingCodes: newGhostingCodes(secInput.GhostingCodes),
	}
}

func (s *security) isAuthorized(ke *DeviceEvent) bool {
	auth := s.ghostingCodes.isAuthorized(ke)
	if !auth {
		slog.Info(
			fmt.Sprintf(
				"Code %d of device %s is in ghost list, not authorized to keylog",
				ke.Code,
				ke.DeviceId,
			),
		)
		return false
	}
	auth = s.baggage.isAuthorized(ke)
	if !auth {
		slog.Info(
			fmt.Sprintf(
				"Baggage of device %s not filled yet, not authorized to keylog",
				ke.DeviceId,
			),
		)
		return false
	}

	return true
}
