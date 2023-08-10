package main

import (
	"log"

	"github.com/aminasadiam/ChatterBox/datalayer"
)

func main() {
	db, err := datalayer.CreateDbConnection("")
	if err != nil {
		log.Fatalln(err)
		return
	}
}
