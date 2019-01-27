# gobuki

## go driver for kobuki bot based platforms (e.g. turtlebot v2)

![picture missing :(](https://github.com/tinogoehlert/gobuki/raw/master/resources/kobuki.webp "Kobuki Bot")


still in developement, interface may change.


## examples

#### read raw gyroscope data
```go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tinogoehlert/gobuki"
	"github.com/tinogoehlert/gobuki/sensors"
)

var (
	tcpAddr    = flag.String("tcp", "", "ip adress:port")
	serialPort = flag.String("serial", "", "serial port")
)

func main() {
	var (
		bot    *gobuki.Bot
		conErr error
	)

	flag.Parse()

	switch {
	case *tcpAddr != "":
		bot, conErr = gobuki.NewBotTCP(*tcpAddr)
	case *serialPort != "":
		bot, conErr = gobuki.NewBotSerial(*serialPort)
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
```

#### Send sound sequence command

```go
package main

import (
	"github.com/tinogoehlert/gobuki"
	"github.com/tinogoehlert/gobuki/commands"
)

func main() {
	bot, err := gobuki.NewBotTCP("127.0.0.1:3333")
    if err != nil {
        panic(err)
    }

	bot.Start()
	defer bot.Stop()

	bot.Send(commands.SoundSequenceCmd(commands.ON))
}
```

#### use kobuki bot with gobot

```go
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
```

TODO:

- Implement GPIO.
- Add differential drive helpers.

Resouces:

http://kobuki.yujinrobot.com/about2/

http://yujinrobot.github.io/kobuki/enAppendixProtocolSpecification.html

