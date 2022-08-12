package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WorkloadRequest struct {
	Algorithm        string `json:"algorithm" binding:"required"`
	DataPath         string `json:"data_path" binding:"required"`
	OutputPath       string `json:"output_path" binding:"required"`
	SchemaPath       string `json:"schema_path" binding:"required"`
	OrgStructurePath string `json:"org_structure_path" binding:"required"`
	WorkloadPath     string `json:"workload_path" binding:"required"`
}

type WorkloadResponse struct {
	Message string `json:"message"`
}

func WorkloadCreate(c *gin.Context) {
	var request WorkloadRequest
	err := c.BindJSON(&request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	fmt.Printf("%+v\n", request)
	c.JSON(http.StatusCreated, gin.H{"message": "workload created"})
}

func WorkloadList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "listing workloads"})
}
