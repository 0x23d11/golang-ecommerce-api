package database

import (
	"context"
	"log"
	"time"

	"github.com/0x23d11/go-ecommerce-project/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConnections holds the databse connections
type DBConnections struct {
	Postgres *gorm.DB
	Mongo    *mongo.Client
}

// InitDatabaseConnections initializes the database connections to Postgres and MongoDB
func InitDatabaseConnections(cfg *config.Config) (*DBConnections, error) {
	// Initialize Postgres connection
	log.Println("Connecting to Postgres...")
	pgDB, err := gorm.Open(postgres.Open(cfg.PostgresURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log SQL queries
	})
	if err != nil {
		log.Printf("Failed to connect to Postgres: %v", err)
		return nil, err
	}

	// Configure connection pool for PostgreSQL
	sqlDB, err := pgDB.DB()
	if err != nil {
		log.Printf("Failed to get underlying sql.DB: %v", err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to Postgres")

	// Initialize MongoDB connection
	// --- Initialize MongoDB ---
	log.Println("Connecting to MongoDB...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Ping MongoDB to verify connection
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		// Optionally disconnect if ping fails
		// _ = mongoClient.Disconnect(context.Background())
		return nil, err
	}

	log.Println("Successfully connected to MongoDB.")

	return &DBConnections{
		Postgres: pgDB,
		Mongo:    mongoClient,
	}, nil
}

// CloseDatabaseConnections closes the database connections gracefully
func CloseDatabaseConnections(connections *DBConnections) {
	log.Println("Closing database connections...")

	// Close PostgreSQL
	if connections.Postgres != nil {
		sqlDB, err := connections.Postgres.DB()
		if err == nil {
			err = sqlDB.Close()
			if err != nil {
				log.Printf("Error closing PostgreSQL connection: %v", err)
			} else {
				log.Println("PostgreSQL connection closed.")
			}
		}
	}

	// Close MongoDB
	if connections.Mongo != nil {
		err := connections.Mongo.Disconnect(context.Background())
		if err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		} else {
			log.Println("MongoDB connection closed.")
		}
	}
}
