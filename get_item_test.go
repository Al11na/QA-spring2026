package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/assert"
)

func (s *APISuite) TestTC101_GetItem_Success(t provider.T) {
	t.Title("TC-101: Получение существующего объявления")
	t.Description("Проверка структуры и полей ответа при получении объявления по ID")
	t.Tags("GET", "positive")

	sellerID := generateSellerID()
	itemID := createItem(t, ItemRequest{
		Name:       "Ноутбук",
		Price:      80000,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 3, ViewCount: 10, Contacts: 1},
	})

	resp, err := http.Get(fmt.Sprintf("%s/api/1/item/%s", baseURL, itemID))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var items []ItemResponse
	err = json.NewDecoder(resp.Body).Decode(&items)
	assert.NoError(t, err)
	assert.NotEmpty(t, items)

	item := items[0]
	assert.NotEmpty(t, item.ID)
	assert.NotEmpty(t, item.CreatedAt)
	assert.Equal(t, sellerID, item.SellerID)
	assert.Equal(t, "Ноутбук", item.Name)
	assert.Equal(t, 80000, item.Price)
}

func (s *APISuite) TestTC102_GetItem_NotFound(t provider.T) {
	t.Title("TC-102: Получение несуществующего объявления")
	t.Description("Сервер должен вернуть 404 для несуществующего UUID")
	t.Tags("GET", "negative")

	resp, err := http.Get(baseURL + "/api/1/item/e4a685e1-b27b-4954-9540-66726fc6e73d")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func (s *APISuite) TestTC103_GetItem_InvalidID(t provider.T) {
	t.Title("TC-103: Получение объявления с невалидным ID")
	t.Description("Сервер должен вернуть 400 при передаче строки вместо UUID")
	t.Tags("GET", "negative", "validation")

	resp, err := http.Get(baseURL + "/api/1/item/abc")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func (s *APISuite) TestTC104_GetItem_AfterDelete(t provider.T) {
	t.Title("TC-104: Получение удалённого объявления (E2E)")
	t.Description("После удаления объявление должно возвращать 404")
	t.Tags("GET", "DELETE", "e2e", "negative")

	sellerID := generateSellerID()
	itemID := createItem(t, ItemRequest{
		Name:       "Велосипед",
		Price:      15000,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 1, ViewCount: 1, Contacts: 1},
	})

	delResp := deleteItem(t, itemID)
	defer delResp.Body.Close()
	assert.Equal(t, http.StatusOK, delResp.StatusCode)

	resp, err := http.Get(fmt.Sprintf("%s/api/1/item/%s", baseURL, itemID))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func (s *APISuite) TestTC105_GetItem_FieldTypes(t provider.T) {
	t.Title("TC-105: Типы данных полей в ответе")
	t.Description("Поля price, sellerId, likes, viewCount, contacts должны быть числами")
	t.Tags("GET", "positive", "validation")

	sellerID := generateSellerID()
	itemID := createItem(t, ItemRequest{
		Name:       "Смартфон на Android",
		Price:      25000,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 5, ViewCount: 20, Contacts: 2},
	})

	resp, err := http.Get(fmt.Sprintf("%s/api/1/item/%s", baseURL, itemID))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var raw []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&raw)
	assert.NoError(t, err)
	assert.NotEmpty(t, raw)

	item := raw[0]
	_, priceIsNum := item["price"].(float64)
	_, sellerIsNum := item["sellerId"].(float64)
	assert.True(t, priceIsNum, "price должен быть числом")
	assert.True(t, sellerIsNum, "sellerId должен быть числом")

	if stats, ok := item["statistics"].(map[string]interface{}); ok {
		_, likesIsNum := stats["likes"].(float64)
		_, viewIsNum := stats["viewCount"].(float64)
		_, contactsIsNum := stats["contacts"].(float64)
		assert.True(t, likesIsNum, "likes должен быть числом")
		assert.True(t, viewIsNum, "viewCount должен быть числом")
		assert.True(t, contactsIsNum, "contacts должен быть числом")
	} else {
		t.Error("поле statistics отсутствует или имеет неверный тип")
	}
}
