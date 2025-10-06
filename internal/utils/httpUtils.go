package utils

import (
	"encoding/json"
	"errors"
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

func IsJSONContentType(contentType string) bool {
	return contentType == "application/json"
}

func DataValidator(Login string, Password string) error {
	if Login == "" {
		return errors.New("логин не может быть пустым")
	}

	if Password == "" {
		return errors.New("пароль не может быть пустым")
	}

	if len(Password) < 8 || len(Password) > 72 {
		return errors.New("пароль должен быть между 8 и 72 символами")
	}

	return nil
}

func ParseRequestBody(body io.Reader) (*RegisterRequest, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, errors.New("ошибка чтения тела запроса")
	}

	if len(bodyBytes) == 0 {
		return nil, errors.New("тело запроса не может быть пустым")
	}

	var requestData RegisterRequest

	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		return nil, errors.New("некорректный формат JSON")
	}

	return &requestData, nil
}

func PasswordHash(password string) ([]byte, error) {
	bufPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return bufPass, err
}
