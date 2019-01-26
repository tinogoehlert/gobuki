package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tinogoehlert/go-kobuki/kobuki"
	"github.com/tinogoehlert/go-kobuki/kobuki/sensors"
)

var (
	tcpAddr    = flag.String("tcp", "", "ip adress:port")
	serialPort = flag.String("serial", "", "serial port")
)

func main() {
	var (
		bot    *kobuki.Bot
		conErr error
	)

	flag.Parse()

	switch {
	case *tcpAddr != "":
		bot, conErr = kobuki.NewBotTCP(*tcpAddr)
	case *serialPort != "":
		bot, conErr = kobuki.NewBotSerial(*serialPort)
	default:
		log.Fatalf("no adress or serial port given")
	}

	if conErr != nil {
		log.Fatalf(conErr.Error())
	}

	bot.Start()
	defer bot.Stop()

	bot.On("Gyro", func(data interface{}) {
		d := data.(*sensors.GyroData)
		fmt.Println(d)
	})

	for {
		log.Println(bot.LogChannel())
	}
}
