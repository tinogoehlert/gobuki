package main

import (
	"log"
	"time"

	gk "github.com/tinogoehlert/gobuki/gobot"
	ks "github.com/tinogoehlert/gobuki/sensors"
	"gobot.io/x/gobot"
)

func main() {
	a := gk.NewAdaptorTCP("127.0.0.1:3333")
	kb := gk.NewDriver(a)

	kb.OnStart(func() {
		kb.SetGyroTolerance(0)
	})

	work := func() {

		kb.OnGyro(func(g *ks.GyroData) {
			log.Printf("%f : %f : %f", g.X, g.Y, g.Z)
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
