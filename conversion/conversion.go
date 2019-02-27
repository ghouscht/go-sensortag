package conversion

import (
	"encoding/binary"
)

func BarometerPressure(b []byte) float64 {
	if len(b) != 6 {
		return -1
	}

	raw := binary.LittleEndian.Uint32(b[2:])
	pressureMask := (int(uint32(raw)) >> 8) & 0x00ffffff

	return float64(pressureMask) / 100.0
}
