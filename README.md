# go-kobuki

## go driver for kobuki bot based platforms (e.g. turtlebot v2)

![picture missing :(](https://github.com/tinogoehlert/go-kobuki/raw/master/resources/kobuki.webp "Kobuki Bot")


still in developement, interface may change.


## examples

#### read raw gyroscope data
```go
package main

import (
	"log"
	"time"

	"github.com/tinogoehlert/go-kobuki/kobuki"
	"github.com/tinogoehlert/go-kobuki/kobuki/sensors"
)

func main() {
	// conenct via TCP
    bot, err := kobuki.NewBotTCP("127.0.0.1:3333")
        if err != nil {
        panic(err)
    }
    // connect via serial port
    //bot := kobuki.NewBotSerial("/dev/ttyUSB0")

	bot.On("Gyro", func(data interface{}) {
		bs := data.(*sensors.GyroData)
		log.Printf("Gyro: %f : %F : %f", bs.X, bs.Y, bs.Z)
	})

	bot.Start()
	defer bot.Stop()

	for {
		time.Sleep(1 * time.Second)
	}
}∏
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
- Export internal log stream as channel.

Resouces:

http://kobuki.yujinrobot.com/about2/

http://yujinrobot.github.io/kobuki/enAppendixProtocolSpecification.html

