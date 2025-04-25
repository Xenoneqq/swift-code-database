package importer

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/Xenoneqq/swift-code-database/handler"
	"github.com/Xenoneqq/swift-code-database/models"
	"gorm.io/gorm"
)

func LoadCSV(path string, db *gorm.DB) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i := 1; i < len(records); i++ {
		bank := models.Bank{
			CountryISO2: records[i][0],
			SwiftCode:   records[i][1],
			BankName:    records[i][3],
			Address:     records[i][4],
			CountryName: records[i][6],
		}
		if bank.SwiftCode[8:] == "XXX" {
			bank.IsHeadquarter = true
		} else {
			bank.IsHeadquarter = false
		}

		err := handler.CheckBankData(bank)
		if err != nil {
			fmt.Printf("(csv) failed to insert bank, wrong details : %s\nerror: %s",
				bank.SwiftCode, err.Error())
			continue
		}

		res := db.Create(&bank)
		if res.Error != nil {
			fmt.Printf("(csv) failed to insert bank, database error : %s\nerror: %s",
				bank.SwiftCode, err.Error())
		}
	}

	defer file.Close()
	return nil
}
