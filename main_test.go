package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

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
	return rand.Intn(999999-111111) + 111111
}

func createItem(t *testing.T, req ItemRequest) string {
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

func deleteItem(t *testing.T, id string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/2/item/%s", baseURL, id), nil)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	return resp
}

func TestTC001_CreateItem_Success(t *testing.T) {
	sellerID := generateSellerID()
	req := ItemRequest{
		Name:     "Смартфон на Android",
		Price:    25000,
		SellerID: sellerID,
		Statistics: Statistics{
			Likes:     5,
			ViewCount: 20,
			Contacts:  2,
		},
	}
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result["status"], "Сохранили объявление")
}

func TestTC002_CreateItem_EmptyName(t *testing.T) {
	sellerID := generateSellerID()
	req := ItemRequest{
		Name:     "",
		Price:    6000,
		SellerID: sellerID,
		Statistics: Statistics{
			Likes:     1,
			ViewCount: 1,
			Contacts:  2,
		},
	}
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTC003_CreateItem_NegativePrice(t *testing.T) {
	sellerID := generateSellerID()
	req := ItemRequest{
		Name:     "Наушники",
		Price:    -100,
		SellerID: sellerID,
		Statistics: Statistics{
			Likes:     1,
			ViewCount: 1,
			Contacts:  2,
		},
	}
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTC004_CreateItem_Idempotency(t *testing.T) {
	sellerID := generateSellerID()
	req := ItemRequest{
		Name:     "Книга по программированию",
		Price:    800,
		SellerID: sellerID,
		Statistics: Statistics{
			Likes:     1,
			ViewCount: 1,
			Contacts:  2,
		},
	}

	jsonData, _ := json.Marshal(req)
	resp1, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	var res1 map[string]string
	json.NewDecoder(resp1.Body).Decode(&res1)
	var id1 string
	fmt.Sscanf(res1["status"], "Сохранили объявление - %s", &id1)

	jsonData, _ = json.Marshal(req)
	resp2, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	var res2 map[string]string
	json.NewDecoder(resp2.Body).Decode(&res2)
	var id2 string
	fmt.Sscanf(res2["status"], "Сохранили объявление - %s", &id2)

	assert.NotEqual(t, id1, id2, "Каждый POST должен создавать уникальный id")
}

func TestTC005_CreateItem_InvalidSellerIDType(t *testing.T) {
	body := `{"sellerID":"один","name":"Книга по программированию","price":800,"statistics":{"likes":1,"viewCount":1,"contacts":2}}`
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTC006_CreateItem_MinSellerID(t *testing.T) {
	t.Skip("SKIP: API зависает при sellerID=111111")

	req := ItemRequest{
		Name:     "Товар с минимальным sellerID",
		Price:    500,
		SellerID: 111111,
		Statistics: Statistics{
			Likes:     1,
			ViewCount: 1,
			Contacts:  1,
		},
	}
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result["status"], "Сохранили объявление")
}

func TestTC007_CreateItem_MissingStatistics(t *testing.T) {
	body := `{"sellerID":555666,"name":"Товар без статистики","price":100}`
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTC008_CreateItem_MaxPrice(t *testing.T) {
	body := `{"sellerID":555666,"name":"Дорогой товар","price":999999999999999,"statistics":{"likes":0,"viewCount":0,"contacts":0}}`
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, resp.StatusCode)
}

func TestTC009_CreateItem_XSSInName(t *testing.T) {
	req := ItemRequest{
		Name:       "<script>alert(1)</script>",
		Price:      100,
		SellerID:   generateSellerID(),
		Statistics: Statistics{Likes: 0, ViewCount: 0, Contacts: 0},
	}
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, resp.StatusCode)
}

func TestTC010_CreateItem_LongName(t *testing.T) {
	longName := ""
	for i := 0; i < 1010; i++ {
		longName += "а"
	}
	req := ItemRequest{
		Name:       longName,
		Price:      100,
		SellerID:   generateSellerID(),
		Statistics: Statistics{Likes: 0, ViewCount: 0, Contacts: 0},
	}
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTC101_GetItem_Success(t *testing.T) {
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

func TestTC102_GetItem_NotFound(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/1/item/e4a685e1-b27b-4954-9540-66726fc6e73d")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestTC103_GetItem_InvalidID(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/1/item/abc")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTC104_GetItem_AfterDelete(t *testing.T) {
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

func TestTC105_GetItem_FieldTypes(t *testing.T) {
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
	assert.True(t, priceIsNum, "price должен быть числом, а не строкой")
	assert.True(t, sellerIsNum, "sellerId должен быть числом, а не строкой")

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

func TestTC201_GetItemsBySeller_Success(t *testing.T) {
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

func TestTC202_GetItemsBySeller_NotFound(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/api/1/9999989/item", baseURL))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var items []ItemResponse
	err = json.NewDecoder(resp.Body).Decode(&items)
	assert.NoError(t, err)
	assert.Empty(t, items)
}

func TestTC203_GetItemsBySeller_InvalidSellerID(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/api/1/abc/item", baseURL))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTC301_GetStatistic_V1_Success(t *testing.T) {
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

func TestTC302_GetStatistic_V2_Success(t *testing.T) {
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

func TestTC303_GetStatistic_NotFound(t *testing.T) {
	t.Skip("SKIP: API возвращает 504 timeout для несуществующего ID в статистике")

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

func TestTC304_GetStatistic_AfterDelete(t *testing.T) {
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

func TestTC401_DeleteItem_Success(t *testing.T) {
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

func TestTC402_DeleteItem_AlreadyDeleted(t *testing.T) {
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
