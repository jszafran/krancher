package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.GET("/", WorkloadList)
	log.Fatal(r.Run(":8000"))
}
