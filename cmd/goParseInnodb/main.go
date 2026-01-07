package main

import (
	"fmt"
	"goParseInnodb/pkg/innodb/page"
	"log"
)

func main() {
	space, err := innodb.OpenSpace("/mysql-data/test/t.ibd")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total Pages:", space.Pages)

	for page := range space.IteratePages() {
		if page.Err != nil {
			panic(page.Err)
		}
		fmt.Println(page.Page.FilTrailer)
	}

}