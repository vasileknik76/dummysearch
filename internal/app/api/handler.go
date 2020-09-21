package api

import "net/http"

type handler func(r *http.Request) response
