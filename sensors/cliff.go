package sensors

import (
	"encoding/binary"

	"github.com/tinogoehlert/gobuki/utils"
)

// CliffADC provides ADC data of PSD sensor, which is facing the floor.
// This value is related with distance between sensor and floor surface. See the datasheet for more detailed information.
// This data also available in the Cliff field of Basic Sensor Data , as a boolean type, processed on the kobuki.
type CliffADC struct {
	Left   uint16
	Center uint16
	Right  uint16
}

// NewCliffADCFromBytes generates a new CliffADC Objects from bytes
func NewCliffADCFromBytes(p []byte) (*CliffADC, error) {
	cliff := CliffADC{}
	cliff.Left = binary.LittleEndian.Uint16(p[:2])
	cliff.Center = binary.LittleEndian.Uint16(p[2:4])
	cliff.Right = binary.LittleEndian.Uint16(p[4:6])
	return &cliff, nil
}

// Equals compares current object with argument
func (c *CliffADC) Equals(v *CliffADC, tolerance int) bool {
	if v == nil {
		return true
	}
	if utils.Cmp(int32(c.Center), int32(v.Center), tolerance) &&
		utils.Cmp(int32(c.Right), int32(v.Right), tolerance) &&
		utils.Cmp(int32(c.Left), int32(v.Left), tolerance) {
		return true
	}
	return false
}

// GetID Gets the Name of the SubPacket
func (c *CliffADC) GetID() byte {
	return 0x05
}

// GetName Gets the Name of the SubPacket
func (c CliffADC) GetName() string {
	return "CliffADC"
}
