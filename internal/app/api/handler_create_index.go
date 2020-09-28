package api

import (
	"context"
	"encoding/json"
	"fmt"
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
		Language     text.Language `json:"language"`
		CustomIDs    bool          `json:"customIds"`
		UpdatePeriod Duration      `json:"updatePeriod"`
		AutoUpdate   bool          `json:"autoUpdate"`
	}
	var request struct {
		Name   string `json:"name"`
		Config _Cfg   `json:"config"`
	}
	err = json.Unmarshal(data, &request)
	if err != nil {
		return errorResponseWithText(fmt.Sprintf("Can't parse json: %s", err.Error()), 400)
	}

	if request.Config.AutoUpdate && request.Config.UpdatePeriod.Seconds() == 0 {
		return errorResponseWithText("Duration must be positive", 400)
	}

	if _, ok := s.indexes[request.Name]; ok {
		return errorResponseWithText("Index already exists", 400)
	}
	ctx, cancel := context.WithCancel(s.ctx)
	s.indexes[request.Name] = IndexInfo{
		indexer.NewIndex(
			ctx,
			&indexer.IndexConfig{
				Language:     request.Config.Language,
				CustomIDs:    request.Config.CustomIDs,
				AutoUpdate:   request.Config.AutoUpdate,
				UpdatePeriod: request.Config.UpdatePeriod.Duration,
			},
		),
		cancel,
	}
	return successResponse(responseData{Status: true, Payload: struct {
		Message string `json:"message"`
	}{"OK"}})
}
