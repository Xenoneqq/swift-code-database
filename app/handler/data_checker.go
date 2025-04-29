package handler

import (
	"errors"
	"strings"

	"github.com/Xenoneqq/swift-code-database/models"
)

func CheckBankData(bank models.Bank) error {

	if strings.TrimSpace(bank.BankName) == "" {
		return errors.New("bank name cannot be left empty")
	}

	if strings.TrimSpace(bank.CountryISO2) == "" {
		return errors.New("bank iso2 code cannot be left empty")
	}

	if strings.TrimSpace(bank.CountryName) == "" {
		return errors.New("bank country cannot be left empty")
	}

	if strings.TrimSpace(bank.SwiftCode) == "" {
		return errors.New("bank SWIFT code cannot be left empty")
	}

	if len(bank.SwiftCode) != 11 {
		return errors.New("SWIFT code must be 11 characters long")
	}

	if bank.SwiftCode != strings.ToUpper(bank.SwiftCode) {
		return errors.New("SWIFT code must consist only of UPPER LETTERS")
	}

	if strings.ContainsAny(bank.SwiftCode[:6], "0123456789") {
		return errors.New("first 6 characters of the SWIFT code cannot contain numbers")
	}

	if bank.CountryName != strings.ToUpper(bank.CountryName) {
		return errors.New("country name must be all UPPER CASE")
	}

	if bank.CountryISO2 != bank.SwiftCode[4:6] {
		return errors.New("country ISO2 does not match the one inside the SWIFT code : " + bank.CountryISO2 + " != " + bank.SwiftCode[4:6])
	}

	if bank.SwiftCode[8:] == "XXX" && !bank.IsHeadquarter {
		return errors.New("bank with XXX as the last 3 letters of the SWIFT code has to be a headquarter")
	}

	if bank.SwiftCode[8:] != "XXX" && bank.IsHeadquarter {
		return errors.New("bank without XXX as the last 3 letters of the SWIFT code cannot to be a headquarter")
	}

	err := CheckCountry(bank.SwiftCode[4:6], bank.CountryName)
	switch err {
	case 1:
		return errors.New("country with selected ISO2 id does not exist")
	case 2:
		return errors.New("ISO2 does not match the selected country")
	}

	return nil
}
