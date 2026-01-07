package fsphdr

import (
	"io"
	"encoding/binary"
)

type XdesState uint32

const (
	XdesFree      XdesState = 0
	XdesFreeFrag  XdesState = 1
	XdesFullFrag  XdesState = 2
	XdesFseg      XdesState = 3
)
type XdesEntry struct{
	FileSegmentID uint32
	// List node for XDES list: Pointers to previous and next extents in a 
	// doubly-linked extent descriptor list.
	XdesListNode ListNode
	State XdesState
	// Page State Bitmap: A bitmap of 2 bits per page in the extent (64 x 2 = 128 bits, or 16 bytes). 
	// The first bit indicates whether the page is free. 
	// The second bit is reserved to indicate whether the page is clean (has no un-flushed data), 
	// but this bit is currently unused and is always set to 1.
	PageBitmap  [16]byte
} 

func NewXdesEntry(data []byte, offset int) (*XdesEntry, error) {
	start := offset

	if start + XDES_ENTRY_SIZE > len(data) {
		return nil, io.ErrUnexpectedEOF
	}

	// File Segment ID (8 bytes)
	fsegID := binary.BigEndian.Uint32(data[start:])
	start += 8

	// List node (12 bytes)
	listNode, err := NewListNode(data, start)
	if err != nil {
		return nil, err
	}
	start += 12

	// State (4 bytes)
	state := XdesState(binary.BigEndian.Uint32(data[start:]))
	start += 4

	// Page state bitmap (16 bytes)
	var bitmap [16]byte
	copy(bitmap[:], data[start:start+16])
	start += 16

	return &XdesEntry{
		FileSegmentID: fsegID,
		XdesListNode:  *listNode,
		State:         state,
		PageBitmap:    bitmap,
	}, nil
}
