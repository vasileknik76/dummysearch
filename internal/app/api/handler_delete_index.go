package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) deleteIndexHandler(r *http.Request) response {
	vars := mux.Vars(r)
	name := vars["index"]
	if _, ok := s.indexes[name]; !ok {
		return errorResponseWithText("index not exist", 404)
	}
	s.indexes[name].destroy()
	delete(s.indexes, name)
	return successResponse(responseData{Status: true, Payload: struct{ Message string }{"OK"}})
}
