package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.GET("/", WorkloadList)
	r.POST("/", WorkloadCreate)
	log.Fatal(r.Run(":8000"))
}
