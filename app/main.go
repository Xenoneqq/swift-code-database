package main

import (
	"log"
	"os"

	"net/http"

	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/Xenoneqq/swift-code-database/storage"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Bank struct {
	Adress        string `json:"adress"`
	BankName      string `json:"bankName`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBank(context *fiber.Ctx) error {
	bank := Bank{}
	err := context.BodyParser(&bank)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "request failed"})
		return err
	}

	res := r.DB.Create(&bank)
	if res.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "failed to create bank entry"})
		return res.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "created bank entry"})
	return nil
}

func (r *Repository) GetAllBanks(context *fiber.Ctx) error {
	banks := &[]models.Banks{}

	res := r.DB.Find(banks)
	if res.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Failed to fetch banks",
		})
		return res.Error
	}

	if res.RowsAffected == 0 {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "There are no banks in the database",
			"data":    []models.Banks{},
		})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Fetched banks successfully",
		"data":    banks,
	})
	return nil
}

func (r *Repository) GetBankByID(context *fiber.Ctx) error {
	banks := &[]models.Banks{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	res := r.DB.Where("swift_code = ?", id).Find(banks)
	if res.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "failed to fetch for bank with SWIFT code",
		})
		return res.Error
	}

	if res.RowsAffected == 0 {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "bank with this SWIFT code does not exist...",
			"data":    []models.Banks{},
		})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "found bank with SWIFT code",
		"data":    banks,
	})
	return nil
}

func (r *Repository) DeleteBank(context *fiber.Ctx) error {
	banks := &[]models.Banks{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	res := r.DB.Where("swift_code = ?", id).Delete(banks)

	if res.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "failed to delete bank",
		})
		return res.Error
	}

	if res.RowsAffected == 0 {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "bank with selected SWIFT id does not exist",
		})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "bank was deleted successfully",
	})
	return nil
}

func (r *Repository) GetBankByCountry(context *fiber.Ctx) error {
	country := context.Params("country")
	banks := &models.Banks{}
	if country == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "country cannot be empty",
		})
		return nil
	}

	res := r.DB.Where("country_iso2 = ?", country).Find(banks)
	if res.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "failed to fetch all SWIFT codes with countryISO2",
		})
		return res.Error
	}

	if res.RowsAffected == 0 {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "No banks with selected country code",
		})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "banks fetched successfully",
		"data":    banks,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/v1/swift-codes", r.CreateBank)
	api.Delete("/v1/swift-codes/:id", r.DeleteBank)

	api.Get("/v1/swift-codes/country/:country", r.GetBankByCountry)
	api.Get("/v1/swift-codes/:id", r.GetBankByID)
	api.Get("/v1/swift-codes", r.GetAllBanks)
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbName := os.Getenv("DB_NAME")

	config := &storage.Config{
		Host:     dbHost,
		Port:     dbPort,
		Password: dbPassword,
		User:     dbUser,
		SSLMode:  dbSSLMode,
		DBNAME:   dbName,
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("failed to load the database")
	}

	err = models.MigrateBanks(db)
	if err != nil {
		log.Fatal("failed to migrate the database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
