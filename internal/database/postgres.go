package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgreSQL connection and configuration
type PostgreSQL struct {
	DB *gorm.DB
}

var (
	postgresInstance *PostgreSQL
)

// InitPostgreSQL initializes PostgreSQL connection
func InitPostgreSQL() (*PostgreSQL, error) {
	if postgresInstance != nil {
		return postgresInstance, nil
	}

	// Get database connection string from environment
	host := os.Getenv("NEBULABOX_POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("NEBULABOX_POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("NEBULABOX_POSTGRES_USER")
	if user == "" {
		user = "nebulabox"
	}

	password := os.Getenv("NEBULABOX_POSTGRES_PASSWORD")
	if password == "" {
		password = "nebulabox"
	}

	dbname := os.Getenv("NEBULABOX_POSTGRES_DB")
	if dbname == "" {
		dbname = "nebulabox"
	}

	sslmode := os.Getenv("NEBULABOX_POSTGRES_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Configure GORM logger
	logLevel := logger.Info
	if os.Getenv("NEBULABOX_LOG_LEVEL") == "error" {
		logLevel = logger.Error
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	postgresInstance = &PostgreSQL{DB: db}
	log.Printf("✅ Connected to PostgreSQL database: %s@%s:%s/%s", user, host, port, dbname)

	return postgresInstance, nil
}

// GetPostgreSQL returns the PostgreSQL instance
func GetPostgreSQL() *PostgreSQL {
	return postgresInstance
}

// Close closes the PostgreSQL connection
func (p *PostgreSQL) Close() error {
	if p.DB != nil {
		sqlDB, err := p.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate runs database migrations for all models
func (p *PostgreSQL) AutoMigrate() error {
	if p.DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	log.Println("🔄 Running database migrations...")

	// Auto-migrate all models
	err := p.DB.AutoMigrate(
		&Container{},
		&Image{},
		&Workspace{},
		&WorkspaceMember{},
		&Invite{},
		&Session{},
		&Snapshot{},
		&Deployment{},
		&Node{},
		&ContainerGroup{},
		&Template{},
		&User{},
		&Team{},
		&Tenant{},
		&Network{},
		&Service{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// HealthCheck checks if the database connection is healthy
func (p *PostgreSQL) HealthCheck() error {
	if p.DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

