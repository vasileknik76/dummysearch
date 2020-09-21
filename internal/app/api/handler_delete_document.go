package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) deleteDocumentHandler(r *http.Request) response {
	vars := mux.Vars(r)
	name := vars["index"]
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errorResponseWithText("id must be number", 400)
	}
	if _, ok := s.indexes[name]; !ok {
		return errorResponseWithText("index not exist", 404)
	}
	i := s.indexes[name]
	i.DeleteDocument(id)
	return successResponse(responseData{Status: true, Payload: struct{ Message string }{"OK"}})
}
