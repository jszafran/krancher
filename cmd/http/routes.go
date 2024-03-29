package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateWorkloadRequest struct {
	DataPath         string `json:"data_path" binding:"required"`
	SchemaPath       string `json:"schema_path" binding:"required"`
	OrgStructurePath string `json:"org_structure_path" binding:"required"`
	WorkloadPath     string `json:"workload_path" binding:"required"`
	IndexBuilder     string `json:"index_builder" binding:"required"`
	Processor        string `json:"processor" binding:"required"`
	OutputPath       string `json:"output_path" binding:"required"`
	Persistor        string `json:"persistor" binding:"required"`
}

func WorkloadCreate(c *gin.Context) {
	var request CreateWorkloadRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "workload created"})
}

func WorkloadList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "listing workloads"})
}
