package packets

import (
	"bufio"
	"fmt"

	"github.com/tinogoehlert/go-kobuki/kobuki/utils"
)

const (
	// ErrPayloadLen payload length does not match real lenght
	ErrPayloadLen = "payload length missmatch"
)

const (
	magitHeaderA byte = 0xAA
	magitHeaderB byte = 0x55
)

// PacketReader reads kobuki protocol packages from a stream
type PacketReader struct {
	data       []byte
	buffReader *bufio.Reader
	readIndex  int64
}

// NewPacketReader creates a new PacketReader Instance
func NewPacketReader(buffReader *bufio.Reader) *PacketReader {
	return &PacketReader{
		data:       make([]byte, 256),
		buffReader: buffReader,
	}
}

func (r *PacketReader) ReadData() ([]byte, error) {

	headerA, _ := r.buffReader.ReadByte()
	if headerA != magitHeaderA {
		return nil, nil
	}

	headerB, _ := r.buffReader.ReadByte()
	if headerB != magitHeaderB {
		return nil, nil
	}

	payloadLen, err := r.buffReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("could not read header bytes: %s", err.Error())
	}

	n, err := r.buffReader.Read(r.data[:payloadLen])
	if err != nil {
		return nil, fmt.Errorf("could not read payloads: %s", err.Error())
	}

	// second read attempt
	if n != int(payloadLen) {
		n2, err := r.buffReader.Read(r.data[n:payloadLen])
		n += n2
		if err != nil {
			return nil, fmt.Errorf("could not read payloads: %s", err.Error())
		}
	}

	if n != int(payloadLen) {
		return nil, fmt.Errorf(ErrPayloadLen)
	}

	cs, err := r.buffReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("could not read checksum: %s", err.Error())
	}

	tb := []byte{payloadLen}
	tb = append(tb, r.data[:payloadLen]...)
	if cs != utils.Checksum(tb) {
		return nil, fmt.Errorf("checksum missmatch")
	}
	return r.data[:payloadLen], nil
}
