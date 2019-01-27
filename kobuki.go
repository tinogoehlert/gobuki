package gobuki

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/tinogoehlert/gobuki/commands"
	"github.com/tinogoehlert/gobuki/packets"
	"github.com/tinogoehlert/gobuki/sensors"
)

// FeedbackData represents a collection of Feedback Data
type FeedbackData struct {
	BasicSensor  *sensors.BasicSensor
	DockingIR    *sensors.DockingIR
	InerTial     *sensors.Inertial
	CliffADC     *sensors.CliffADC
	CurrenWheels *sensors.CurrentWheels
}

// Callback represents a Callback that carries changed data from kobuki bot
type Callback func(interface{})

// AllCallback represents a Callback that carries changed data from kobuki bot
type AllCallback func(string, interface{})

// ToleranceCfg represents the tolerance settings for specific sensors
type ToleranceCfg struct {
	gyro          float64
	cliffADC      int
	currentWheels int
}

// Bot Represents a Kobuki Bot
type Bot struct {
	conn            io.ReadWriteCloser
	lastFrame       FeedbackData
	currentFrame    FeedbackData
	Gyro            *sensors.Gyro
	callBacks       map[string][]Callback
	allCallback     []AllCallback
	cmdChan         chan commands.Command
	logChan         chan string
	toleranceCfg    ToleranceCfg
	hardwareVersion Version
	firmwareVersion Version
	uid             UniqueID
}

// NewBotTCP creates a new Bot instance and connects to a Kobuki Bot
func NewBotTCP(address string) (*Bot, error) {

	client, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("could not dial targetHost: %s", err.Error())
	}

	client.(*net.TCPConn).SetLinger(2)
	client.(*net.TCPConn).SetKeepAlive(true)
	client.(*net.TCPConn).SetReadBuffer(1)

	bot := Bot{conn: client}
	bot.initBot()
	return &bot, nil
}

// NewBotSerial creates a new Bot instance and connects to a Kobuki Bot via serial port
func NewBotSerial(dev string) (*Bot, error) {
	// Open the port.
	port, err := serial.Open(serial.OpenOptions{
		PortName:        dev,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	})

	if err != nil {
		return nil, fmt.Errorf("could not open serial: %s", err.Error())
	}

	bot := Bot{conn: port}
	bot.initBot()

	return &bot, nil
}

func (k *Bot) initBot() {
	k.Gyro = sensors.NewGyroADC(64, 8)
	k.cmdChan = make(chan commands.Command, 1)
	k.callBacks = make(map[string][]Callback)
	k.logChan = make(chan string)
	k.allCallback = []AllCallback{}
	k.toleranceCfg = ToleranceCfg{
		gyro:     0.01,
		cliffADC: 50,
	}
}

// Stop disconnects from a Bot
func (k *Bot) Stop() {
	k.conn.Close()
}

// Start read loop
func (k *Bot) Start() {
	packetReader := packets.NewPacketReader(bufio.NewReader(k.conn))

	go func() {
		reqTick := time.Tick(10 * time.Second)
		reqCmd := commands.RequestCmd()
		k.conn.Write(reqCmd.Serialize())
		for {
			select {
			case cmd := <-k.cmdChan:
				_, err := k.conn.Write(cmd.Serialize())
				if err != nil {
					fmt.Printf("could not send command: %s\n", err.Error())
				}
			case <-reqTick:
				k.conn.Write(reqCmd.Serialize())
			default:
			}
			b, err := packetReader.ReadData()
			if err != nil {
				k.logChan <- err.Error()
			}
			if b != nil {
				currentFrame := k.parseFrame(b)
				k.hasChangedData(currentFrame)
				k.lastFrame = currentFrame
			}
		}
	}()

	// warmup
	time.Sleep(500 * time.Millisecond)
	return
}

// HardwareVersion returns HardwareVersion string
func (k *Bot) HardwareVersion() string {
	return k.hardwareVersion.String()
}

// FirmwareVersion returns FirmwareVersion string
func (k *Bot) FirmwareVersion() string {
	return k.firmwareVersion.String()
}

// UniqueID returns UID string
func (k *Bot) UniqueID() string {
	return k.uid.String()
}

// SetCliffADCTolerance set tolerance for Cliff ADC
func (k *Bot) SetCliffADCTolerance(t int) {
	k.toleranceCfg.cliffADC = t
}

// SetGyroTolerance set tolerance for gyroscope
func (k *Bot) SetGyroTolerance(t float64) {
	k.toleranceCfg.gyro = t
}

// SetCurrentWheelsTolerance set tolerance for wheels current
func (k *Bot) SetCurrentWheelsTolerance(t int) {
	k.toleranceCfg.currentWheels = t
}

//On registers a new Callback
func (k *Bot) On(event string, cb Callback) {
	if _, ok := k.callBacks[event]; !ok {
		k.callBacks[event] = []Callback{}
	}
	k.callBacks[event] = append(k.callBacks[event], cb)
}

//OnAll registers a new Callback for all events
func (k *Bot) OnAll(cb AllCallback) {
	k.allCallback = append(k.allCallback, cb)
}

// Send sends Command to Bot
func (k *Bot) Send(cmd commands.Command) {
	k.cmdChan <- cmd
}

// LogChannel waits and returns log entry. blocking
func (k *Bot) LogChannel() string {
	return <-k.logChan
}

func (k *Bot) emitEvent(name string, data interface{}) {
	if callbacks, ok := k.callBacks[name]; ok {
		for _, callback := range callbacks {
			callback(data)
		}
	}

	for _, callback := range k.allCallback {
		callback(name, data)
	}
}

func (k *Bot) hasChangedData(current FeedbackData) {
	if current.BasicSensor != nil && k.lastFrame.BasicSensor != nil {
		if !current.BasicSensor.Wheels.Drop.Equals(
			&k.lastFrame.BasicSensor.Wheels.Drop,
		) {
			k.emitEvent("WheelsDrop", &current.BasicSensor.Wheels.Drop)
		}
		if !current.BasicSensor.Wheels.Encoder.Equals(
			&k.lastFrame.BasicSensor.Wheels.Encoder,
		) {
			k.emitEvent("WheelsEncoder", &current.BasicSensor.Wheels.Encoder)
		}
		if !current.BasicSensor.Wheels.PWM.Equals(
			&k.lastFrame.BasicSensor.Wheels.PWM,
		) {
			k.emitEvent("WheelsPWM", &current.BasicSensor.Wheels.PWM)
		}
		if !current.BasicSensor.Cliff.Equals(&k.lastFrame.BasicSensor.Cliff) {
			k.emitEvent("Cliff", &current.BasicSensor.Cliff)
		}
		if !current.BasicSensor.Bumper.Equals(&k.lastFrame.BasicSensor.Bumper) {
			k.emitEvent("Bumper", &current.BasicSensor.Bumper)
		}
		if !current.BasicSensor.Buttons.Equals(&k.lastFrame.BasicSensor.Buttons) {
			k.emitEvent("Buttons", &current.BasicSensor.Buttons)
		}
		if current.BasicSensor.ChargeState != k.lastFrame.BasicSensor.ChargeState {
			k.emitEvent("ChargeState", &current.BasicSensor.ChargeState)
		}
		if current.BasicSensor.BatteryVoltage != k.lastFrame.BasicSensor.BatteryVoltage {
			k.emitEvent("BatteryVoltage", &current.BasicSensor.BatteryVoltage)
		}
	}

	if current.DockingIR != nil {
		if !current.DockingIR.Equals(k.lastFrame.DockingIR) {
			k.emitEvent("DockingIR", current.DockingIR)
		}
	}

	if current.CliffADC != nil {
		if !current.CliffADC.Equals(k.lastFrame.CliffADC, k.toleranceCfg.cliffADC) {
			k.emitEvent("CliffADC", current.CliffADC)
		}
	}

	if current.CurrenWheels != nil {
		if !current.CurrenWheels.Equals(k.lastFrame.CurrenWheels, k.toleranceCfg.currentWheels) {
			k.emitEvent("CurrentWheels", current.CurrenWheels)
		}
	}

	if k.Gyro != nil && k.Gyro.Changed(k.toleranceCfg.gyro) {
		k.emitEvent("Gyro", k.Gyro.GetNewData())
	}

	if current.InerTial != nil {
		if !current.InerTial.Equals(k.lastFrame.InerTial) {
			k.emitEvent("Inertial", current.InerTial)
		}
	}
}

func (k *Bot) parseFrame(buffer []byte) FeedbackData {
	var data FeedbackData
	for offset := 0; (offset + 1) < len(buffer); {
		var (
			subID  = int(buffer[offset])
			subLen = int(buffer[offset+1])
		)

		var (
			err error
		)
		switch {
		case (subID == 0x01 && subLen == 0x0F):
			data.BasicSensor, err = sensors.NewBasicSensorFromBytes(buffer[offset+2 : offset+subLen+2])
		case (subID == 0x03 && subLen == 0x03):
			data.DockingIR, err = sensors.NewDockingIRFromBytes(buffer[offset+2 : offset+subLen+2])
		case (subID == 0x04 && subLen == 0x07):
			data.InerTial, err = sensors.NewInertialFromBytes(buffer[offset+2 : offset+subLen+2])
		case (subID == 0x05 && subLen == 0x06):
			data.CliffADC, err = sensors.NewCliffADCFromBytes(buffer[offset+2 : offset+subLen+2])
		case (subID == 0x06 && subLen == 0x02):
			data.CurrenWheels, err = sensors.NewCurrentWheelsFromBytes(buffer[offset+2 : offset+subLen+2])
		case (subID == 0x0A && subLen == 0x04):
			k.hardwareVersion.FromBytes(buffer[offset+2 : offset+subLen+2])
		case (subID == 0x0B && subLen == 0x04):
			k.firmwareVersion.FromBytes(buffer[offset+2 : offset+subLen+2])
		case (subID == 0x0D):
			k.Gyro.Read(buffer[offset+2 : offset+subLen+2])
		case (subID == 0x13 && subLen == 0x0C):
			k.uid.FromBytes(buffer[offset+2 : offset+subLen+2])
		}

		if err != nil {
			fmt.Println(err)
			k.logChan <- fmt.Sprintf("could not parse data: %s", err.Error())
		}

		offset += subLen + 2
	}
	return data
}
