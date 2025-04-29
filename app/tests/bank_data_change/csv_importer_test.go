package main_test

import (
	"testing"

	"github.com/Xenoneqq/swift-code-database/importer"
	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/Xenoneqq/swift-code-database/testutils"
)

func TestDataImportViaCSV(t *testing.T) {
	baseURL := "http://app:8080/api/v1/swift-codes"

	// remove any banks that could break the tests
	defer testutils.DeleteBankBySwiftCodeSafe("IIBKPLKRXXX", baseURL)
	defer testutils.DeleteBankBySwiftCodeSafe("IIBKPLKRAAB", baseURL)
	defer testutils.DeleteBankBySwiftCodeSafe("LCKBPLKRXXX", baseURL)
	defer testutils.DeleteBankBySwiftCodeSafe("TOBKPLKRXXX", baseURL)
	defer testutils.DeleteBankBySwiftCodeSafe("LCKBPLKRABC", baseURL)

	db, err := testutils.SetupTestDatabase()
	if err != nil {
		t.Fatalf("Failed to connect to the database for csv testing")
		return
	}

	err = importer.LoadCSV("./simple_data.csv", db)
	if err != nil {
		t.Fatalf("Failed to import banks to the database\nerror: %v", err)
		return
	}

	var swiftcodes = [5]string{"IIBKPLKRXXX", "IIBKPLKRAAB", "LCKBPLKRXXX", "TOBKPLKRXXX", "LCKBPLKRABC"}

	for _, swiftcode := range swiftcodes {
		var bank []models.Bank
		db.Where("swift_code = ?", swiftcode).Find(&bank)
		if len(bank) == 0 {
			t.Errorf("Bank with swift code %s is missing", swiftcode)
		}
	}

	defer testutils.DeleteBankBySwiftCodeSafe("IIBKPLKRXXX", baseURL)
	defer testutils.DeleteBankBySwiftCodeSafe("IIBKPLKRAAB", baseURL)
	defer testutils.DeleteBankBySwiftCodeSafe("LCKBPLKRXXX", baseURL)
	defer testutils.DeleteBankBySwiftCodeSafe("TOBKPLKRXXX", baseURL)
	defer testutils.DeleteBankBySwiftCodeSafe("LCKBPLKRABC", baseURL)
}
