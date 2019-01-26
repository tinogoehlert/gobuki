package sensors

import (
	"encoding/binary"
)

// Inertial represents Z Axis Gyro of Inertial sensor
type Inertial struct {
	Angle     uint16
	AngleRate uint16
}

// NewInertialFromBytes generates a new Inertial Objects from bytes
func NewInertialFromBytes(p []byte) (*Inertial, error) {
	inertial := Inertial{}
	inertial.Angle = binary.LittleEndian.Uint16(p[:2])
	inertial.AngleRate = binary.LittleEndian.Uint16(p[2:4])
	return &inertial, nil
}

// GetID Gets the Name of the SubPacket
func (i *Inertial) GetID() byte {
	return 0x04
}

func (i *Inertial) Equals(v *Inertial) bool {
	if v == nil {
		return true
	}

	if i.Angle != v.Angle {
		return false
	}
	return false
}

// GetName Gets the Name of the SubPacket
func (i Inertial) GetName() string {
	return "Inertial"
}
