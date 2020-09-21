package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/vasileknik76/dummysearch/internal/app/indexer"
	"github.com/vasileknik76/dummysearch/internal/app/text"
)

func (s *Server) createIndexHandler(r *http.Request) response {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errorResponseWithText("Can't read body", 400)
	}
	type _Cfg struct {
		Language text.Language `json:"language"`
	}
	var request struct {
		Name   string `json:"name"`
		Config _Cfg   `json:"config"`
	}
	err = json.Unmarshal(data, &request)
	if err != nil {
		return errorResponseWithText("Can't parse json", 400)
	}
	if _, ok := s.indexes[request.Name]; ok {
		return errorResponseWithText("Index already exists", 400)
	}
	s.indexes[request.Name] = indexer.NewIndex(
		&indexer.IndexConfig{
			Language: request.Config.Language,
		},
	)
	return successResponse(responseData{Status: true, Payload: struct{ Message string }{"OK"}})
}
