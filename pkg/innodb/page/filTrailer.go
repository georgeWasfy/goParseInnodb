package innodb

import (
	"encoding/binary"
	"fmt"
)
const (
	// File Trailer
	FIL_Trailer_START = 16376
	FIL_TRAILER_SIZE  = 8

)

type FilTrailer struct {
	Checksum uint32
	Low32lsn uint32
}

func NewFilTrailer(data []byte) (*FilTrailer,error) {
		if len(data) < FIL_TRAILER_SIZE {
		return nil, fmt.Errorf("data too short for fil trailer: got %d bytes", len(data))
	}

	f := &FilTrailer{}

	offset := 0

	// Checksum (4 bytes)
	f.Checksum = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// low32lsn (4 bytes)
	f.Low32lsn = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	return f, nil

}