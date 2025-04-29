package main_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/Xenoneqq/swift-code-database/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGettingBankWithID_branch(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "DDDDPLGHAAB",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bankowa 1",
		IsHeadquarter: false,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankSafe(bank, baseURL)

	resPost := testutils.PostBank(assert, bank, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank: PostBank returned nil response")
		return
	}

	assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank")
	if resPost.StatusCode != http.StatusCreated {
		t.FailNow()
		return
	}

	checkURL := fmt.Sprintf("%s/%s", baseURL, bank.SwiftCode)
	resGet, err := http.Get(checkURL)
	if err != nil {
		t.Fatalf("HTTP GET request to %s failed: %v", checkURL, err)
		return
	}
	if resGet == nil {
		t.Fatalf("HTTP GET request to %s returned nil response", checkURL)
		return
	}

	assert.Equal(http.StatusOK, resGet.StatusCode, "Expected status code 200 for getting bank by SWIFT code %s", bank.SwiftCode)
	if resGet.StatusCode != http.StatusOK {
		return
	}
	defer resGet.Body.Close()

	bankRes, err := testutils.ReadBankFromResponse(assert, resGet)
	if assert.NoError(err, "Failed to read bank from response") {
		assert.Equal(bank, bankRes, "Expected the bank data from DB to match local")
	}

	defer testutils.DeleteBank(assert, bank, baseURL)
}

func TestGettingBankWithID_headquarter(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bankHead := models.Bank{
		SwiftCode:     "CTBNPLGHXXX",
		BankName:      "Central Test Bank",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bank Street 1",
		IsHeadquarter: true,
	}

	bankBranchOne := models.Bank{
		SwiftCode:     "CTBNPLGHAAB",
		BankName:      "Central Test Bank A",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Inwestment Valley 7",
		IsHeadquarter: false,
	}

	bankBranchTwo := models.Bank{
		SwiftCode:     "CTBNPLGHGGJ",
		BankName:      "Central Test Bank B",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Street of Profit 15",
		IsHeadquarter: false,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankSafe(bankHead, baseURL)
	defer testutils.DeleteBankSafe(bankBranchOne, baseURL)
	defer testutils.DeleteBankSafe(bankBranchTwo, baseURL)

	resPost := testutils.PostBank(assert, bankHead, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank: PostBank returned nil response (headquarter)")
		return
	}

	assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank (headquarter)")
	if resPost.StatusCode != http.StatusCreated {
		return
	}

	resPost = testutils.PostBank(assert, bankBranchOne, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank: PostBank returned nil response (branch one)")
		return
	}

	assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank (branch one)")
	if resPost.StatusCode != http.StatusCreated {
		return
	}

	resPost = testutils.PostBank(assert, bankBranchTwo, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank: PostBank returned nil response (branch two)")
		return
	}

	assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank (branch two)")
	if resPost.StatusCode != http.StatusCreated {
		return
	}

	checkURL := fmt.Sprintf("%s/%s", baseURL, bankHead.SwiftCode)
	resGet, err := http.Get(checkURL)
	if err != nil {
		t.Fatalf("HTTP GET request to %s failed: %v", checkURL, err)
		return
	}
	if resGet == nil {
		t.Fatalf("HTTP GET request to %s returned nil response", checkURL)
		return
	}

	assert.Equal(http.StatusOK, resGet.StatusCode, "Expected status code 200 for getting bank by SWIFT code %s", bankHead.SwiftCode)
	if resGet.StatusCode != http.StatusOK {
		return
	}
	defer resGet.Body.Close()

	banks := []models.Bank{}
	banks = append(banks, bankBranchOne, bankBranchTwo)
	var branches []models.Branch = models.ConvertBanksToBranches(banks)

	headquarter := models.Headquarter{}
	headquarter.Address = bankHead.Address
	headquarter.BankName = bankHead.BankName
	headquarter.CountryISO2 = bankHead.CountryISO2
	headquarter.CountryName = bankHead.CountryName
	headquarter.IsHeadquarter = bankHead.IsHeadquarter
	headquarter.SwiftCode = bankHead.SwiftCode
	headquarter.Branches = branches

	bankRes, err := testutils.ReadHeadquarterFromResponse(assert, resGet)
	if assert.NoError(err, "Failed to read bank from response") {
		assert.Equal(headquarter, bankRes, "Expected the bank data from DB to match local")
	}

	defer testutils.DeleteBank(assert, bankHead, baseURL)
	defer testutils.DeleteBank(assert, bankBranchOne, baseURL)
	defer testutils.DeleteBank(assert, bankBranchTwo, baseURL)
}
