package innodb

import (
	"io"
)
const (
	SIZE_LIST_BASE_NODE = 16
)
type ListBaseNode struct {
	ListLength  uint32

	FirstPageNumber uint32
	LastPageNumber uint32

	FirstOffset     uint16
	LastOffset     uint16
}

func NewListBaseNode(data []byte, offset int) (*ListBaseNode, int, error) {
	if offset + SIZE_LIST_BASE_NODE > len(data) {
		return nil, offset, io.ErrUnexpectedEOF
	}

	length := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	firstPage := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	firstOff := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	lastPage := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	lastOff := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	return &ListBaseNode{
		Length:          length,
		FirstPageNumber: firstPage,
		FirstOffset:     firstOff,
		LastPageNumber:  lastPage,
		LastOffset:      lastOff,
	}, offset, nil
}
