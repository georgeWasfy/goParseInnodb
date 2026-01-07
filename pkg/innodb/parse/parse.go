package parse

import (
	"goParseInnodb/pkg/innodb/page"
	"goParseInnodb/pkg/innodb/fsphdr"
)

func ParsePage(data []byte) (interface{}, error) {
	p, err := page.NewPage(data)
	if err != nil {
		return nil, err
	}

	switch p.FilHeader.PageType {
		case page.PageTypeFspHdr, page.PageTypeXdes:
			fspPage, err := fsphdr.NewFspHdrXdesPage(p)
			if err != nil {
				return nil, err
			}
			return fspPage, nil
		default:
			return p, nil
    }
}
