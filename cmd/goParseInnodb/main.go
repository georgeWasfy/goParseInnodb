package main

import (
	"fmt"
	"log"

	"goParseInnodb/pkg/innodb/space"
	"goParseInnodb/pkg/innodb/page"
	"goParseInnodb/pkg/innodb/fsphdr"
)

func main() {
	sp, err := space.OpenSpace("/mysql-data/test/t.ibd")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total Pages:", sp.Pages)

	for item := range sp.IteratePages() {
		if item.Err != nil {
			log.Fatal(item.Err)
		}

		switch p := item.Page.(type) {

		case *fsphdr.FspHdrXdesPage:
			fmt.Printf(
				"Page %d: FSP/XDES, SpaceID=%d, XDES=%d\n",
				item.Number,
				p.FspHeader.SpaceID,
				len(p.Xdes),
			)

		case *page.Page:
			// generic page
			// fmt.Printf("Page %d: type=%v\n", item.Number, p.FilHeader.PageType)

		default:
			// should never happen
			fmt.Printf("Page %d: unknown type %T\n", item.Number, p)
		}
	}
}
