package fsphdr

import (
	"fmt"
	"encoding/binary"
)
const (
	FSP_HEADER_SIZE = 112
)
type FspHeader struct {
	NextUnusedSegID uint64

	SpaceID uint32
	//  The highest valid page number, and is incremented when the file is grown. 
	//  However, not all of these pages are initialized (some may be zero-filled)
	Size uint32
	// The highest page number for which the FIL header has been initialized
	FreeLimit uint32
	Flags uint32
	NumberOfPagesUsedFreeFrag uint32

	Free  ListBaseNode
	FreeFrag  ListBaseNode
	FullFrag  ListBaseNode

	FullInodes  ListBaseNode
	FreeInodes  ListBaseNode
}

func NewFspHeader(data []byte, startOffset int) (*FspHeader,error) {
	if len(data) < FSP_HEADER_SIZE {
		return nil, fmt.Errorf("data too short for FSP header: got %d bytes", len(data))
	}

	f := &FspHeader{}
	offset := startOffset

	// SpaceID (4 bytes)
	f.SpaceID = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// Skip Unused (4 bytes)
	offset += 4

	// Size (4 bytes)
	f.Size = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// FreeLimit (4 bytes)
	f.FreeLimit = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// Flags (4 bytes)
	f.Flags = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// NumberOfPagesUsedFreeFrag (4 bytes)
	f.NumberOfPagesUsedFreeFrag = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	lists := []*ListBaseNode{
		&f.Free,
		&f.FreeFrag,
		&f.FullFrag,
		&f.FullInodes,
		&f.FreeInodes,
	}

	for i := 0; i < len(lists); i++ {
		node, err := NewListBaseNode(data, offset)
		if err != nil {
			return nil, err
		}
		*lists[i] = *node
		offset += 16
	}

	return f, nil
}