package conversion

import (
	"encoding/binary"
	"math"
)

func AmbientDegreeCelsius(b []byte) float64 {
	if len(b) != 4 {
		return -1
	}

	return float64(binary.LittleEndian.Uint16(b[2:])) / 128.0
}

func ObjectDegreeCelsius(b []byte) float64 {
	if len(b) != 4 {
		return -1
	}

	return float64(binary.LittleEndian.Uint16(b[:2])) / 128.0
}

func HumidityRelative(b []byte) float64 {
	if len(b) != 4 {
		return -1
	}

	return float64(binary.LittleEndian.Uint16(b[2:])) * 100 / 65536.0
}

func HumidityDegreeCelsius(b []byte) float64 {
	if len(b) != 4 {
		return -1
	}

	return -40 + ((165 * float64(binary.LittleEndian.Uint16(b[:2]))) / 65536.0)
}

func BarometerPressure(b []byte) float64 {
	if len(b) != 6 {
		return -1
	}

	raw := binary.LittleEndian.Uint32(b[2:])
	pressureMask := (int(uint32(raw)) >> 8) & 0x00ffffff

	return float64(pressureMask) / 100.0
}

func OpticalLux(b []byte) float64 {
	if len(b) != 2 {
		return -1
	}
	raw := binary.LittleEndian.Uint16(b)

	exponent := float64((int(raw) & 0xF000) >> 12)
	mantissa := float64((int(raw) & 0x0FFF))

	return mantissa * math.Pow(2, exponent) / 100.0
}
