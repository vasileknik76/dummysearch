package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) createDocumentHandler(r *http.Request) response {
	vars := mux.Vars(r)
	name := vars["index"]
	if _, ok := s.indexes[name]; !ok {
		return errorResponseWithText("index not exist", 404)
	}
	i := s.indexes[name].i
	var request struct {
		ID      string      `json:"id"`
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
	id := i.AddDocument(request.ID, request.Content, request.Meta)
	return successResponse(responseData{
		Status: true,
		Payload: struct {
			Message    string `json:"message"`
			DocumentId string `json:"documentId"`
		}{"OK", id},
	})
}
