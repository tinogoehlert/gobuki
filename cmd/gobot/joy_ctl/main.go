package main

import (
	"encoding/binary"
	"flag"
	"math"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
	"gobot.io/x/gobot/platforms/nats"
)

var (
	joyConf  = flag.String("joy", "d4.json", "joystick config string")
	dz       = flag.Uint("dz", 6000, "axis deadzone for x and y")
	natsAddr = flag.String("nats", "127.0.0.1:4222", "-nats host:port")
)

func main() {
	flag.Parse()
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, *joyConf)
	natsAdaptor := nats.NewAdaptor(*natsAddr, 0)

	var (
		vx   float64
		wz   float64
		buff = make([]byte, 16)
	)
	work := func() {
		// joysticks
		stick.On(joystick.LeftX, func(data interface{}) {
			wz = -float64(data.(int16)) / 10000
		})

		stick.On(joystick.L2, func(data interface{}) {
			vx = -float64(triggerValue(data)) / 45000
		})

		stick.On(joystick.R2, func(data interface{}) {
			vx = float64(triggerValue(data)) / 45000
		})

		gobot.Every(10*time.Millisecond, func() {
			binary.LittleEndian.PutUint64(buff[:8], math.Float64bits(vx))
			binary.LittleEndian.PutUint64(buff[8:16], math.Float64bits(wz))
			natsAdaptor.Publish("lovoobot/feedback/move", buff)
		})

	}

	robot := gobot.NewRobot("botControl",
		[]gobot.Connection{joystickAdaptor, natsAdaptor},
		[]gobot.Device{stick},
		work,
	)

	robot.Start()
}

func triggerValue(data interface{}) uint16 {
	return uint16(data.(int16)) + 32768
}

func axisDejitter(data interface{}, deadZone uint) int16 {
	v := abs16(data.(int16))
	if v < int16(deadZone) {
		return 0
	}
	return data.(int16)
}

func abs16(n int16) int16 {
	y := n >> 15
	return int16((n ^ y) - y)
}
