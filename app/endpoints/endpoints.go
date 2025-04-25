package endpoints

import (
	"net/http"

	"github.com/Xenoneqq/swift-code-database/handler"
	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBank(context *fiber.Ctx) error {
	bank := models.Bank{}

	err := context.BodyParser(&bank)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to parse bank data to database entry (incorrect bank details)"})
		return err
	}

	err = handler.CheckBankData(bank)
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "failed to create bank entry"})
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
	banks := &[]models.Bank{}

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
			"data":    []models.Bank{},
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
	bank := &models.Bank{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	res := r.DB.Where("swift_code = ?", id).First(bank)
	if res.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "failed to fetch for bank with SWIFT code",
		})
		return res.Error
	}

	if res.RowsAffected == 0 {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "bank with this SWIFT code does not exist...",
			"data":    []models.Bank{},
		})
		return nil
	}

	if bank.IsHeadquarter {

		branchesFull := &[]models.Bank{}
		first8 := id[:8]
		res = r.DB.Where("LEFT(swift_code, 8) = ?", first8).Where("is_headquarter != ?", bank.IsHeadquarter).Find(branchesFull)
		if res.Error != nil {
			context.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "failed to fetch for branches of bank with SWIFT code",
			})
			return res.Error
		}

		branches := models.ConvertBanksToBranches(*branchesFull)
		headquarter := models.Headquarter{}

		headquarter.Address = bank.Address
		headquarter.BankName = bank.BankName
		headquarter.CountryISO2 = bank.CountryISO2
		headquarter.CountryName = bank.CountryName
		headquarter.IsHeadquarter = bank.IsHeadquarter
		headquarter.SwiftCode = bank.SwiftCode
		headquarter.Branches = branches

		context.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "found bank with SWIFT code",
			"data":    headquarter,
		})
	} else {

		context.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "found bank with SWIFT code",
			"data": fiber.Map{
				"address":       bank.Address,
				"bankName":      bank.BankName,
				"countryISO2":   bank.CountryISO2,
				"countryName":   bank.CountryName,
				"isHeadquarter": bank.IsHeadquarter,
				"swiftCode":     bank.SwiftCode,
			},
		})
	}
	return nil
}

func (r *Repository) DeleteBank(context *fiber.Ctx) error {
	banks := &[]models.Bank{}
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
	banks := &[]models.Bank{}
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

	countryCode := (*banks)[0].CountryISO2
	countryName := (*banks)[0].CountryName

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "banks fetched successfully",
		"data": fiber.Map{
			"countryISO2": countryCode,
			"countryName": countryName,
			"swiftCodes":  banks,
		},
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
