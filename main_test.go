package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
)

const baseURL = "https://qa-internship.avito.com"

type Statistics struct {
	Contacts  int `json:"contacts"`
	Likes     int `json:"likes"`
	ViewCount int `json:"viewCount"`
}

type ItemRequest struct {
	Name       string     `json:"name"`
	Price      int        `json:"price"`
	SellerID   int        `json:"sellerID"`
	Statistics Statistics `json:"statistics"`
}

type ItemResponse struct {
	ID         string     `json:"id"`
	SellerID   int        `json:"sellerId"`
	Name       string     `json:"name"`
	Price      int        `json:"price"`
	Statistics Statistics `json:"statistics"`
	CreatedAt  string     `json:"createdAt"`
}

type StatisticsResponse struct {
	Likes     int `json:"likes"`
	ViewCount int `json:"viewCount"`
	Contacts  int `json:"contacts"`
}

func generateSellerID() int {
	return rand.Intn(999999-111112) + 111112
}

func createItem(t provider.T, req ItemRequest) string {
	t.Helper()
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Logf("Ошибка создания объявления: статус %d", resp.StatusCode)
		t.FailNow()
	}

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Logf("Ошибка декодирования ответа: %v", err)
		t.FailNow()
	}

	var itemID string
	fmt.Sscanf(result["status"], "Сохранили объявление - %s", &itemID)
	if itemID == "" {
		t.Logf("UUID объявления пустой. Ответ: %v", result)
		t.FailNow()
	}
	return itemID
}

func deleteItem(t provider.T, id string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/2/item/%s", baseURL, id), nil)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	return resp
}

type APISuite struct {
	suite.Suite
}

func TestAPISuite(t *testing.T) {
	runner.NewSuiteRunner(t, "api_tests", "APISuite", new(APISuite)).RunTests()
}
