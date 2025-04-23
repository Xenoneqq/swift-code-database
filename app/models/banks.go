package models

import "gorm.io/gorm"

type Bank struct {
	Adress        string `json:"adress"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `gorm:"primaryKey" json:"swiftCode"`
}

func MigrateBanks(db *gorm.DB) error {
	err := db.AutoMigrate(&Bank{})
	return err
}
