package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	port := os.Getenv("SERVER_PORT")

	if port == "" {
		port = "8081"
	}

	addr := fmt.Sprintf(":%s", port)

	router.Run(addr)
}
