package main

import (
	"log"

	gk "github.com/tinogoehlert/go-kobuki/gobot"
	"github.com/tinogoehlert/go-kobuki/kobuki/sensors"
	"gobot.io/x/gobot"
)

func main() {
	adapter := gk.NewAdaptorTCP("127.0.0.1:3333")
	kobukiBot := gk.NewDriver(adapter)

	work := func() {
		kobukiBot.OnGyro(func(g *sensors.GyroData) {
			log.Printf("%f:%f:%f", g.X, g.Y, g.Z)
		})
	}

	robot := gobot.NewRobot("kobuki",
		[]gobot.Connection{adapter},
		[]gobot.Device{kobukiBot},
		work,
	)

	robot.Start()
}
