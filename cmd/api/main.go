package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/condemo/raspi-home-service/api"
	"github.com/condemo/raspi-home-service/store"
)

func main() {
	addr := flag.String("p", ":4000", "addr")
	flag.Parse()

	sqlStorage := store.NewPostgresStore()
	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStorage(db)

	apiServer := api.NewAPIServer(*addr, store)
	fmt.Println("Starting API server at", *addr)
	apiServer.Run()
}
