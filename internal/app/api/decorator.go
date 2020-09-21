package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerDecorator(handler handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := handler(r)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		data, err := json.Marshal(resp.ResponseData)
		if err != nil {
			log.Printf("Error when encode response: %#v\n", err)
			return
		}
		w.Write(data)
	}
}
