package innodb

import (
	"encoding/binary"
	"fmt"
)

const (
	// File header
	FIL_HEADER_START = 0
	FIL_HEADER_SIZE  = 38

	// Page header
	PAGE_HEADER_START = FIL_HEADER_START + FIL_HEADER_SIZE
)
type FilHeader struct {
	LastModifiedLsn uint64
	FlushLsn uint64
	Checksum uint32
	Offset  uint32
	/**
   * Pointers to the logical previous and next page for this page type are stored in the header.
   * This allows doubly-linked lists of pages to be built, and this is used for INDEX pages to link all pages at
   * the same level, which allows for e.g. full index scans to be efficient.
   * Many page types do not use these fields.
   */
	PreviousPage uint32
	NextPage uint32
	SpaceID uint32
	PageType PageType
}

func NewFilHeader(data []byte) (*FilHeader,error) {
		if len(data) < FIL_HEADER_SIZE {
		return nil, fmt.Errorf("data too short for fil header: got %d bytes", len(data))
	}

	f := &FilHeader{}

	offset := 0

	// Checksum (4 bytes)
	f.Checksum = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// PageNumber / Offset (4 bytes)
	f.Offset = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// PreviousPage (4 bytes)
	f.PreviousPage = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// NextPage (4 bytes)
	f.NextPage = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// LastModifiedLsn (8 bytes)
	f.LastModifiedLsn = binary.BigEndian.Uint64(data[offset:])
	offset += 8

	// PageType (2 bytes) — assuming PageType is uint16
	f.PageType = PageType(binary.BigEndian.Uint16(data[offset:]))
	offset += 2

	// FlushLsn (8 bytes)
	f.FlushLsn = binary.BigEndian.Uint64(data[offset:])
	offset += 8

	// SpaceID (4 bytes)
	f.SpaceID = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	return f, nil

}