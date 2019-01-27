package sensors

import (
	"github.com/tinogoehlert/gobuki/utils"
)

// CurrentWheels represents Current sensor readings of wheel motors in mA.
type CurrentWheels struct {
	LeftMotor  uint8
	RightMotor uint8
}

// NewCurrentWheelsFromBytes generates a new CurrentWheels Objects from bytes
func NewCurrentWheelsFromBytes(p []byte) (*CurrentWheels, error) {
	cw := CurrentWheels{}
	cw.LeftMotor = uint8(p[0])
	cw.RightMotor = uint8(p[1])
	return &cw, nil
}

// Equals compares current object with argument
func (i *CurrentWheels) Equals(v *CurrentWheels, tolerance int) bool {
	if v == nil {
		return true
	}
	if utils.Cmp(int32(i.LeftMotor), int32(i.LeftMotor), tolerance) &&
		utils.Cmp(int32(i.RightMotor), int32(i.RightMotor), tolerance) {
		return true
	}
	return false
}

// GetID Gets the Name of the SubPacket
func (i *CurrentWheels) GetID() byte {
	return 0x06
}

// GetName Gets the Name of the SubPacket
func (i CurrentWheels) GetName() string {
	return "CurrentWheels"
}
