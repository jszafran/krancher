package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"krancher/survey"
	"net/http"
)

type WorkloadRequest struct {
	Algorithm        string `json:"algorithm"`
	DataPath         string `json:"data_path"`
	OutputPath       string `json:"output_path"`
	SchemaPath       string `json:"schema_path"`
	OrgStructurePath string `json:"org_structure_path"`
	WorkloadPath     string `json:"workload_path"`
}

func (w *WorkloadRequest) IsValid() bool {
	if w.Algorithm == "" {
		return false
	}

	_, err := survey.ProcessingAlgorithmFactory(w.Algorithm)
	if err != nil {
		return false
	}

	if w.DataPath == "" {
		return false
	}

	if w.OutputPath == "" {
		return false
	}

	if w.SchemaPath == "" {
		return false
	}

	if w.OrgStructurePath == "" {
		return false
	}

	if w.WorkloadPath == "" {
		return false
	}

	return true
}

type WorkloadResponse struct {
	Message string `json:"message"`
}

func CreateWorkloadHandler(w http.ResponseWriter, r *http.Request) {
	var request WorkloadRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(&request)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if !request.IsValid() {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "processing request")
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
