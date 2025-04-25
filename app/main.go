package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Xenoneqq/swift-code-database/endpoints"
	"github.com/Xenoneqq/swift-code-database/importer"
	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/Xenoneqq/swift-code-database/storage"
	"github.com/gofiber/fiber/v2"
)

func main() {

	if os.Getenv("MODE") != "" {
		fmt.Printf("Launching application in mode : %s", os.Getenv("MODE"))
	}

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

	// Importing Data from csv (DEFAULT BUILD) (STARTUP ONLY)
	flagPath := "/app/data/import_done.flag"
	if os.Getenv("MODE") == "" {
		if _, err := os.Stat(flagPath); os.IsNotExist(err) {
			fmt.Println("importing data from csv...")
			err := importer.LoadCSV("./data/bank_data.csv", db)

			if err != nil {
				fmt.Println("import failed!")
			} else {
				fmt.Println("done importing!")

				f, err := os.Create(flagPath)
				if err != nil {
					fmt.Println("Failed to create a flag file (data will be imported again on next startup)", err)
				}
				defer f.Close()
			}
		} else {
			fmt.Println("the data has already been imported before")
		}
	} else if os.Getenv("MODE") == "TEST" {
		fmt.Println("launching app in test mode...")
		fmt.Println("preparing application for testing...")
		// test code will be called here once reade :D
	} else if os.Getenv("MODE") == "DEBUG" {
		fmt.Println("launching app in debug mode...")
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
