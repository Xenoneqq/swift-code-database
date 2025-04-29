package main_test

import (
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
		BankName:      "Centralny Bank Testowy",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Bankowa 1",
		IsHeadquarter: true,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankSafe(headquarterBank, baseURL)

	res := testutils.PostBank(assert, headquarterBank, baseURL)
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
		BankName:      "Oddział Testowego Banku",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		Address:       "Boczna 2",
		IsHeadquarter: false,
	}

	// deleting bank if left from old tests
	defer testutils.DeleteBankSafe(branchBank, baseURL)

	res := testutils.PostBank(assert, branchBank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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

	res := testutils.PostBank(assert, bank, baseURL)
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
