package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/assert"
)

func (s *APISuite) TestTC001_CreateItem_Success(t provider.T) {
	t.Title("TC-001: Успешное создание объявления")
	t.Description("Проверка создания объявления с валидными данными")
	t.Tags("POST", "positive")

	sellerID := generateSellerID()
	req := ItemRequest{
		Name:     "Смартфон на Android",
		Price:    25000,
		SellerID: sellerID,
		Statistics: Statistics{Likes: 5, ViewCount: 20, Contacts: 2},
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

func (s *APISuite) TestTC002_CreateItem_EmptyName(t provider.T) {
	t.Title("TC-002: Создание объявления с пустым именем")
	t.Description("Сервер должен вернуть 400 при пустом поле name")
	t.Tags("POST", "negative", "validation")

	req := ItemRequest{
		Name:     "",
		Price:    6000,
		SellerID: generateSellerID(),
		Statistics: Statistics{Likes: 1, ViewCount: 1, Contacts: 2},
	}
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func (s *APISuite) TestTC003_CreateItem_NegativePrice(t provider.T) {
	t.Title("TC-003: Создание объявления с отрицательной ценой")
	t.Description("Сервер должен вернуть 400 при отрицательной цене")
	t.Tags("POST", "negative", "validation")

	req := ItemRequest{
		Name:     "Наушники",
		Price:    -100,
		SellerID: generateSellerID(),
		Statistics: Statistics{Likes: 1, ViewCount: 1, Contacts: 2},
	}
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func (s *APISuite) TestTC004_CreateItem_Idempotency(t provider.T) {
	t.Title("TC-004: Идемпотентность POST")
	t.Description("Два одинаковых POST-запроса должны создавать объявления с разными UUID")
	t.Tags("POST", "corner-case", "idempotency")

	sellerID := generateSellerID()
	req := ItemRequest{
		Name:     "Книга по программированию",
		Price:    800,
		SellerID: sellerID,
		Statistics: Statistics{Likes: 1, ViewCount: 1, Contacts: 2},
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

	assert.NotEqual(t, id1, id2)
}

func (s *APISuite) TestTC005_CreateItem_InvalidSellerIDType(t provider.T) {
	t.Title("TC-005: Неверный тип sellerID (строка вместо числа)")
	t.Description("Сервер должен вернуть 400 при передаче строки в поле sellerID")
	t.Tags("POST", "negative", "validation")

	body := `{"sellerID":"один","name":"Книга по программированию","price":800,"statistics":{"likes":1,"viewCount":1,"contacts":2}}`
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func (s *APISuite) TestTC006_CreateItem_MinSellerID(t provider.T) {
	t.Title("TC-006: Граничное значение sellerID = 111111")
	t.Description("Нижняя граница допустимого диапазона sellerID должна приниматься")
	t.Tags("POST", "boundary")

	t.Skip()

	req := ItemRequest{
		Name:     "Товар с минимальным sellerID",
		Price:    500,
		SellerID: 111111,
		Statistics: Statistics{Likes: 1, ViewCount: 1, Contacts: 1},
	}
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func (s *APISuite) TestTC007_CreateItem_MissingStatistics(t provider.T) {
	t.Title("TC-007: Отсутствует обязательное поле statistics")
	t.Description("Сервер должен вернуть 400 при отсутствии поля statistics")
	t.Tags("POST", "negative", "validation")

	body := `{"sellerID":555666,"name":"Товар без статистики","price":100}`
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func (s *APISuite) TestTC008_CreateItem_MaxPrice(t provider.T) {
	t.Title("TC-008: Максимальная цена (стресс-тест)")
	t.Description("Сервер должен корректно обработать очень большое значение price")
	t.Tags("POST", "boundary", "stress")

	body := `{"sellerID":555666,"name":"Дорогой товар","price":999999999999999,"statistics":{"likes":0,"viewCount":0,"contacts":0}}`
	resp, err := http.Post(baseURL+"/api/1/item", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, resp.StatusCode)
}

func (s *APISuite) TestTC009_CreateItem_XSSInName(t provider.T) {
	t.Title("TC-009: XSS-символы в поле name")
	t.Description("Сервер должен принять с экранированием или отклонить XSS в name")
	t.Tags("POST", "security")

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

func (s *APISuite) TestTC010_CreateItem_LongName(t provider.T) {
	t.Title("TC-010: Имя длиной 1000+ символов")
	t.Description("Сервер должен отклонить слишком длинное значение поля name")
	t.Tags("POST", "boundary", "negative")

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
