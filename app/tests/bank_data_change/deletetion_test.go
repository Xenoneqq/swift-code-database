package main_test

import (
	"net/http"
	"testing"

	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/Xenoneqq/swift-code-database/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBankDeletion(t *testing.T) {
	assert := assert.New(t)
	baseURL := "http://app:8080/api/v1/swift-codes"

	bank := models.Bank{
		SwiftCode:     "BOTSDEGHXXX",
		BankName:      "Bank of Testing",
		CountryISO2:   "DE",
		CountryName:   "GERMANY",
		Address:       "Bank Street 7",
		IsHeadquarter: true,
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

	var resDel *http.Response = testutils.DeleteBank(assert, bank, baseURL)
	if resDel == nil {
		t.Fatalf("Failed to delete test bank: DeeteBank returned nil response")
	}

	passed = assert.Equal(http.StatusOK, resDel.StatusCode, "Expected status code 200 for deleting test bank")
	if !passed {
		testutils.PrintMessageError(t, resDel)
		return
	}
}
