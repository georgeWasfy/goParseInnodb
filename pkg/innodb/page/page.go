package innodb

import (
	"fmt"
)

type PageType uint16
type Page struct {
	Data []byte  // full 16 KB
	FilHeader  *FilHeader
	FilTrailer *FilTrailer
}
type PageWrapper struct {
	Number int64
	Page   *Page
	Err    error
}
func NewPage(data []byte) (*Page, error) {
	if len(data) < PAGE_SIZE {
		return nil, fmt.Errorf("page data too short: got %d bytes", len(data))
	}

	filHeader, err := NewFilHeader(data)
	if err != nil {
		return nil, fmt.Errorf("fil header parse failed: %w", err)
	}

	filTrailer, err := NewFilTrailer(data[len(data)-FIL_TRAILER_SIZE:])
	if err != nil {
		return nil, fmt.Errorf("fil trailer parse failed: %w", err)
	}

	return &Page{
		Data:       data,
		FilHeader:  filHeader,
		FilTrailer: filTrailer,
	}, nil
}
const (
	PageTypeAllocated     PageType = 0
	PageTypeUndoLog       PageType = 2
	PageTypeINode         PageType = 3
	PageTypeIBufFreeList  PageType = 4
	PageTypeIBufBitmap   PageType = 5
	PageTypeSys           PageType = 6
	PageTypeTrxSys        PageType = 7
	PageTypeFspHdr        PageType = 8
	PageTypeXdes          PageType = 9
	PageTypeBlob          PageType = 10
	PageTypeZBlob         PageType = 11
	PageTypeZBlob2        PageType = 12
	PageTypeIndex         PageType = 17855
)
func (pt PageType) String() string {
	switch pt {
	case PageTypeAllocated:
		return "ALLOCATED"
	case PageTypeUndoLog:
		return "UNDO_LOG"
	case PageTypeINode:
		return "INODE"
	case PageTypeIBufFreeList:
		return "IBUF_FREE_LIST"
	case PageTypeIBufBitmap:
		return "IBUF_BITMAP"
	case PageTypeSys:
		return "SYS"
	case PageTypeTrxSys:
		return "TRX_SYS"
	case PageTypeFspHdr:
		return "FSP_HDR"
	case PageTypeXdes:
		return "XDES"
	case PageTypeBlob:
		return "BLOB"
	case PageTypeZBlob:
		return "ZBLOB"
	case PageTypeZBlob2:
		return "ZBLOB2"
	case PageTypeIndex:
		return "INDEX"
	default:
		return "UNKNOWN"
	}
}



