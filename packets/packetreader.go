package packets

import (
	"bufio"
	"fmt"

	"github.com/tinogoehlert/gobuki/utils"
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

	payloadLenByte, err := r.buffReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("could not read header bytes: %s", err.Error())
	}

	var payloadLen = int(payloadLenByte)
	for i := 0; i < payloadLen; {
		nb, err := r.buffReader.Read(r.data[i:payloadLen])
		i += nb
		if err != nil {
			return nil, fmt.Errorf("could not read payloads: %s", err.Error())
		}
	}

	checksumR, err := r.buffReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("could not read checksum: %s", err.Error())
	}

	checksumL := (payloadLenByte ^ utils.Checksum(r.data[:payloadLen]))
	if checksumL != checksumR {
		return nil, fmt.Errorf("checksum missmatch (%d != %d)", checksumL, checksumR)
	}

	return r.data[:payloadLen], nil
}
