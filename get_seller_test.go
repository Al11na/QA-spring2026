package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/assert"
)

func (s *APISuite) TestTC201_GetItemsBySeller_Success(t provider.T) {
	t.Title("TC-201: Получение всех объявлений продавца")
	t.Description("Все объявления в ответе должны принадлежать указанному sellerID")
	t.Tags("GET", "positive")

	sellerID := generateSellerID()
	createItem(t, ItemRequest{
		Name:       "Товар А",
		Price:      1000,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 1, ViewCount: 1, Contacts: 1},
	})
	createItem(t, ItemRequest{
		Name:       "Товар Б",
		Price:      2000,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 2, ViewCount: 2, Contacts: 2},
	})

	resp, err := http.Get(fmt.Sprintf("%s/api/1/%d/item", baseURL, sellerID))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var items []ItemResponse
	err = json.NewDecoder(resp.Body).Decode(&items)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(items), 2)

	for _, item := range items {
		assert.Equal(t, sellerID, item.SellerID)
	}
}

func (s *APISuite) TestTC202_GetItemsBySeller_NotFound(t provider.T) {
	t.Title("TC-202: Несуществующий продавец")
	t.Description("Для несуществующего sellerID должен вернуться пустой массив")
	t.Tags("GET", "negative")

	resp, err := http.Get(fmt.Sprintf("%s/api/1/9999989/item", baseURL))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var items []ItemResponse
	err = json.NewDecoder(resp.Body).Decode(&items)
	assert.NoError(t, err)
	assert.Empty(t, items)
}

func (s *APISuite) TestTC203_GetItemsBySeller_InvalidSellerID(t provider.T) {
	t.Title("TC-203: Невалидный формат sellerID")
	t.Description("Сервер должен вернуть 400 при передаче строки вместо числа в sellerID")
	t.Tags("GET", "negative", "validation")

	resp, err := http.Get(fmt.Sprintf("%s/api/1/abc/item", baseURL))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
