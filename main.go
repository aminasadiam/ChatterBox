package main

import (
	"log"

	"github.com/aminasadiam/ChatterBox/datalayer"
	"github.com/aminasadiam/ChatterBox/restapi"
)

func main() {
	db, err := datalayer.CreateDbConnection("")
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = restapi.RunApi("localhost:80", *db)
	if err != nil {
		log.Println(err)
		return
	}
}
