package main

import (
	"context"
	"log"

	"github.com/vasileknik76/dummysearch/internal/app/api"
)

func main() {
	srv := api.NewServer(context.Background(), api.DefaultConfig)
	log.Fatal(srv.ListenAndServe())
}
