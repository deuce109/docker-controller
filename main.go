package main

import (
	"log"
	"net/http"
	"time"

	"github.com/deuce109/docker-controller/v2/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r = handlers.SetContainerRoutes(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
