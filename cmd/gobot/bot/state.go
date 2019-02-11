package main

import "gobot.io/x/gobot/platforms/keyboard"

type botState struct {
	Velocity  float64
	Raduis    float64
	IsMovable bool
}

func (b *botState) MoveByKey(key int) (float64, float64) {
	if b.IsMovable == false {
		return 0, 0
	}
	switch key {
	case keyboard.W:
		return b.Velocity, 0
	case keyboard.S:
		return -b.Velocity, 0
	case keyboard.A:
		return 0, 1
	case keyboard.D:
		return 0, -1
	case keyboard.Q:
		return b.Velocity, 0.5
	case keyboard.E:
		return b.Velocity, -0.5
	case keyboard.X:
		return b.Velocity, 0.5
	case keyboard.Y:
		return -b.Velocity, -0.5
	case keyboard.U:
		if b.Velocity <= 1 {
			b.Velocity += 0.05
		}
	case keyboard.I:
		if b.Velocity >= 0 {
			b.Velocity -= 0.05
		}
	}
	return 0, 0
}
