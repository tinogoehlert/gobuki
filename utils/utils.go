package utils

import "math"

func BitHas(v, test byte) bool {
	return v&test != 0
}

func Checksum(data []byte) byte {
	var cs byte
	for i := 0; i < len(data); i++ {
		cs ^= data[i]
	}
	return cs
}

func ChecksumChk(data []byte) bool {
	var cs uint8
	for i := 0; i < len(data); i++ {
		cs ^= data[i]
	}
	return cs == 0
}

func Abs32(n int32) int {
	y := n >> 31
	return int((n ^ y) - y)
}

func Cmp(a, b int32, offset int) bool {
	return !(Abs32(a-b) > offset)
}

func Cmpu(a, b uint16, offset int) bool {
	return !(Abs32(int32(a)) > offset)
}

func Cmpf(a, b float64, tolerance float64) bool {
	return !(math.Abs(a-b) > tolerance)
}

func CheckMinMax16(val int16) int16 {
	switch {
	case val < math.MinInt16:
		return math.MinInt16
	case val > math.MaxInt16:
		return math.MaxInt16
	default:
		return val
	}
}
