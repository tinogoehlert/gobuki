# go-kobuki

## go driver for kobuki bot based platforms (e.g. turtlebot v2)

![picture missing :(](https://github.com/tinogoehlert/go-kobuki/raw/master/resources/kobuki.webp "Kobuki Bot")


still in developement, interface may change.


## examples

#### read raw gyroscope data
```go
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
```

#### Send sound sequence command

```go
package main

import (
	"github.com/tinogoehlert/go-kobuki/kobuki"
	"github.com/tinogoehlert/go-kobuki/kobuki/commands"
)

func main() {
	// conenct via TCP
    bot, err := kobuki.NewBotTCP("127.0.0.1:3333")
    if err != nil {
        panic(err)
    }
    // connect via serial port
    //bot := kobuki.NewBotSerial("/dev/ttyUSB0")

	bot.Start()
	defer bot.Stop()

	bot.Send(commands.SoundSequenceCmd(commands.ON))
}
```

TODO:

- Make tolerance settings configurable.
- Implement GPIO.

Resouces:

http://kobuki.yujinrobot.com/about2/

http://yujinrobot.github.io/kobuki/enAppendixProtocolSpecification.html

