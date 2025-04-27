package main_test

import (
	"testing"

	"github.com/Xenoneqq/swift-code-database/handler"
	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/stretchr/testify/assert"
)

func TestBankCheckFunction_Positive(t *testing.T) {
	assert := assert.New(t)

	bankPL := models.Bank{
		SwiftCode:     "IKOPPLPWXXX",
		BankName:      "InvestKomerc Polski Bank",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "ul. Nowy Świat Finansów 5",
		IsHeadquarter: true,
	}
	assert.Nil(handler.CheckBankData(bankPL), "Valid bank data for Poland should return nil")

	bankBE := models.Bank{
		SwiftCode:     "BNAGBEBBXXX",
		BankName:      "Belgian National Bank",
		CountryISO2:   "BE",
		CountryName:   "BELGIUM",
		Address:       "Grand Place 1",
		IsHeadquarter: true,
	}
	assert.Nil(handler.CheckBankData(bankBE), "Valid bank data for Belgium should return nil")

	bankDE := models.Bank{
		SwiftCode:     "DEUTDEFF420",
		BankName:      "Deutsche Einheit Bank",
		CountryISO2:   "DE",
		CountryName:   "GERMANY",
		Address:       "Berliner Allee 42",
		IsHeadquarter: false,
	}
	assert.Nil(handler.CheckBankData(bankDE), "Valid bank data for Germany (non-headquarter) should return nil")

	bankIT := models.Bank{
		SwiftCode:     "INTCITAA123",
		BankName:      "Banca Nazionale Popolare",
		CountryISO2:   "IT",
		CountryName:   "ITALY",
		Address:       "Via Roma 10",
		IsHeadquarter: false,
	}
	assert.Nil(handler.CheckBankData(bankIT), "Valid bank data for Italy (non-headquarter) should return nil")

	bankFR := models.Bank{
		SwiftCode:     "SOGEFRPPARI",
		BankName:      "Societe Generale Paris",
		CountryISO2:   "FR",
		CountryName:   "FRANCE",
		Address:       "Champs-Élysées 16",
		IsHeadquarter: false,
	}
	assert.Nil(handler.CheckBankData(bankFR), "Valid bank data for France (non-headquarter) should return nil")
}

func TestBankCheckFunction_Negative(t *testing.T) {
	assert := assert.New(t)

	bank := models.Bank{
		SwiftCode:     "ALIOPLKRXXXXXX",
		BankName:      "All In One Bank",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Amazing Street 8",
		IsHeadquarter: true,
	}
	assert.EqualError(handler.CheckBankData(bank), "SWIFT code must be 11 characters long")

	bank = models.Bank{
		SwiftCode:     "123OPLKRXXX",
		BankName:      "All In One Bank",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Amazing Street 8",
		IsHeadquarter: true,
	}
	assert.EqualError(handler.CheckBankData(bank), "first 6 characters of the SWIFT code cannot contain numbers")

	bank = models.Bank{
		SwiftCode:     "aliOPLKRxxx",
		BankName:      "All In One Bank",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Amazing Street 8",
		IsHeadquarter: true,
	}
	assert.EqualError(handler.CheckBankData(bank), "SWIFT code must consist only of UPPER LETTERS")

	bank = models.Bank{
		SwiftCode:     "ALIOPLKRXXX",
		BankName:      "All In One Bank",
		CountryISO2:   "PL",
		CountryName:   "poland",
		Address:       "Amazing Street 8",
		IsHeadquarter: true,
	}
	assert.EqualError(handler.CheckBankData(bank), "country name must be all UPPER CASE")

	bank = models.Bank{
		SwiftCode:     "ALIOPLKRXXX",
		BankName:      "All In One Bank",
		CountryISO2:   "DE",
		CountryName:   "POLAND",
		Address:       "Amazing Street 8",
		IsHeadquarter: true,
	}
	assert.EqualError(handler.CheckBankData(bank), "country ISO2 does not match the one inside the SWIFT code : DE != PL")

	bank = models.Bank{
		SwiftCode:     "ALIOPLKRXXX",
		BankName:      "All In One Bank",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Amazing Street 8",
		IsHeadquarter: false,
	}
	assert.EqualError(handler.CheckBankData(bank), "bank with XXX as the last 3 letters of the SWIFT code has to be a headquarter")

	bank = models.Bank{
		SwiftCode:     "ALIOPLKRAAB",
		BankName:      "All In One Bank",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Amazing Street 8",
		IsHeadquarter: true,
	}
	assert.EqualError(handler.CheckBankData(bank), "bank without XXX as the last 3 letters of the SWIFT code cannot to be a headquarter")
}

func TestBankToBranchConversion(t *testing.T) {
	assert := assert.New(t)

	banks := []models.Bank{
		{
			SwiftCode:     "PLBIPLPWAAA",
			BankName:      "Polski Bank Inwestycyjny",
			CountryISO2:   "PL",
			CountryName:   "POLAND",
			Address:       "ul. Nowy Świat Finansów 5",
			IsHeadquarter: false,
		},
		{
			SwiftCode:     "INGBPLPWXXX",
			BankName:      "ING Bank Śląski",
			CountryISO2:   "PL",
			CountryName:   "POLAND",
			Address:       "ul. Sokolska 34",
			IsHeadquarter: true,
		},
		{
			SwiftCode:     "PKOOPWPPXXX",
			BankName:      "PKO Bank Polski",
			CountryISO2:   "PL",
			CountryName:   "POLAND",
			Address:       "ul. Puławska 15",
			IsHeadquarter: true,
		},
		{
			SwiftCode:     "BREXPLPWMBK",
			BankName:      "mBank S.A.",
			CountryISO2:   "PL",
			CountryName:   "POLAND",
			Address:       "ul. Senatorska 18",
			IsHeadquarter: false,
		},
		{
			// does not check for shorter code should return as normal!
			SwiftCode:     "RBOSPLPW",
			BankName:      "Bank Zachodni WBK",
			CountryISO2:   "PL",
			CountryName:   "POLAND",
			Address:       "pl. Grunwaldzki 10",
			IsHeadquarter: true,
		},
	}

	outcomeBranches := []models.Branch{
		{
			SwiftCode:     "PLBIPLPWAAA",
			BankName:      "Polski Bank Inwestycyjny",
			CountryISO2:   "PL",
			Address:       "ul. Nowy Świat Finansów 5",
			IsHeadquarter: false,
		},
		{
			SwiftCode:     "INGBPLPWXXX",
			BankName:      "ING Bank Śląski",
			CountryISO2:   "PL",
			Address:       "ul. Sokolska 34",
			IsHeadquarter: true,
		},
		{
			SwiftCode:     "PKOOPWPPXXX",
			BankName:      "PKO Bank Polski",
			CountryISO2:   "PL",
			Address:       "ul. Puławska 15",
			IsHeadquarter: true,
		},
		{
			SwiftCode:     "BREXPLPWMBK",
			BankName:      "mBank S.A.",
			CountryISO2:   "PL",
			Address:       "ul. Senatorska 18",
			IsHeadquarter: false,
		},
		{
			SwiftCode:     "RBOSPLPW",
			BankName:      "Bank Zachodni WBK",
			CountryISO2:   "PL",
			Address:       "pl. Grunwaldzki 10",
			IsHeadquarter: true,
		},
	}

	branches := models.ConvertBanksToBranches(banks)

	assert.Equal(len(outcomeBranches), len(branches), "Number of branches should match the number of banks")

	for i := 0; i < len(banks); i++ {
		assert.Equal(outcomeBranches[i], branches[i], "Branch at index %d should match", i)
	}
}
