package testutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Xenoneqq/swift-code-database/models"
	"github.com/stretchr/testify/assert"
)

func PostBank(assert *assert.Assertions, bank models.Bank, baseURL string) *http.Response {
	jsonData, err := json.Marshal(bank)
	assert.Nil(err, "Failed to marshal bank data to JSON: %v", err)
	if err != nil {
		return nil
	}

	res, err := http.Post(baseURL, "application/json", bytes.NewBuffer(jsonData))
	assert.Nil(err, "HTTP POST request for bank failed: %v", err)
	if err != nil {
		return nil
	}
	return res
}

func DeleteBank(assert *assert.Assertions, bank models.Bank, baseURL string) *http.Response {
	deleteURL := fmt.Sprintf("%s/%s", baseURL, bank.SwiftCode)
	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	assert.Nil(err, "Failed to create DELETE request for %s: %v", bank.SwiftCode, err)
	if err != nil {
		return nil
	}

	client := &http.Client{}
	res, err := client.Do(req)
	assert.Nil(err, "HTTP DELETE request for %s failed: %v", bank.SwiftCode, err)
	if err != nil {
		return nil
	}
	return res
}

func DeleteBankSafe(bank models.Bank, baseURL string) {
	deleteURL := fmt.Sprintf("%s/%s", baseURL, bank.SwiftCode)
	req, _ := http.NewRequest(http.MethodDelete, deleteURL, nil)
	client := &http.Client{}
	deleteRes, _ := client.Do(req)
	if deleteRes != nil {
		deleteRes.Body.Close()
	}
}

type BankResponse struct {
	Data    models.Bank `json:"data"`
	Message string      `json:"message"`
}

func ReadBankFromResponse(assert *assert.Assertions, res *http.Response) (models.Bank, error) {
	assert.NotNil(res, "Response is nil (failed to convert to Bank data)")
	if res == nil {
		return models.Bank{}, errors.New("HTTP response failed and is nil (failed to convert to Bank data)")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	assert.Nil(err, "Failed to read response body: %v", err)
	if err != nil {
		return models.Bank{}, err
	}

	var bankResponse BankResponse
	err = json.Unmarshal(bodyBytes, &bankResponse)
	assert.Nil(err, "Failed to unmarshal response body to BankResponse struct: %v", err)
	if err != nil {
		return models.Bank{}, err
	}

	return bankResponse.Data, nil
}

type HeadquarterResponse struct {
	Data    models.Headquarter `json:"data"`
	Message string             `json:"message"`
}

func ReadHeadquarterFromResponse(assert *assert.Assertions, res *http.Response) (models.Headquarter, error) {
	assert.NotNil(res, "Response is nil (failed to convert to Bank data)")
	if res == nil {
		return models.Headquarter{}, errors.New("HTTP response failed and is nil (failed to convert to Bank data)")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	assert.Nil(err, "Failed to read response body: %v", err)
	if err != nil {
		return models.Headquarter{}, err
	}

	var headquarterResponse HeadquarterResponse
	err = json.Unmarshal(bodyBytes, &headquarterResponse)
	assert.Nil(err, "Failed to unmarshal response body to BankResponse struct: %v", err)
	if err != nil {
		return models.Headquarter{}, err
	}

	return headquarterResponse.Data, nil
}
