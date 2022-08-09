package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	r := Router()
	addr := "127.0.0.1:8000"
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %s...", addr)
	log.Fatal(srv.ListenAndServe())
}
