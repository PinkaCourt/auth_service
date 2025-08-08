package handlers

// Шаг 4: HTTP-обработчики

// Вызов сервиса для регистрации

import (
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func isJSONContentType(contentType string) bool {
	return contentType == "application/json"
}

func dataValidator(Login string, Password string) string {
	err := ""

	if Login == "" {
		err = "Field Login is required"
	}

	if Password == "" {
		err = "Field Password is required"
	}

	if len(Password) < 8 || len(Password) > 72 {
		err = "Password must be between 8 and 72 characters"
	}

	return err

}

func parseRequestBody(body io.Reader) (*RegisterRequest, string) {
	// Читаем всё тело запроса
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, "Ошибка чтения тела запроса"
	}

	// Проверяем пустое тело
	if len(bodyBytes) == 0 {
		return nil, "empty request body"
	}

	// Декодируем JSON
	var requestData RegisterRequest

	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		return nil, "Invalid JSON format"

	}

	return &requestData, ""
}

func passwordHash(password string) ([]byte, string) {
	bufPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, "Password Too Long"
	}

	return bufPass, ""
}

func Register(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !isJSONContentType(contentType) {
		SendError(w, http.StatusUnsupportedMediaType, "only application/json is supported")
		return
	}

	requestData, bodyError := parseRequestBody(r.Body)

	if bodyError != "" {
		SendError(w, http.StatusBadRequest, bodyError)
		return
	}

	errData := dataValidator(requestData.Login, requestData.Password)

	if errData != "" {
		SendError(w, http.StatusBadRequest, errData)
		return
	}

	hash, hashErr := passwordHash((requestData.Password))

	if hashErr != "" {
		SendError(w, http.StatusBadRequest, hashErr)
		return
	}

	println(string(hash))

	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte("register OK"))
}

// Текущий фокус:
// Базовый рефакторинг без глубоких изменений
// Сохранение работоспособности кода
// Постепенное улучшение
// Отложенные задачи (напоминания для будущего):
//  Валидация: Возврат множественных ошибок (п.3)
//  Content-Type: Учет параметров в проверке (п.4)
//  Тело запроса: Восстановление r.Body (п.5)
//  Ошибки bcrypt: Детализация сообщений (п.6)
//  Длина пароля: Удаление верхней границы (п.7)
//  Структуры: Вынос в отдельный пакет (п.9)
// Следующие шаги:
//  нэйминг
// Переход к обработке ошибок
// Работа с логами и ответами (на финальных этапах)
