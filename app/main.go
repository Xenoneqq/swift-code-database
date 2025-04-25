package main

import (
	"log"
	"os"

	"github.com/Xenoneqq/swift-code-database/endpoints"
	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/Xenoneqq/swift-code-database/storage"
	"github.com/gofiber/fiber/v2"
)

func main() {

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBNAME:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("failed to load the database")
	}

	err = models.MigrateBanks(db)
	if err != nil {
		log.Fatal("failed to migrate the database")
	}

	res := db.Exec(`CREATE INDEX IF NOT EXISTS idx_swiftcode_8 ON banks ((LEFT("swift_code", 8)));`)
	if res.Error != nil {
		log.Fatal("failed to create index for fast searches")
	}

	r := endpoints.Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
