package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vasileknik76/dummysearch/internal/app/indexer"
)

var DefaultConfig = Config{
	Listen:       "0.0.0.0:6745",
	WriteTimeout: 15 * time.Second,
	ReadTimeout:  15 * time.Second,
}

// Config for API server.
type Config struct {
	Listen       string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	*http.Server
	indexes map[string]IndexInfo
	ctx     context.Context
}

type IndexInfo struct {
	i       *indexer.Index
	destroy func()
}

func NewServer(ctx context.Context, config Config) *Server {
	server := &Server{
		indexes: make(map[string]IndexInfo),
		ctx:     ctx,
	}
	r := mux.NewRouter()
	r.HandleFunc("/", handlerDecorator(server.createIndexHandler)).Methods("POST")
	r.HandleFunc("/{index}/", handlerDecorator(server.createDocumentHandler)).Methods("POST")
	r.HandleFunc("/{index}/", handlerDecorator(server.deleteIndexHandler)).Methods("DELETE")
	r.HandleFunc("/{index}/batch", handlerDecorator(server.createDocumentBatchHandler)).Methods("POST")
	r.HandleFunc("/{index}/search", handlerDecorator(server.searchHandler)).Methods("GET")
	r.HandleFunc("/{index}/compare", handlerDecorator(server.compareHandler)).Methods("GET")
	r.HandleFunc("/{index}/update", handlerDecorator(server.updateIndexHandler)).Methods("GET")
	r.HandleFunc("/{index}/{id}", handlerDecorator(server.deleteDocumentHandler)).Methods("DELETE")
	r.HandleFunc("/{index}/{id}", handlerDecorator(server.getDocumentHandler)).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(handlerDecorator(server.notFoundHandler))

	srv := &http.Server{
		Handler:      r,
		Addr:         config.Listen,
		WriteTimeout: config.WriteTimeout,
		ReadTimeout:  config.ReadTimeout,
	}
	server.Server = srv
	return server
}
