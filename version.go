package gobuki

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Version represents kobuki hardware version
type Version struct {
	Patch uint8
	Minor uint8
	Major uint8
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d-%d", v.Major, v.Minor, v.Patch)
}

// FromBytes reads version from bytes
func (v *Version) FromBytes(b []byte) error {
	if len(b) < 3 {
		return errors.New("length missmatch")
	}
	v.Patch = b[0]
	v.Minor = b[1]
	v.Major = b[2]
	return nil
}

// UniqueID Contains Unique Device IDentifier of robot.
// This value is unique for all robot in the world.
// It can be represented by triplet form: <UDID0>-<UDID1>-<UDID2>
type UniqueID struct {
	UDID0 uint32
	UDID1 uint32
	UDID2 uint32
}

// FromBytes reads UID from bytes
func (u *UniqueID) FromBytes(b []byte) (*UniqueID, error) {
	if len(b) < 3 {
		return nil, errors.New("length missmatch")
	}
	return &UniqueID{
		UDID0: binary.LittleEndian.Uint32(b[:4]),
		UDID1: binary.LittleEndian.Uint32(b[4:8]),
		UDID2: binary.LittleEndian.Uint32(b[8:12]),
	}, nil
}

func (u *UniqueID) String() string {
	return fmt.Sprintf("%d-%d-%d", u.UDID0, u.UDID1, u.UDID2)
}
