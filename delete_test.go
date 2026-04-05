package main

import (
	"encoding/json"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/assert"
)

func (s *APISuite) TestTC401_DeleteItem_Success(t provider.T) {
	t.Title("TC-401: Удаление существующего объявления")
	t.Description("Успешное удаление должно возвращать 200")
	t.Tags("DELETE", "positive")

	sellerID := generateSellerID()
	itemID := createItem(t, ItemRequest{
		Name:       "Товар для удаления",
		Price:      3000,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 1, ViewCount: 1, Contacts: 1},
	})

	resp := deleteItem(t, itemID)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(t, result)
}

func (s *APISuite) TestTC402_DeleteItem_AlreadyDeleted(t provider.T) {
	t.Title("TC-402: Повторное удаление (идемпотентность DELETE)")
	t.Description("Повторное удаление уже удалённого объявления должно возвращать 404")
	t.Tags("DELETE", "corner-case", "idempotency")

	sellerID := generateSellerID()
	itemID := createItem(t, ItemRequest{
		Name:       "Товар для двойного удаления",
		Price:      1500,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 1, ViewCount: 1, Contacts: 1},
	})

	resp1 := deleteItem(t, itemID)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	resp2 := deleteItem(t, itemID)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)
}
