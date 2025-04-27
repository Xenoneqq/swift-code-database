package testutils

import (
	"fmt"
	"os"

	"github.com/Xenoneqq/swift-code-database/storage"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func SetupTestDatabase() (*gorm.DB, error) {
	if TestDB != nil {
		return TestDB, nil
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBNAME:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}
	TestDB = db
	return TestDB, nil
}

func CleanupTestDatabase(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
