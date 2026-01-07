package innodb

import (
	"io"
)
const (
	SIZE_LIST_NODE = 12
)
// All pointers point to the beginning (that is, N+0) of the list node, 
// not necessarily of the structure being linked together. For example, 
// when extent descriptor entries are linked in a list, 
// since the list node is at offset 8 within the XDES entry structure, 
// the code reading a list entry must “know” that the descriptor 
// structure starts 8 bytes before the offset of the list node, and read the structure from there. 
type ListNode struct {
	PrevPageNumber uint32
	NextPageNumber uint32

	PrevOffset     uint16
	NextOffset     uint16
}

func NewListNode(data []byte, offset int) (*ListNode, int, error) {
	if offset + SIZE_LIST_NODE > len(data) {
		return nil, offset, io.ErrUnexpectedEOF
	}

	prevPage := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	prevOff := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	nextPage := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	nextOff := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	return &ListNode{
		PrevPageNumber: prevPage,
		PrevOffset:     prevOff,
		NextPageNumber: nextPage,
		NextOffset:     nextOff,
	}, offset, nil
}
