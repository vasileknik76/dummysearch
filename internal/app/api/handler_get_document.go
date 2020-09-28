package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) getDocumentHandler(r *http.Request) response {
	vars := mux.Vars(r)
	name := vars["index"]
	id := vars["id"]
	if id == "" {
		return errorResponseWithText("id must be not empty", 400)
	}
	if _, ok := s.indexes[name]; !ok {
		return errorResponseWithText("index not exist", 404)
	}
	i := s.indexes[name].i
	doc := i.GetDocument(id)
	if doc == nil {
		return errorResponseWithText("Document not found", 404)
	}
	return successResponse(responseData{
		Status: true,
		Payload: struct {
			Doc interface{}
		}{
			Doc: struct {
				Meta interface{}
			}{doc.Meta},
		},
	})
}
