package gobuki

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Version represents kobuki hardware version
type Version struct {
	Minor uint8
	Major uint8
}

func (v *Version) String() string {
	return fmt.Sprintf("%d-%d", v.Major, v.Minor)
}

// FromBytes reads version from bytes
func (v *Version) FromBytes(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("length missmatch (%d)", len(b))
	}

	v.Major = b[0]
	v.Minor = b[1]
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
func (u *UniqueID) FromBytes(b []byte) error {
	if len(b) < 3 {
		return errors.New("length missmatch")
	}

	u.UDID0 = binary.LittleEndian.Uint32(b[:4])
	u.UDID1 = binary.LittleEndian.Uint32(b[4:8])
	u.UDID2 = binary.LittleEndian.Uint32(b[8:12])

	return nil
}

func (u *UniqueID) String() string {
	return fmt.Sprintf("%d-%d-%d", u.UDID0, u.UDID1, u.UDID2)
}
