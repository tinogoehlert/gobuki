package gobuki

import (
	"math"

	"github.com/tinogoehlert/gobuki/utils"

	"github.com/tinogoehlert/gobuki/commands"
)

const (
	defaultWheelBase   float64 = 0.23
	defaultWheelRadius float64 = 0.035
	tickToRad          float64 = 0.002436916871363930187454
	epsilon                    = 0.0001
)

// Driver represents a differential driving module
type Mover struct {
	bias float64
}

func NewMover() *Mover {
	return &Mover{
		bias: defaultWheelBase,
	}
}

func (d *Mover) Move(vx, wz float64) commands.Command {
	s, r := d.convert(vx, wz)
	speed := utils.CheckMinMax16(int16(s))
	radius := utils.CheckMinMax16(int16(r))
	return commands.MoveCmd(uint16(speed), uint16(radius))
}

func (d *Mover) convert(vx, wz float64) (speed, radius float64) {
	// vx: in m/s
	// wz: in rad/s

	// Special Case #1 : Straight Run
	if math.Abs(wz) < epsilon {
		radius = 0.0
		speed = 1000.0 * vx
		return
	}

	radius = vx * 1000.0 / wz
	// Special Case #2 : Pure Rotation or Radius is less than or equal to 1.0 mm
	if math.Abs(vx) < epsilon || math.Abs(radius) <= 1.0 {
		speed = 1000.0 * d.bias * wz / 2.0
		radius = 1.0
		return
	}

	// General Case :
	if radius > 0.0 {
		speed = (radius + 1000.0*d.bias/2.0) * wz
	} else {
		speed = (radius - 1000.0*d.bias/2.0) * wz
	}
	return
}
