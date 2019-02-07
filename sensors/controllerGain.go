package sensors

import (
	"encoding/binary"
	"fmt"
)

// ControllerGain Contains PID gain of wheel velocity controller of robot.
type ControllerGain struct {
	Type uint8
	P    uint32
	I    uint32
	D    uint32
}

// FromBytes Objects from bytes
func (c *ControllerGain) FromBytes(p []byte) error {

	if len(p) < 13 {
		return fmt.Errorf("ControllerGain len missmatch (%d)", len(p))
	}
	c.Type = uint8(p[0])
	c.P = binary.LittleEndian.Uint32(p[1:5])
	c.P = binary.LittleEndian.Uint32(p[5:9])
	c.P = binary.LittleEndian.Uint32(p[9:13])
	return nil
}

// GetID Gets the Name of the SubPacket
func (c *ControllerGain) GetID() byte {
	return 0x01
}

// GetName Gets the Name of the SubPacket
func (c ControllerGain) GetName() string {
	return "ControllerGain"
}
