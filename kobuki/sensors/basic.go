package sensors

import (
	"bufio"
	"bytes"
	"encoding/binary"

	"github.com/tinogoehlert/go-kobuki/kobuki/utils"
)

// BasicSensor Basic core sensor data.
type BasicSensor struct {
	Timestamp   uint16
	Bumper      Bumper
	Wheels      Wheels
	Cliff       Cliff
	Buttons     Buttons
	ChargeState ChargeState
	// Voltage of battery in 0.1 V. Typically 16.7 V when fully charged
	BatteryVoltage uint8
	Overvoltage    uint8
}

// GetID Gets the ID of the SubPacket
func GetID() byte {
	return 0x01
}

// GetName Gets the Name of the SubPacket
func (b BasicSensor) GetName() string {
	return "BasicSensors"
}

// Bumper flags will be setted when bumper is pressed.
type Bumper struct {
	Left   bool
	Center bool
	Right  bool
}

// FromByte reads Bumper values from byte
func (b *Bumper) FromByte(val byte) {
	b.Right = utils.BitHas(val, 0x01)
	b.Center = utils.BitHas(val, 0x02)
	b.Left = utils.BitHas(val, 0x04)
}

// Equals compares current object with argument
func (b *Bumper) Equals(v *Bumper) bool {
	if b.Center == v.Center &&
		b.Left == v.Left &&
		b.Right == v.Right {
		return true
	}
	return false
}

// Cliff Flag will be setted when cliff is detected
type Cliff struct {
	Left   bool
	Center bool
	Right  bool
}

// FromByte reads Cliff values from byte
func (c *Cliff) FromByte(val byte) {
	c.Right = utils.BitHas(val, 0x01)
	c.Center = utils.BitHas(val, 0x02)
	c.Left = utils.BitHas(val, 0x04)
}

// Equals compares current object with argument
func (c *Cliff) Equals(v *Cliff) bool {
	if c.Center == v.Center &&
		c.Left == v.Left &&
		c.Right == v.Right {
		return true
	}
	return false
}

// WheelsDrop Flag will be set when wheel is dropped
type WheelsDrop struct {
	Left  bool
	Right bool
}

func (wd *WheelsDrop) FromByte(val byte) {
	wd.Right = utils.BitHas(val, 0x01)
	wd.Left = utils.BitHas(val, 0x02)
}

// Equals compares current object with argument
func (wd *WheelsDrop) Equals(v *WheelsDrop) bool {
	if wd.Left == v.Left &&
		wd.Right == v.Right {
		return true
	}
	return false
}

// WheelsEncoder Accumulated encoder data of left and right wheels in ticks
// Increments of this value means forward direction
// It circulates from 0 to 65535
type WheelsEncoder struct {
	Left  uint16
	Right uint16
}

// Equals compares current object with argument
func (we *WheelsEncoder) Equals(v *WheelsEncoder) bool {
	if we.Left == v.Left &&
		we.Right == v.Right {
		return true
	}
	return false
}

// WheelsPWM value that applied to left and right wheel motor
// This data should be converted signed type to represent correctly
// Negative sign indicates backward direction
type WheelsPWM struct {
	Left  int8
	Right int8
}

// Equals compares current object with argument
func (wp *WheelsPWM) Equals(v *WheelsPWM) bool {
	if wp.Left == v.Left &&
		wp.Right == v.Right {
		return true
	}
	return false
}

// Wheels collects wheel data
type Wheels struct {
	Drop    WheelsDrop
	Encoder WheelsEncoder
	PWM     WheelsPWM
}

// Equals compares current object with argument
func (w *Wheels) Equals(v *Wheels) bool {
	if w.Drop.Equals(&v.Drop) && w.Encoder.Equals(&v.Encoder) && w.PWM.Equals(&v.PWM) {
		return true
	}
	return false
}

// Buttons Flag will be setted when button is pressed
type Buttons struct {
	Button0 bool
	Button1 bool
	Button2 bool
}

// FromByte reads Button values from byte
func (b *Buttons) FromByte(val byte) {
	b.Button0 = utils.BitHas(val, 0x01)
	b.Button1 = utils.BitHas(val, 0x02)
	b.Button2 = utils.BitHas(val, 0x02)
}

// Equals compares current object with argument
func (b *Buttons) Equals(v *Buttons) bool {
	if b.Button0 == v.Button0 &&
		b.Button1 == v.Button1 &&
		b.Button2 == v.Button2 {
		return true
	}
	return false
}

type ChargeState int8

const (
	Discharging     ChargeState = 0
	DockingCharged  ChargeState = 2
	DockingCharging ChargeState = 6
	AdapterCharged  ChargeState = 18
	AdapterCharging ChargeState = 22
)

// NewBasicSensorFromBytes creates BasicSensors Object from byte
func NewBasicSensorFromBytes(p []byte) (*BasicSensor, error) {
	buf := bufio.NewReader(bytes.NewReader(p))
	bs := &BasicSensor{}
	cutError := func(b byte, err error) byte {
		return b
	}
	binary.Read(buf, binary.LittleEndian, &bs.Timestamp)
	bs.Bumper.FromByte(cutError(buf.ReadByte()))
	bs.Wheels.Drop.FromByte(cutError(buf.ReadByte()))
	bs.Cliff.FromByte(cutError(buf.ReadByte()))
	binary.Read(buf, binary.LittleEndian, &bs.Wheels.Encoder.Left)
	binary.Read(buf, binary.LittleEndian, &bs.Wheels.Encoder.Right)
	bs.Buttons.FromByte(cutError(buf.ReadByte()))
	binary.Read(buf, binary.LittleEndian, &bs.ChargeState)
	bs.Overvoltage = cutError(buf.ReadByte())
	return bs, nil
}
