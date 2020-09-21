package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) createDocumentBatchHandler(r *http.Request) response {
	vars := mux.Vars(r)
	name := vars["index"]
	if _, ok := s.indexes[name]; !ok {
		return errorResponseWithText("index not exist", 404)
	}
	i := s.indexes[name]
	var request []struct {
		Content string      `json:"content"`
		Meta    interface{} `json:"meta"`
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errorResponseWithText("Can't read body", 400)
	}
	err = json.Unmarshal(data, &request)
	if err != nil {
		return errorResponseWithText("Can't parse json", 400)
	}
	var ids []int
	for _, req := range request {
		id := i.AddDocument(req.Content, req.Meta)
		ids = append(ids, id)
	}

	return successResponse(responseData{
		Status: true,
		Payload: struct {
			Message     string
			DocumentIds []int
		}{"OK", ids},
	})
}
