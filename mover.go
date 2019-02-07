package gobuki

import (
	"math"
	"time"

	"github.com/tinogoehlert/gobuki/utils"

	"github.com/tinogoehlert/gobuki/commands"
)

const (
	defaultWheelBase   float64 = 0.23
	defaultWheelRadius float64 = 0.035
	tickToRad          float64 = 0.002436916871363930187454
	epsilon            float64 = 0.0001
)

// Mover represents a differential driving module
type Mover struct {
	bias              float64
	bot               *Bot
	lastTickLeft      uint16
	lastTickRight     uint16
	lastTimestamp     uint16
	lastDiffTime      float64
	lastVelocityLeft  float64
	lastVelocityRight float64
	lastRadLeft       float64
	lastRadRight      float64
	initLeft          bool
	initRight         bool
	speed             int16
	radius            int16
	lastMovedTime     time.Time
}

// NewMover Creates a new Mover
func NewMover(bot *Bot) *Mover {
	return &Mover{
		bias: defaultWheelBase,
		bot:  bot,
	}
}

// Tick gets called when new Feedback Data is ready
func (m *Mover) Tick(data FeedbackData) *commands.Command {
	if data.BasicSensor == nil {
		return nil
	}
	m.update(
		data.BasicSensor.Timestamp,
		data.BasicSensor.Wheels.Encoder.Left,
		data.BasicSensor.Wheels.Encoder.Right,
	)
	cmd := commands.MoveCmd(uint16(m.speed), uint16(m.radius))
	if time.Now().Sub(m.lastMovedTime) > (200 * time.Millisecond) {
		m.speed = 0
		m.radius = 0
	}
	return &cmd
}

// Move throws new move values into move command buffer jo
func (m *Mover) Move(vx, wz float64) {
	s, r := m.convert(vx, wz)
	m.speed = utils.CheckMinMax16(int16(s))
	m.radius = utils.CheckMinMax16(int16(r))
	m.lastMovedTime = time.Now()
}

func (m *Mover) MoveRaw(speed, radius int16) {
	m.speed = utils.CheckMinMax16(int16(speed))
	m.radius = utils.CheckMinMax16(int16(radius))
	m.lastMovedTime = time.Now()
}

func (m *Mover) update(timestamp, leftEncoder, rightEncoder uint16) {
	var (
		leftDiffTicks  float64
		rightDiffTicks float64
		currTickLeft   uint16
		currTickRight  uint16
		currTimestamp  uint16
	)

	currTimestamp = timestamp
	currTickLeft = leftEncoder

	if !m.initLeft {
		m.lastTickLeft = currTickLeft
		m.initLeft = true
	}

	leftDiffTicks = float64((currTickLeft - m.lastTickLeft) & 0xffff)
	m.lastTickLeft = currTickLeft
	m.lastRadLeft += tickToRad * leftDiffTicks

	currTickRight = rightEncoder
	if !m.initRight {
		m.lastTickLeft = currTickLeft
		m.initRight = true
	}

	rightDiffTicks = float64((currTickRight - m.lastTickRight) & 0xffff)
	m.lastTickLeft = currTickRight
	m.lastRadRight += tickToRad * rightDiffTicks

	if currTimestamp != m.lastTimestamp {
		m.lastDiffTime = float64((currTimestamp-m.lastTimestamp)&0xffff) / 1000.0
		m.lastTimestamp = currTimestamp
		m.lastVelocityLeft = (tickToRad * leftDiffTicks) / m.lastDiffTime
		m.lastVelocityRight = (tickToRad * rightDiffTicks) / m.lastDiffTime
	}
}

func (m *Mover) convert(vx, wz float64) (speed, radius float64) {
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
		speed = 1000.0 * m.bias * wz / 2.0
		radius = 1.0
		return
	}

	// General Case :
	if radius > 0.0 {
		speed = (radius + 1000.0*m.bias/2.0) * wz
	} else {
		speed = (radius - 1000.0*m.bias/2.0) * wz
	}
	return
}
