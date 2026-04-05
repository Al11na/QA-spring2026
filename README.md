# API Tests - qa-internship.avito.com

Автоматизированные тесты API на Go для микросервиса объявлений.

## Требования

- Go 1.24+
- golangci-lint (для статического анализа кода)
- Доступ к интернету (тесты обращаются к `https://qa-internship.avito.com`)

## Установка зависимостей

Клонировать репозиторий:

```bash
git clone https://github.com/Al11na/QA-spring2026.git
cd QA-spring2026
```

Установить зависимости Go:

```bash
go mod download
```

Установить golangci-lint:

```bash
# Windows (PowerShell)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

# Mac / Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.8
```

## Запуск тестов

Запустить все тесты:

```bash
go test -v
```

Запустить с таймаутом:

```bash
go test -v -timeout 120s
```

Запустить конкретный тест:

```bash
go test -v -run TestAPISuite/APISuite/Tests/TestTC001
```

## Allure отчёт

Установить Allure CLI (требуется Java):

```bash
# Windows (Scoop)
scoop install allure

# Mac
brew install allure
```

Запустить тесты с сохранением результатов:

```bash
# Windows (PowerShell)
$env:ALLURE_OUTPUT_PATH="."; go test -v -timeout 120s

# Mac / Linux
ALLURE_OUTPUT_PATH=. go test -v -timeout 120s
```

Открыть отчёт:

```bash
allure serve allure-results
```

## Статический анализ кода

```bash
golangci-lint run
```

Конфигурация линтера находится в `.golangci.yml`.

## Структура проекта

| Файл                | Описание                                  |
|---------------------|-------------------------------------------|
| `main_test.go`      | Структуры, хелперы, точка входа suite     |
| `post_test.go`      | Тесты POST /api/1/item (TC-001..TC-010)   |
| `get_item_test.go`  | Тесты GET /api/1/item/{id} (TC-101..TC-105) |
| `get_seller_test.go`| Тесты GET /api/1/{sellerID}/item (TC-201..TC-203) |
| `statistic_test.go` | Тесты GET /api/1/statistic/{id} (TC-301..TC-304) |
| `delete_test.go`    | Тесты DELETE /api/2/item/{id} (TC-401..TC-402) |
| `TESTCASES.md`      | Описание тест-кейсов                      |
| `BUGS.md`           | Найденные дефекты (9 багов)               |
| `Task1_bugs.md`     | Баги из задания 1 (UI)                    |
| `.golangci.yml`     | Конфигурация линтера                      |

## Покрытие

| Метод  | Эндпоинт                    | Тестов |
|--------|-----------------------------|--------|
| POST   | /api/1/item                 | 10     |
| GET    | /api/1/item/{id}            | 5      |
| GET    | /api/1/{sellerID}/item      | 3      |
| GET    | /api/1/statistic/{id}       | 4      |
| DELETE | /api/2/item/{id}            | 2      |

**Всего: 24 теста**

## Ожидаемые результаты

При запуске тестов часть из них упадёт - это ожидаемо, так как они фиксируют баги API:

| Тест                                              | Статус | Причина                                     |
|---------------------------------------------------|--------|---------------------------------------------|
| TestTC003_CreateItem_NegativePrice                | FAIL   | BUG-002: API принимает отрицательную цену   |
| TestTC304_GetStatistic_AfterDelete                | FAIL   | BUG-007: статистика доступна после удаления |
| TestTC401_DeleteItem_Success                      | FAIL   | BUG-008: DELETE не возвращает поле status   |
| TestTC006_CreateItem_MinSellerID                  | SKIP   | BUG-005: API зависает при sellerID=111111   |
| TestTC303_GetStatistic_NotFound                   | SKIP   | BUG-006: API возвращает 504 timeout         |

Подробное описание всех дефектов — в файле `BUGS.md`.

## Окружение

- API: `https://qa-internship.avito.com`
- OS: Windows 11
- Go: 1.24
