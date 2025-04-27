package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize gin router with default middleware
	router := gin.Default()

	// define a simple router for the root endpoint
	router.GET("/", func(c *gin.Context) {
		// return a json response with status code 200
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Go E-Commerce API",
		})
	})

	// start the server on port 8080
	log.Println("Server is running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
