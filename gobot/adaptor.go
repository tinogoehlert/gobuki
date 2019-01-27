package kobuki

import "github.com/tinogoehlert/gobuki"

// Adaptor represents a gobot compatible kobuki adaptor
type Adaptor struct {
	name       string
	bot        *gobuki.Bot
	connectStr string
	connect    func(string) (*gobuki.Bot, error)
}

// NewAdaptorTCP creates a new Adaptor with TCP connection
func NewAdaptorTCP(addr string) *Adaptor {
	return &Adaptor{
		name:       "Kobuki",
		connectStr: addr,
		connect: func(cs string) (*gobuki.Bot, error) {
			return gobuki.NewBotTCP(cs)
		},
	}
}

// NewAdaptorSerial creates a new Adaptor with Serial connection
func NewAdaptorSerial(dev string) *Adaptor {
	return &Adaptor{
		name:       "Kobuki",
		connectStr: dev,
		connect: func(cs string) (*gobuki.Bot, error) {
			return gobuki.NewBotSerial(cs)
		},
	}
}

// Name returns the Adaptor's name
func (k *Adaptor) Name() string {
	return k.name
}

// SetName sets the Adaptor's name
func (k *Adaptor) SetName(name string) {
	k.name = name
}

// Connect initiates a connection to the Kobuki Bot.
func (k *Adaptor) Connect() error {
	bot, err := k.connect(k.connectStr)
	if err == nil {
		k.bot = bot
	}
	return err
}

// Finalize finalizes the Kobuki Adaptor
func (k *Adaptor) Finalize() error {
	k.bot.Stop()
	return nil
}

// Ping what?
func (k *Adaptor) Ping() string {
	return "pong"
}
