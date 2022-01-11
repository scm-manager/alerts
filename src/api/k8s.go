package api

import (
	"log"
	"net/http"
)

func OkHandlerFunc(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Println("Failed to write response body")
	}
}

func CreateOkEndpoint() http.Handler {
	return http.HandlerFunc(OkHandlerFunc)
}
