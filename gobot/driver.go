package kobuki

import (
	"fmt"

	"github.com/tinogoehlert/go-kobuki/kobuki/sensors"

	"gobot.io/x/gobot"
)

const (
	// GyroEvent event
	GyroEvent = "Gyro"
	// CliffEvent event
	CliffEvent = "Cliff"
	// WheelsDropEvent event
	WheelsDropEvent = "WheelsDrop"
	// WheelsEncoderEvent event
	WheelsEncoderEvent = "WheelsEncoder"
	// WheelsPWMEvent event
	WheelsPWMEvent = "wheelsPWM"
	// BumperEvent event
	BumperEvent = "Bumper"
	// ButtonsEvent event
	ButtonsEvent = "Button"
	// ChargeStateEvent event
	ChargeStateEvent = "ChargeState"
	// DockingIREvent event
	DockingIREvent = "DockingIR"
	// CliffADCEvent event
	CliffADCEvent = "CliffADC"
	// InertialEvent event
	InertialEvent = "Inertial"
)

// Driver is the interface that describes a driver in gobot
type Driver struct {
	name    string
	adaptor *Adaptor
	gobot.Eventer
}

// NewDriver creates a new Kobuki Bot driver
func NewDriver(a *Adaptor) *Driver {

	d := &Driver{
		name:    "Kobuki",
		adaptor: a,
		Eventer: gobot.NewEventer(),
	}

	d.AddEvent(GyroEvent)
	d.AddEvent(CliffEvent)
	d.AddEvent(CliffADCEvent)
	d.AddEvent(BumperEvent)
	d.AddEvent(ButtonsEvent)
	d.AddEvent(WheelsDropEvent)
	d.AddEvent(WheelsEncoderEvent)
	d.AddEvent(WheelsPWMEvent)
	d.AddEvent(InertialEvent)
	d.AddEvent(DockingIREvent)
	d.AddEvent(ChargeStateEvent)

	return d
}

// Name returns the label for the Driver
func (d *Driver) Name() string {
	return d.name
}

// SetName sets the label for the Driver
func (d *Driver) SetName(s string) {
	d.name = s
}

// Start initiates the Driver
func (d *Driver) Start() error {
	d.adaptor.bot.Start()

	d.adaptor.bot.OnAll(func(name string, data interface{}) {
		d.Publish(name, data)
	})

	return nil
}

// Halt terminates the Driver
func (d *Driver) Halt() error {
	fmt.Println("halt called")
	return nil
}

// Connection returns the Connection associated with the Driver
func (d *Driver) Connection() gobot.Connection {
	return d.adaptor
}

func (d *Driver) OnGyro(f func(*sensors.GyroData)) {
	d.adaptor.bot.On(GyroEvent, func(data interface{}) {
		f(data.(*sensors.GyroData))
	})
}

func (d *Driver) OnCliff(f func(*sensors.Cliff)) {
	d.adaptor.bot.On(CliffEvent, func(data interface{}) {
		f(data.(*sensors.Cliff))
	})
}

func (d *Driver) OnWheelEncoder(f func(*sensors.WheelsEncoder)) {
	d.adaptor.bot.On(WheelsEncoderEvent, func(data interface{}) {
		f(data.(*sensors.WheelsEncoder))
	})
}

func (d *Driver) OnWheelDrop(f func(*sensors.WheelsDrop)) {
	d.adaptor.bot.On(WheelsDropEvent, func(data interface{}) {
		f(data.(*sensors.WheelsDrop))
	})
}

func (d *Driver) OnWheelPWM(f func(*sensors.WheelsPWM)) {
	d.adaptor.bot.On(WheelsPWMEvent, func(data interface{}) {
		f(data.(*sensors.WheelsPWM))
	})
}

func (d *Driver) OnInertial(f func(*sensors.Inertial)) {
	d.adaptor.bot.On(InertialEvent, func(data interface{}) {
		f(data.(*sensors.Inertial))
	})
}

func (d *Driver) OnBumper(f func(*sensors.Bumper)) {
	d.adaptor.bot.On(BumperEvent, func(data interface{}) {
		f(data.(*sensors.Bumper))
	})
}

func (d *Driver) OnButtons(f func(*sensors.Buttons)) {
	d.adaptor.bot.On(ButtonsEvent, func(data interface{}) {
		f(data.(*sensors.Buttons))
	})
}

func (d *Driver) OnChargeState(f func(*sensors.ChargeState)) {
	d.adaptor.bot.On(ChargeStateEvent, func(data interface{}) {
		f(data.(*sensors.ChargeState))
	})
}

func (d *Driver) OnDockingIR(f func(*sensors.DockingIR)) {
	d.adaptor.bot.On(DockingIREvent, func(data interface{}) {
		f(data.(*sensors.DockingIR))
	})
}
