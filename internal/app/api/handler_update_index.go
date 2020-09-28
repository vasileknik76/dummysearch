package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) updateIndexHandler(r *http.Request) response {
	vars := mux.Vars(r)
	name := vars["index"]
	if _, ok := s.indexes[name]; !ok {
		return errorResponseWithText("Index not exist", 404)
	}
	i := s.indexes[name]
	go i.UpdateTFIDF()
	return successResponse(responseData{Status: true, Payload: struct{ Message string }{"Index updating"}})
}
