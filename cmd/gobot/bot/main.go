package main

import (
	"encoding/binary"
	"flag"
	"log"
	"math"

	"github.com/tinogoehlert/gobuki"

	"gobot.io/x/gobot/platforms/nats"

	gokubi_driver "github.com/tinogoehlert/gobuki/gobot"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
)

var (
	tcpAddr    = flag.String("tcp", "", "-tcp host:port")
	serialPort = flag.String("serial", "", "-serial port")
	natsAddr   = flag.String("nats", "127.0.0.1:4222", "-nats host:port")
)

func main() {
	flag.Parse()

	var gobukiAdaptor *gokubi_driver.Adaptor
	switch {
	case *tcpAddr != "":
		gobukiAdaptor = gokubi_driver.NewAdaptorTCP("192.168.1.176:4161")
	case *serialPort != "":
		gobukiAdaptor = gokubi_driver.NewAdaptorSerial(*serialPort)
	default:
		log.Fatalf("no adress or serial port given")
	}

	gobukiDriver := gokubi_driver.NewDriver(gobukiAdaptor)
	natsAdaptor := nats.NewAdaptor(*natsAddr, 0)
	keys := keyboard.NewDriver()
	state := botState{Velocity: 0.3, IsMovable: true}

	work := func() {
		gobukiDriver.On(gokubi_driver.FeedbackEvent, func(data interface{}) {
			event := data.(gokubi_driver.Feedback)
			msg, err := gobuki.EncodeGob(event)
			if err != nil {
				log.Println(err.Error())
			}
			natsAdaptor.Publish("lovoobot/feedback/"+event.Name, msg)
		})

		natsAdaptor.On("lovoobot/feedback/move_raw", func(msg nats.Message) {
			if len(msg.Data) == 4 {
				s := int16(binary.LittleEndian.Uint16(msg.Data[0:2]))
				r := int16(binary.LittleEndian.Uint16(msg.Data[2:4]))
				gobukiDriver.MoveRaw(s, r)
			}
		})

		natsAdaptor.On("lovoobot/feedback/move", func(msg nats.Message) {
			if len(msg.Data) == 16 {
				vx := math.Float64frombits(binary.LittleEndian.Uint64(msg.Data[0:8]))
				wz := math.Float64frombits(binary.LittleEndian.Uint64(msg.Data[8:16]))
				if vx != 0 || wz != 0 {
					gobukiDriver.Move(vx, wz)
				}
			}
		})

		keys.On(keyboard.Key, func(data interface{}) {
			ks := data.(keyboard.KeyEvent)
			gobukiDriver.Move(state.MoveByKey(ks.Key))
		})
	}

	robot := gobot.NewRobot("kobuki",
		[]gobot.Connection{gobukiAdaptor, natsAdaptor},
		[]gobot.Device{gobukiDriver, keys},
		work,
	)

	robot.Start()
}
