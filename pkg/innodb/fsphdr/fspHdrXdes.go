package fsphdr

import (
	"fmt"
	"goParseInnodb/pkg/innodb/page"
)

const (
	MAX_EXTENT_DESCRIPTORS = 256
	XDES_ENTRY_SIZE  = 40
	XDES_START_OFFSET = 150
	FSP_HEADER_START_OFFSET = 38
)
type FspHdrXdesPage struct {
	*page.Page
	// FSP_HDR / XDES
	FspHeader *FspHeader // filled with zeros form xes pages
	Xdes      []XdesEntry
}
func NewFspHdrXdesPage(p *page.Page) (*FspHdrXdesPage,error) {
	pt := p.FilHeader.PageType
	if pt != page.PageTypeFspHdr && pt != page.PageTypeXdes {
		return nil, fmt.Errorf("page is not FSP_HDR or XDES")
	}
	// parse FSP header
	fspHeader, err := NewFspHeader(p.Data, FSP_HEADER_START_OFFSET)
	if err != nil {
		return nil, err
	}

	// parse XDES entries
	xdes, err := ParseXdesEntries(p.Data, XDES_START_OFFSET)
	if err != nil {
		return nil, err
	}

	return &FspHdrXdesPage{
		Page:      p,
		FspHeader: fspHeader,
		Xdes:      xdes,
	}, nil
}

func ParseXdesEntries(data []byte, startOffset int) ([]XdesEntry, error) {
	xdes := make([]XdesEntry, 0, MAX_EXTENT_DESCRIPTORS)

	offset := startOffset

	for i := 0; i < MAX_EXTENT_DESCRIPTORS; i++ {
		entry, err := NewXdesEntry(data, offset)
		if err != nil {
			return nil, err
		}

		xdes = append(xdes, *entry)
		offset += XDES_ENTRY_SIZE
	}

	return xdes, nil
}
