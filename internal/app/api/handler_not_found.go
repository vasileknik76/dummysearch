package api

import "net/http"

func (s *Server) notFoundHandler(r *http.Request) response {
	return errorResponseWithText("Not Found", 404)
}
