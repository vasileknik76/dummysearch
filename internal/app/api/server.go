package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vasileknik76/dummysearch/internal/app/indexer"
)

var DefaultConfig = Config{
	Listen:       "127.0.0.1:6745",
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
	indexes map[string]*indexer.Index
}

func NewServer(config Config) *Server {
	server := &Server{
		indexes: make(map[string]*indexer.Index),
	}
	r := mux.NewRouter()
	r.HandleFunc("/", handlerDecorator(server.createIndexHandler)).Methods("POST")
	r.HandleFunc("/{index}/", handlerDecorator(server.createDocumentHandler)).Methods("POST")
	r.HandleFunc("/{index}/batch", handlerDecorator(server.createDocumentBatchHandler)).Methods("POST")
	r.HandleFunc("/{index}/search", handlerDecorator(server.searchHandler)).Methods("GET")
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
