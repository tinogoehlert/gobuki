package sensors

import "github.com/tinogoehlert/go-kobuki/kobuki/utils"

type DockingSignal struct {
	None       bool
	FarLeft    bool
	FarCenter  bool
	FarRight   bool
	NearLeft   bool
	NearCenter bool
	NearRight  bool
}

func (ds *DockingSignal) FromByte(val byte) {
	ds.NearLeft = utils.BitHas(val, 0x01)
	ds.NearCenter = utils.BitHas(val, 0x02)
	ds.NearRight = utils.BitHas(val, 0x04)
	ds.FarCenter = utils.BitHas(val, 0x08)
	ds.NearLeft = utils.BitHas(val, 0x10)
	ds.NearRight = utils.BitHas(val, 0x20)
	ds.None = (val == 0x00)
}

//DockingIR Collection of Docking IR Signals for Right, Center and Left
type DockingIR struct {
	Left   DockingSignal
	Center DockingSignal
	Right  DockingSignal
}

//NewDockingIRFromBytes generates a new DockingIR Objects from bytes
func NewDockingIRFromBytes(p []byte) (*DockingIR, error) {
	dockIR := DockingIR{}
	dockIR.Right.FromByte(p[0])
	dockIR.Center.FromByte(p[1])
	dockIR.Left.FromByte(p[2])

	return &dockIR, nil
}

func (ds *DockingSignal) equalsSignal(v *DockingSignal) bool {
	if ds.FarCenter != v.FarCenter ||
		ds.FarLeft != v.FarLeft ||
		ds.FarRight != v.FarRight ||
		ds.NearCenter != v.NearCenter ||
		ds.NearLeft != v.NearLeft ||
		ds.NearRight != v.NearRight ||
		ds.None != v.None {
		return true
	}
	return false
}

// Equals compares current object with argument
func (d *DockingIR) Equals(v *DockingIR) bool {
	if v == nil {
		return true
	}
	if d.Center.equalsSignal(&v.Center) &&
		d.Right.equalsSignal(&v.Right) &&
		d.Left.equalsSignal(&v.Left) {
		return true
	}
	return false
}

// GetID Gets the Name of the SubPacket
func (d *DockingIR) GetID() byte {
	return 0x03
}

// GetName Gets the Name of the SubPacket
func (d *DockingIR) GetName() string {
	return "DockingIR"
}
