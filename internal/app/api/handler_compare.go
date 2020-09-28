package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vasileknik76/dummysearch/internal/app/search"
)

func (srv *Server) compareHandler(r *http.Request) response {
	vars := mux.Vars(r)
	name := vars["index"]
	if _, ok := srv.indexes[name]; !ok {
		return errorResponseWithText("index not exist", 404)
	}
	i := srv.indexes[name].i
	doc1 := r.URL.Query().Get("doc1")
	if doc1 == "" {
		return errorResponseWithText("doc1 is empty", 400)
	}

	doc2 := r.URL.Query().Get("doc2")
	if doc2 == "" {
		return errorResponseWithText("doc2 is empty", 400)
	}

	s := search.NewSearcher(i)

	for _, id := range []string{doc1, doc2} {
		if !i.HasDoc(id) {
			return errorResponseWithText(fmt.Sprintf("Document with id #%d not found", id), 404)
		}
	}

	score := s.Score(doc1, doc2)

	var payload = struct {
		Score float64 `json:"score"`
	}{score}

	return successResponse(responseData{
		Status:  true,
		Payload: payload,
	})
}
