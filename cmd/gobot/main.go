package main

import (
	"time"

	"github.com/tinogoehlert/gobuki/commands"

	gk "github.com/tinogoehlert/gobuki/gobot"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
)

func main() {
	//a := gk.NewAdaptorTCP("192.168.1.176:4161")
	a := gk.NewAdaptorSerial("/dev/ttyUSB0")
	kb := gk.NewDriver(a)
	keys := keyboard.NewDriver()

	kb.OnStart(func() {
		kb.PlaySoundSequence(commands.ON)
	})

	work := func() {
		keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)
			switch key.Key {
			case keyboard.W:
				kb.Move(1, 0)
			case keyboard.A:
				kb.Move(0, 1)
			case keyboard.S:
				kb.Move(-1, 0)
			case keyboard.D:
				kb.Move(-1, 0)
			}
		})
		gobot.Every(1*time.Second, func() {
			kb.Move(0, 0)
		})
	}

	robot := gobot.NewRobot("kobuki",
		[]gobot.Connection{a},
		[]gobot.Device{kb, keys},
		work,
	)

	robot.Start()
}
