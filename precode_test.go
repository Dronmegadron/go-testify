package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Fatalf("expected status code: %d, got %d", http.StatusOK, status)
	}

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	if len(list) != totalCount {
		t.Errorf("expected cafe count: %d, got %d", totalCount, len(list))
	}
}

// Тест 1: Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerSuccessfulRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
	require.NoError(t, err, "Failed to create HTTP request") //Ошибка при создании запроса

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status 200 OK")   //Ожидался статус 200 OK
	assert.NotEmpty(t, responseRecorder.Body.String(), "Response should not be empty") //Ответ не должен быть пустым
}

// Тест 2: Город, который передаётся в параметре city, не поддерживается.
// Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerUnsupportedCity(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=2&city=unknown", nil)
	require.NoError(t, err, "Failed to create HTTP request") //Ошибка при создании запроса

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status 400 Bad Request") //Ожидался статус 400 Bad Request
	assert.Equal(t, "wrong city value", responseRecorder.Body.String(), "Incorrect error message")    //Неверное сообщение об ошибке
}

// Тест 3: Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerCountExceedsAvailable(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	require.NoError(t, err, "Failed to create HTTP request") //Ошибка при создании запроса

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status 200 OK")

	expectedResponse := strings.Join(cafeList["moscow"], ",")
	assert.Equal(t, expectedResponse, responseRecorder.Body.String(), "Response should contain all cafes") //Ответ должен содержать все кафе
}
