package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type WorkloadRequest struct {
	Algorithm string `json:"algorithm"`
}

type WorkloadResponse struct {
	Message string `json:"message"`
}

func CreateWorkloadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "processing workload")
}

func ListWorkloadsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "listing workloads")
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/workloads", ListWorkloadsHandler).Methods("GET")
	r.HandleFunc("/workloads", CreateWorkloadHandler).Methods("POST")
	return r
}
