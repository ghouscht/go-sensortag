package sensortag

import (
	"encoding/binary"
)

// Movement sensor
type Movement struct {
	*sensorConfig
}

// empty declarationto check interface
var _ Sensor = &Movement{}

func NewMovement(s *sensorConfig) *Movement {
	return &Movement{s}
}

func (m *Movement) StartNotify(period []byte) (<-chan SensorEvent, error) {
	if err := m.setPeriod(period); err != nil {
		return nil, err
	}

	// enable the sensor
	// bit 0-2 = Gyroscope axis enable
	// bit 3-5 = Accelerometer axis enable
	// bit 6 = Magnetometer axis enable
	// bit 7 = Wake on motion
	// bit 8-9 = Accelerometer range
	// bit 10-15 = not used
	// http://processors.wiki.ti.com/index.php/CC2650_SensorTag_User's_Guide#Configuration_2
	if err := m.enable([]byte{0xFE, 0x40}); err != nil {
		return nil, err
	}

	return m.notify(m.convert)
}

func (m *Movement) convert(data []byte) *[]SensorEvent {
	if len(data) != 18 {
		return nil
	}

	xG := binary.LittleEndian.Uint16(data[0:2])
	yG := binary.LittleEndian.Uint16(data[2:4])
	zG := binary.LittleEndian.Uint16(data[4:6])

	xA := binary.LittleEndian.Uint16(data[6:8])
	yA := binary.LittleEndian.Uint16(data[8:10])
	zA := binary.LittleEndian.Uint16(data[10:12])

	xM := binary.LittleEndian.Uint16(data[12:14])
	yM := binary.LittleEndian.Uint16(data[14:16])
	zM := binary.LittleEndian.Uint16(data[16:18])

	events := []SensorEvent{}

	events = append(events, SensorEvent{
		Name: "Gyroscope",
		Unit: "deg/s",
		X:    pointTo(float64(uint16(xG)) / 128.0),
		Y:    pointTo(float64(uint16(yG)) / 128.0),
		Z:    pointTo(float64(uint16(zG)) / 128.0),
	})

	events = append(events, SensorEvent{
		Name: "Accelerometer",
		Unit: "G",
		X:    pointTo(float64(uint16(xA)) / (8192.0 * 2)), // 4G
		Y:    pointTo(float64(uint16(yA)) / (8192.0 * 2)), // 4G
		Z:    pointTo(float64(uint16(zA)) / (8192.0 * 2)), // 4G
	})

	events = append(events, SensorEvent{
		Name: "Magnetometer",
		Unit: "uT",
		X:    pointTo(float64(uint16(xM))),
		Y:    pointTo(float64(uint16(yM))),
		Z:    pointTo(float64(uint16(zM))),
	})

	return &events
}

// little helper to get a pointer to a float as go does not allow
// to take the address of a numeric constant
func pointTo(f float64) *float64 {
	return &f
}
