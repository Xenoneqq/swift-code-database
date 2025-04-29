package main_test

import (
	"encoding/json"
	"fmt"
	"io"
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

	var resPost *http.Response = testutils.PostBank(assert, bank, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank: PostBank returned nil response")
		return
	}

	var passed bool = assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank")
	if !passed {
		testutils.PrintMessageError(t, resPost)
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

	passed = assert.Equal(http.StatusOK, resGet.StatusCode, "Expected status code 200 for getting bank by SWIFT code %s", bank.SwiftCode)
	if !passed {
		testutils.PrintMessageError(t, resGet)
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

	var resPost *http.Response = testutils.PostBank(assert, bankHead, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank: PostBank returned nil response (headquarter)")
		return
	}

	var passed bool = assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank (headquarter)")
	if !passed {
		testutils.PrintMessageError(t, resPost)
		return
	}

	resPost = testutils.PostBank(assert, bankBranchOne, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank: PostBank returned nil response (branch one)")
		return
	}

	passed = assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank (branch one)")
	if !passed {
		testutils.PrintMessageError(t, resPost)
		return
	}

	resPost = testutils.PostBank(assert, bankBranchTwo, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank: PostBank returned nil response (branch two)")
		return
	}

	passed = assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank (branch two)")
	if !passed {
		testutils.PrintMessageError(t, resPost)
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

	passed = assert.Equal(http.StatusOK, resGet.StatusCode, "Expected status code 200 for getting bank by SWIFT code %s", bankHead.SwiftCode)
	if !passed {
		testutils.PrintMessageError(t, resGet)
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

type iso2Info struct {
	CountryISO2 string        `json:"countryISO2"`
	CountryName string        `json:"countryName"`
	SwiftCodes  []models.Bank `json:"swiftCodes"`
}

type iso2Response struct {
	Message string   `json:"message"`
	Data    iso2Info `json:"data"`
}

func TestGettingBanksWithISO2(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bankOne := models.Bank{
		SwiftCode:     "BOTSDEGHXXX",
		BankName:      "Bank of Testing",
		CountryISO2:   "DE",
		CountryName:   "GERMANY",
		Address:       "Bank Street 7",
		IsHeadquarter: true,
	}

	bankTwo := models.Bank{
		SwiftCode:     "BOTSDEGHAAA",
		BankName:      "Bank of Testing",
		CountryISO2:   "DE",
		CountryName:   "GERMANY",
		Address:       "Investment Street 25",
		IsHeadquarter: false,
	}

	bankThree := models.Bank{
		SwiftCode:     "UNBSPLGAXXX",
		BankName:      "Universal Banking Solutions",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bright Street 10",
		IsHeadquarter: true,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankSafe(bankOne, baseURL)
	defer testutils.DeleteBankSafe(bankTwo, baseURL)
	defer testutils.DeleteBankSafe(bankThree, baseURL)

	// first bank (GERMANY)
	var resPost *http.Response = testutils.PostBank(assert, bankOne, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank one: PostBank returned nil response")
		return
	}

	var passed bool = assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank one")
	if !passed {
		testutils.PrintMessageError(t, resPost)
		return
	}

	// second bank (GERMANY)
	resPost = testutils.PostBank(assert, bankTwo, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank two: PostBank returned nil response")
		return
	}

	passed = assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank two")
	if !passed {
		testutils.PrintMessageError(t, resPost)
		return
	}

	// third bank (POLAND)
	resPost = testutils.PostBank(assert, bankThree, baseURL)
	if resPost == nil {
		t.Fatalf("Failed to create test bank three: PostBank returned nil response")
		return
	}

	passed = assert.Equal(http.StatusCreated, resPost.StatusCode, "Expected status code 201 for creating test bank three")
	if !passed {
		testutils.PrintMessageError(t, resPost)
		return
	}

	banksByCountry, err := http.Get(fmt.Sprintf("%s/%s/%s", baseURL, "country", "DE"))
	if err != nil {
		t.Fatalf("Failed to fetch HTTP : %v", err)
		return
	}

	passed = assert.Equal(http.StatusOK, banksByCountry.StatusCode, "Expected status code 200 for getting banks with ISO2 = DE")
	if !passed {
		testutils.PrintMessageError(t, resPost)
		return
	}

	bodyBytes, err := io.ReadAll(banksByCountry.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
		return
	}

	var iso2res iso2Response
	err = json.Unmarshal(bodyBytes, &iso2res)
	if err != nil {
		t.Fatalf("Failed to convert json to response (unmarshal fail): %v", err)
		return
	}

	var correctRes iso2Response = iso2Response{}
	var correctInfo iso2Info = iso2Info{}
	correctInfo.CountryISO2 = "DE"
	correctInfo.CountryName = "GERMANY"
	correctInfo.SwiftCodes = []models.Bank{}
	correctInfo.SwiftCodes = append(correctInfo.SwiftCodes, bankOne, bankTwo)
	correctRes.Data = correctInfo
	correctRes.Message = "banks fetched successfully"

	assert.Equal(iso2res, correctRes, "Expected different response")

	defer testutils.DeleteBankSafe(bankOne, baseURL)
	defer testutils.DeleteBankSafe(bankTwo, baseURL)
	defer testutils.DeleteBankSafe(bankThree, baseURL)
}
