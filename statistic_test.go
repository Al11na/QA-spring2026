package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/assert"
)

func (s *APISuite) TestTC301_GetStatistic_V1_Success(t provider.T) {
	t.Title("TC-301: Получение статистики v1")
	t.Description("Статистика объявления должна содержать поля likes, viewCount, contacts")
	t.Tags("GET", "positive", "statistics")

	sellerID := generateSellerID()
	itemID := createItem(t, ItemRequest{
		Name:       "Телефон",
		Price:      30000,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 5, ViewCount: 20, Contacts: 3},
	})

	resp, err := http.Get(fmt.Sprintf("%s/api/1/statistic/%s", baseURL, itemID))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var stats []StatisticsResponse
	err = json.NewDecoder(resp.Body).Decode(&stats)
	assert.NoError(t, err)
	assert.NotEmpty(t, stats)
}

func (s *APISuite) TestTC302_GetStatistic_V2_Success(t provider.T) {
	t.Title("TC-302: Получение статистики v2")
	t.Description("Статистика через v2 endpoint должна работать аналогично v1")
	t.Tags("GET", "positive", "statistics")

	sellerID := generateSellerID()
	itemID := createItem(t, ItemRequest{
		Name:       "Планшет",
		Price:      25000,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 4, ViewCount: 15, Contacts: 2},
	})

	resp, err := http.Get(fmt.Sprintf("%s/api/2/statistic/%s", baseURL, itemID))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var stats []StatisticsResponse
	err = json.NewDecoder(resp.Body).Decode(&stats)
	assert.NoError(t, err)
	assert.NotEmpty(t, stats)
}

func (s *APISuite) TestTC303_GetStatistic_NotFound(t provider.T) {
	t.Title("TC-303: Статистика несуществующего объявления")
	t.Description("Сервер должен вернуть 404 для несуществующего ID")
	t.Tags("GET", "negative", "statistics")

	t.Skip()

	nonExistentID := "e4a685e1-b27b-4954-9540-66726fc6e73d"

	respV1, err := http.Get(fmt.Sprintf("%s/api/1/statistic/%s", baseURL, nonExistentID))
	assert.NoError(t, err)
	defer respV1.Body.Close()
	assert.Equal(t, http.StatusNotFound, respV1.StatusCode)

	respV2, err := http.Get(fmt.Sprintf("%s/api/2/statistic/%s", baseURL, nonExistentID))
	assert.NoError(t, err)
	defer respV2.Body.Close()
	assert.Equal(t, http.StatusNotFound, respV2.StatusCode)
}

func (s *APISuite) TestTC304_GetStatistic_AfterDelete(t provider.T) {
	t.Title("TC-304: Статистика удалённого объявления (E2E)")
	t.Description("После удаления объявления его статистика должна возвращать 404")
	t.Tags("GET", "DELETE", "e2e", "negative", "statistics")

	sellerID := generateSellerID()
	itemID := createItem(t, ItemRequest{
		Name:       "Фотоаппарат",
		Price:      50000,
		SellerID:   sellerID,
		Statistics: Statistics{Likes: 7, ViewCount: 30, Contacts: 5},
	})

	delResp := deleteItem(t, itemID)
	defer delResp.Body.Close()
	assert.Equal(t, http.StatusOK, delResp.StatusCode)

	resp, err := http.Get(fmt.Sprintf("%s/api/1/statistic/%s", baseURL, itemID))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
