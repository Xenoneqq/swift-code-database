package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/Xenoneqq/swift-code-database/testutils"
	"github.com/stretchr/testify/assert"
)

func TestInsertingBank_Positive_Headquarter(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	headquarterBank := models.Bank{
		SwiftCode:     "ABCDPLGHXXX",
		BankName:      "Central Bank of Testing",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bank st 1",
		IsHeadquarter: true,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankSafe(headquarterBank, baseURL)

	var res *http.Response = testutils.PostBank(assert, headquarterBank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusCreated, res.StatusCode, "Expected status code 201 for headquarter creation")
	if http.StatusCreated != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}

	defer testutils.DeleteBankSafe(headquarterBank, baseURL)
}

func TestInsertingBank_Positive_Branch(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	branchBank := models.Bank{
		SwiftCode:     "IJKLPLPQAAB",
		BankName:      "Bank of Testing",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Sidestreet 2",
		IsHeadquarter: false,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankSafe(branchBank, baseURL)

	var res *http.Response = testutils.PostBank(assert, branchBank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusCreated, res.StatusCode, "Expected status code 201 for headquarter creation")
	if http.StatusCreated != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}

	defer testutils.DeleteBankSafe(branchBank, baseURL)
}

type BankDetailed struct {
	Address         string `json:"address"`
	BankName        string `json:"bankName" gorm:"not null"`
	CountryISO2     string `json:"countryISO2" gorm:"not null"`
	CountryName     string `json:"countryName" gorm:"not null"`
	IsHeadquarter   bool   `json:"isHeadquarter"`
	SwiftCode       string `gorm:"primaryKey;not null" json:"swiftCode"`
	BankDescription string `json:"bankDescription"`
}

func TestAdditionalBankData(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bankDetails := BankDetailed{
		SwiftCode:       "IJKLPLPQAAB",
		BankName:        "Tester Bank",
		CountryISO2:     "PL",
		CountryName:     "POLAND",
		Address:         "Boczna 2",
		IsHeadquarter:   false,
		BankDescription: "Bank for testing purposes only",
	}

	bankOutcome := models.Bank{
		SwiftCode:     "IJKLPLPQAAB",
		BankName:      "Tester Bank",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Boczna 2",
		IsHeadquarter: false,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankSafe(bankOutcome, baseURL)

	jsonData, err := json.Marshal(bankDetails)
	if err != nil {
		t.Fatalf("Failed to marshal bank data to JSON: %v", err)
		return
	}

	res, err := http.Post(baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("HTTP POST request for bank failed: %v", err)
		return
	}

	assert.Equal(http.StatusCreated, res.StatusCode, "Expected status code 201 for creating a bank with additional info")
	if http.StatusCreated != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}

	res, err = http.Get(fmt.Sprintf("%s/%s", baseURL, bankDetails.SwiftCode))
	if err != nil {
		t.Fatalf("HTTP GET request for bank failed: %v", err)
		return
	}
	addedBank, err := testutils.ReadBankFromResponse(assert, res)
	if assert.NoError(err, "Failed to read bank from response") {
		assert.Equal(bankOutcome, addedBank, "Expected the bank data from DB to match local")
	}

	defer testutils.DeleteBankSafe(bankOutcome, baseURL)
}

type WrongBank struct {
	Address       string `json:"address"`
	BankName      string `json:"bank" gorm:"not null"`
	CountryISO2   string `json:"countryCode" gorm:"not null"`
	CountryName   string `json:"country" gorm:"not null"`
	IsHeadquarter bool   `json:"headquarter"`
	SwiftCode     string `gorm:"primaryKey;not null" json:"swift"`
}

func TestIncorrectTypeOfBankData(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := WrongBank{
		SwiftCode:     "IJKLPLPQAAB",
		BankName:      "Tester Bank",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Sidestreet 2",
		IsHeadquarter: false,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankBySwiftCodeSafe(bank.SwiftCode, baseURL)

	jsonData, err := json.Marshal(bank)
	if err != nil {
		t.Fatalf("Failed to marshal bank data to JSON: %v", err)
		return
	}

	res, err := http.Post(baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("HTTP POST request for bank failed: %v", err)
		return
	}

	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for wrong bank information")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}

	defer testutils.DeleteBankBySwiftCodeSafe(bank.SwiftCode, baseURL)
}

func TestInsertingBank_Duplicate(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	headquarterBank := models.Bank{
		SwiftCode:     "CBOTPLGHXXX",
		BankName:      "Central Bank of Testing",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bank st 1",
		IsHeadquarter: true,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankSafe(headquarterBank, baseURL)

	var res *http.Response = testutils.PostBank(assert, headquarterBank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusCreated, res.StatusCode, "Expected status code 201 for first bank creation")
	if http.StatusCreated != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}

	res = testutils.PostBank(assert, headquarterBank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for attempting to POST the same bank twice")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}

	defer testutils.DeleteBankSafe(headquarterBank, baseURL)
}

func TestInsertingBank_Negative_SwiftCodeSmall(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "cdcdPLGHXXX",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for swift code with lower letters")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}
}

func TestInsertingBank_Negative_SwiftCodeLength(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "ABCDPLGHXXXAA",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for swift code with incorrect length")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}
}

func TestInsertingBank_Negative_SwiftCodeContainsNumber(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "AB12PLGHXXX",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for SWIFT code with numbers in first 6 characters")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}
}

func TestInsertingBank_Negative_CountryLowerLetters(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "BBBGPLGHXXX",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for country with lower letters")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}

}

func TestInsertingBank_Negative_NotRealISO2(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "BBBGPLGHXXX",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "00",
		CountryName:   "POLAND",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for fake ISO2")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}
}

func TestInsertingBank_Negative_NotRealCountry_REALISO2(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "BBBGPLGHXXX",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "PL",
		CountryName:   "POLANDONIA",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for fake country (with real ISO2)")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}
}

func TestInsertingBank_Negative_NotRealCountry_FAKEISO2(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "BBBGPLGHXXX",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "00",
		CountryName:   "POLANDONIA",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for fake country (with fake ISO2)")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}
}

func TestInsertingBank_Negative_IncorrectCountryForISO2(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "BBBGDEGHXXX",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "DE",
		CountryName:   "POLAND",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for incorrect country for ISO2")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}
}

func TestInsertingBank_Negative_ISO2NotMatchingSwiftCode(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "BBBGPLGHXXX",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "DE",
		CountryName:   "GERMANY",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for ISO2 not matching the one in SWIFT code")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}

}

func TestInsertingBank_Negative_HeadquarterSwiftError(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "BBBGPLGHAAB",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bankowa 3",
		IsHeadquarter: true,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for headquarter without XXX in SWIFT code")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}
}

func TestInsertingBank_Negative_BranchSwiftError(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "BBBGPLGUXXX",
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bankowa 3",
		IsHeadquarter: false,
	}

	var res *http.Response = testutils.PostBank(assert, bank, baseURL)
	if res == nil {
		t.FailNow()
		return
	}
	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400 for branch with XXX SWIFT code")
	if http.StatusBadRequest != res.StatusCode {
		testutils.PrintMessageError(t, res)
		return
	}
}
