package main

import (
	"log"

	"github.com/vasileknik76/dummysearch/internal/app/api"
)

func main() {
	srv := api.NewServer(api.DefaultConfig)
	log.Fatal(srv.ListenAndServe())
}
