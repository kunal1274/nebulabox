package database

import (
	"fmt"
	"log"
	"os"
)

// Init initializes both PostgreSQL and MongoDB connections
func Init() error {
	// Check if databases should be initialized
	skipPostgres := os.Getenv("NEBULABOX_SKIP_POSTGRES") == "true"
	skipMongo := os.Getenv("NEBULABOX_SKIP_MONGODB") == "true"

	// Initialize PostgreSQL
	if !skipPostgres {
		postgres, err := InitPostgreSQL()
		if err != nil {
			log.Printf("⚠️  Warning: Failed to initialize PostgreSQL: %v", err)
			log.Println("   Continuing without PostgreSQL (using in-memory storage)")
		} else {
			// Run auto-migrations
			if err := postgres.AutoMigrate(); err != nil {
				log.Printf("⚠️  Warning: Failed to run PostgreSQL migrations: %v", err)
			}
			// Note: Repositories should be initialized separately after database initialization
		}
	} else {
		log.Println("ℹ️  PostgreSQL initialization skipped (NEBULABOX_SKIP_POSTGRES=true)")
	}

	// Initialize MongoDB
	if !skipMongo {
		_, err := InitMongoDB()
		if err != nil {
			log.Printf("⚠️  Warning: Failed to initialize MongoDB: %v", err)
			log.Println("   Continuing without MongoDB (logs/metrics won't be persisted)")
		}
	} else {
		log.Println("ℹ️  MongoDB initialization skipped (NEBULABOX_SKIP_MONGODB=true)")
	}

	return nil
}

// Close closes all database connections
func Close() error {
	var errs []error

	if postgresInstance != nil {
		if err := postgresInstance.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if mongoInstance != nil {
		if err := mongoInstance.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing databases: %v", errs)
	}

	return nil
}

