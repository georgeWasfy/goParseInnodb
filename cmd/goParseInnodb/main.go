package main

import (
	"fmt"
	"goParseInnodb/pkg/innodb"
	"log"
)

func main() {
	space, err := innodb.OpenSpace("/mysql-data/test/t.ibd")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Pages:", space.Pages)

}