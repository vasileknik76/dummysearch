package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vasileknik76/dummysearch/internal/app/search"
)

func (srv *Server) compareHandler(r *http.Request) response {
	vars := mux.Vars(r)
	name := vars["index"]
	if _, ok := srv.indexes[name]; !ok {
		return errorResponseWithText("index not exist", 404)
	}
	i := srv.indexes[name]
	doc1S := r.URL.Query().Get("doc1")
	if doc1S == "" {
		return errorResponseWithText("doc1 is empty", 400)
	}
	doc1, err := strconv.Atoi(doc1S)
	if err != nil {
		return errorResponseWithText("doc1 must be positive number", 400)
	}

	doc2S := r.URL.Query().Get("doc2")
	if doc2S == "" {
		return errorResponseWithText("doc2 is empty", 400)
	}
	doc2, err := strconv.Atoi(doc2S)
	if err != nil {
		return errorResponseWithText("doc2 must be positive number", 400)
	}

	s := search.NewSearcher(i)

	for _, id := range []int{doc1, doc2} {
		if !i.HasDoc(id) {
			return errorResponseWithText(fmt.Sprintf("Document with id #%d not found", id), 404)
		}
	}

	score := s.Score(doc1, doc2)
	if err != nil {
		return errorResponseWithText(err.Error(), 500)
	}

	var payload = struct {
		Score float64 `json:"score"`
	}{score}

	return successResponse(responseData{
		Status:  true,
		Payload: payload,
	})
}
