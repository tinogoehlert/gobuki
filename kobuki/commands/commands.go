package commands

import (
	"encoding/binary"

	"github.com/tinogoehlert/go-kobuki/kobuki/utils"
)

const (
	header0 byte = 0xAA
	header1 byte = 0x55
)

// Command represents a Robot command packet
type Command struct {
	ID   byte
	Data []byte
}

// Serialize packs Command Data in kobuki protocol compatible byte buffer
func (c *Command) Serialize() []byte {
	dlen := byte(len(c.Data))
	buff := make([]byte, 0, dlen+5)
	buff = append(buff, header0, header1, dlen+2) // add protocol header
	buff = append(buff, c.ID, dlen)               // add command control data
	buff = append(buff, c.Data...)                // add payload
	buff = append(buff, utils.Checksum(buff[2:])) // add checksum
	return buff
}

// MoveCmd Creates a Move Command (Wheel Motor Control)
func MoveCmd(Speed, Rotation int16) Command {
	cmd := Command{
		ID:   0x01,
		Data: make([]byte, 4),
	}
	binary.LittleEndian.PutUint16(cmd.Data, uint16(Speed))
	binary.LittleEndian.PutUint16(cmd.Data[2:], uint16(Rotation))
	return cmd
}

// SoundCmd create a Command to play a sound
func SoundCmd(note uint16, duration byte) Command {
	cmd := Command{
		ID:   0x03,
		Data: make([]byte, 3),
	}
	binary.LittleEndian.PutUint16(cmd.Data, note)
	cmd.Data[2] = duration
	return cmd
}

// SoundSequence represents a constant sequence id
type SoundSequence byte

const (
	// ON Plays Sound Sequence "ON"
	ON SoundSequence = 0x00
	// OFF Plays Sound Sequence "OFF"
	OFF SoundSequence = 0x01
	// RECHARGE Plays Sound Sequence "RECHARGE"
	RECHARGE SoundSequence = 0x02
	// BUTTON Plays Sound Sequence "BUTTON"
	BUTTON SoundSequence = 0x03
	// ERROR Plays Sound Sequence "ERROR"
	ERROR SoundSequence = 0x04
	// CLEANINGSTART Plays Sound Sequence "CLEANING START" .. wtf?
	CLEANINGSTART SoundSequence = 0x05
	// CLEANINGEND Plays Sound Sequence "CLEANING END"
	CLEANINGEND SoundSequence = 0x06
)

// SoundSequenceCmd creates command to play given sequence
func SoundSequenceCmd(SoundSequence SoundSequence) Command {
	cmd := Command{
		ID:   0x04,
		Data: make([]byte, 1),
	}
	cmd.Data[0] = byte(SoundSequence)
	return cmd
}
