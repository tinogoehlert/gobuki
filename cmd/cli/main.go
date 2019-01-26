package main

import (
	"github.com/tinogoehlert/go-kobuki/kobuki"
	"github.com/tinogoehlert/go-kobuki/kobuki/commands"
)

func main() {
	//bot := kobuki.NewBotTCP("172.20.100.30:4161")
	bot := kobuki.NewBotTCP("127.0.0.1:3333")
	//bot := kobuki.NewBotSerial("/dev/ttyUSB0")

	bot.Start()
	defer bot.Stop()

	bot.Send(commands.SoundSequenceCmd(commands.ON))
}
