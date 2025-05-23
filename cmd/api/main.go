package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/0x23d11/go-ecommerce-project/internal/config"
	"github.com/0x23d11/go-ecommerce-project/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// load configuration
	cfg := config.LoadConfig()

	// initialize database connections
	dbConnections, err := database.InitDatabaseConnections(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database connections: %v", err)
	}

	// defer closing database connections when main function exits
	defer database.CloseDatabaseConnections(dbConnections)

	log.Printf("Postgres Connection: %p", dbConnections.Postgres)
	log.Printf("Mongo Connection: %p", dbConnections.Mongo)

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
	serverAddress := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server is running on port %s...", serverAddress)
	if err := router.Run(serverAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
