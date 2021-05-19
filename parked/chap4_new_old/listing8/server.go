package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func apiHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func setupHandlers(r *mux.Router) {
	r.HandleFunc("/healthcheck", healthCheckHandler)
	r.HandleFunc("/api", apiHandler)
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	r := mux.NewRouter()
	setupHandlers(r)

	log.Printf("Server attempting to listen on: %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, r)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}