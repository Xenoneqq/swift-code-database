package testutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

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

type ResError struct {
	Message string `json:"message"`
}

func CheckMessageError(t *testing.T, res *http.Response) error {
	if res == nil {
		t.Errorf("failed to get error from response : res is nil")
		return nil
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("failed to get error from response : error on reading body")
		return nil
	}

	var resError ResError
	err = json.Unmarshal(bodyBytes, &resError)
	if err != nil {
		t.Errorf("failed to get error from response : could not parse json via unmarshal")
		return nil
	}

	return errors.New(resError.Message)
}

func PrintMessageError(t *testing.T, res *http.Response) {
	err := CheckMessageError(t, res)
	if err != nil {
		t.Errorf("Error message : %s", err.Error())
	}
}
