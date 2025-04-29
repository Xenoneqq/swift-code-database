package models

import "gorm.io/gorm"

type Bank struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName" gorm:"not null"`
	CountryISO2   string `json:"countryISO2" gorm:"not null"`
	CountryName   string `json:"countryName" gorm:"not null"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `gorm:"primaryKey;not null" json:"swiftCode"`
}

type Branch struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `gorm:"primaryKey" json:"swiftCode"`
}

type Headquarter struct {
	Address       string   `json:"address"`
	BankName      string   `json:"bankName"`
	CountryISO2   string   `json:"countryISO2"`
	CountryName   string   `json:"countryName"`
	IsHeadquarter bool     `json:"isHeadquarter"`
	SwiftCode     string   `gorm:"primaryKey" json:"swiftCode"`
	Branches      []Branch `json:"branches"`
}

func ConvertBanksToBranches(banks []Bank) []Branch {
	branches := make([]Branch, len(banks))
	for i, bank := range banks {
		branches[i].Address = bank.Address
		branches[i].BankName = bank.BankName
		branches[i].CountryISO2 = bank.CountryISO2
		branches[i].IsHeadquarter = bank.IsHeadquarter
		branches[i].SwiftCode = bank.SwiftCode
	}
	return branches
}

func MigrateBanks(db *gorm.DB) error {
	err := db.AutoMigrate(&Bank{})
	return err
}
