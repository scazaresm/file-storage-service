package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/scazaresm/go-file-storage-service/handler"
)

func main() {
	router := gin.Default()
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}
	router.POST("v1/file", handler.UploadFile)
	router.Run(fmt.Sprintf(":%s", port))
}
