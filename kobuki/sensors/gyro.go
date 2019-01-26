package sensors

import (
	"encoding/binary"
	"fmt"

	"github.com/tinogoehlert/go-kobuki/kobuki/utils"
)

// GyroData stores XYZ Data from Gyro
type GyroData struct {
	X float64
	Y float64
	Z float64
}

const digitToDPS float64 = 0.00875

// ToVelocity rotates 90 degree counterclockwise about z-axis.
// Sensing axis of 3d gyro is not match with robot. So, below conversion will needed.
func (gd *GyroData) ToVelocity(x, y, z uint16) {
	gd.X = -digitToDPS * float64(y)
	gd.Y = digitToDPS * float64(x)
	gd.Z = digitToDPS * float64(z)
}

// Gyro Raw ADC data of digital 3D gyro; L3G4200D
// Due to difference of acquisition rate and update rate, 2-3 data will be arrived at once.
// Digit to deg/s ratio is 0.00875, it comes from datasheet of 3d gyro.
type Gyro struct {
	FrameID          uint8
	Data             []GyroData
	newest           *GyroData
	last             *GyroData
	delta            GyroData
	tolerance        GyroData
	bucketSize       int
	readCount        int
	calibrateSamples int
	calibrated       bool
}

// NewGyroADC generates a new Gyro Objects from bytes
func NewGyroADC(bucketSize, calibrateSamples int) *Gyro {
	gyro := Gyro{
		Data:             make([]GyroData, bucketSize),
		bucketSize:       bucketSize,
		calibrateSamples: calibrateSamples,
	}
	return &gyro
}

// Changed check if values changed significantly
func (g *Gyro) Changed(tolerance float64) bool {
	if g.last == nil || g.newest == nil {
		return false
	}

	if !utils.Cmpf(g.newest.X, g.last.X, tolerance) ||
		!utils.Cmpf(g.newest.Z, g.last.Z, tolerance) ||
		!utils.Cmpf(g.newest.Z, g.last.Z, tolerance) {
		return true
	}

	return false
}

// GetNewData returns newest gyro data
func (g *Gyro) GetNewData() *GyroData {
	return &GyroData{
		X: g.newest.X,
		Y: g.newest.Y,
		Z: g.newest.Z,
	}
}

// Read reads raw gyro data from bytes
func (g *Gyro) Read(p []byte) error {
	if len(p) < 8 {
		return fmt.Errorf("gyroscope currupt read")
	}

	g.FrameID = p[0]
	dataLen := int(p[1])

	if dataLen <= 0 || dataLen > len(p) {
		return fmt.Errorf("gyroscope length missmatch: %d > %d", dataLen, len(p))
	}

	ec := int(dataLen / 3)

	for i := 0; i < ec; i++ {
		if g.readCount >= g.bucketSize {
			g.readCount = 0
		}
		v := p[2+i*6:]
		g.Data[g.readCount].ToVelocity(
			binary.LittleEndian.Uint16(v),
			binary.LittleEndian.Uint16(v[2:]),
			binary.LittleEndian.Uint16(v[4:]),
		)

		g.last = g.newest
		if g.calibrated {
			g.Data[g.readCount].X = (g.Data[g.readCount].X - g.delta.X) * digitToDPS
			g.Data[g.readCount].Y = (g.Data[g.readCount].Y - g.delta.Y) * digitToDPS
			g.Data[g.readCount].Z = (g.Data[g.readCount].Y - g.delta.Z) * digitToDPS
			g.newest = &g.Data[g.readCount]
		}

		g.readCount++
	}

	if g.readCount >= g.calibrateSamples && !g.calibrated {
		g.Calibrate(g.calibrateSamples)
	}

	return nil
}

// Calibrate algorithm
func (g *Gyro) Calibrate(samples int) {

	if samples > g.readCount || g.calibrated {
		return
	}

	// Reset values
	var (
		sumX   float64
		sumY   float64
		sumZ   float64
		sigmaX float64
		sigmaY float64
		sigmaZ float64
	)

	// Read n-samples
	for i := 0; i < samples; i++ {
		sumX += g.Data[i].X
		sumY += g.Data[i].Y
		sumZ += g.Data[i].Z

		sigmaX += g.Data[i].X * g.Data[i].X
		sigmaY += g.Data[i].Y * g.Data[i].Y
		sigmaZ += g.Data[i].Z * g.Data[i].Z
	}

	fsamples := float64(samples)
	// Calculate delta vectors
	g.delta.X = sumX / fsamples
	g.delta.Y = sumY / fsamples
	g.delta.Y = sumZ / fsamples

	g.calibrated = true
}

// GetID Gets the Name of the SubPacket
func (g *Gyro) GetID() byte {
	return 0x0D
}

// GetName Gets the Name of the SubPacket
func (g Gyro) GetName() string {
	return "Gyro"
}
