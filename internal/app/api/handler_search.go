package api

import (
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/vasileknik76/dummysearch/internal/app/search"
)

type doc struct {
	DocId int
	Meta  interface{}
	Score float64
}

type doclist []doc

func (p doclist) Len() int           { return len(p) }
func (p doclist) Less(i, j int) bool { return p[i].Score < p[j].Score }
func (p doclist) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (srv *Server) searchHandler(r *http.Request) response {
	vars := mux.Vars(r)
	name := vars["index"]
	if _, ok := srv.indexes[name]; !ok {
		return errorResponseWithText("index not exist", 404)
	}
	i := srv.indexes[name]
	query := r.URL.Query().Get("query")
	if query == "" {
		return errorResponseWithText("Query is empty", 400)
	}

	s := search.NewSearcher(i)

	res, err := s.Search(query)
	if err != nil {
		return errorResponseWithText(err.Error(), 500)
	}

	var payload doclist = make([]doc, 0)

	for docID, score := range res {
		d := i.GetDocument(docID)
		payload = append(payload, doc{
			docID,
			d.Meta,
			score,
		})
	}

	sort.Sort(sort.Reverse(payload))

	return successResponse(responseData{
		Status:  true,
		Payload: payload,
	})
}
