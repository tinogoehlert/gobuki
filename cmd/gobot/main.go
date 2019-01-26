package main

import (
	"log"
	"time"

	gk "github.com/tinogoehlert/go-kobuki/gobot"
	ks "github.com/tinogoehlert/go-kobuki/kobuki/sensors"
	"gobot.io/x/gobot"
)

func main() {
	a := gk.NewAdaptorTCP("127.0.0.1:3333")
	kb := gk.NewDriver(a)

	work := func() {
		kb.OnWheelsCurrent(func(w *ks.CurrentWheels) {
			log.Printf("%d:%d", w.LeftMotor, w.RightMotor)
		})

		gobot.Every(1*time.Minute, func() {
			kb.PlaySoundSequence(gk.SoundOn)
		})
	}

	robot := gobot.NewRobot("kobuki",
		[]gobot.Connection{a},
		[]gobot.Device{kb},
		work,
	)

	robot.Start()
}
