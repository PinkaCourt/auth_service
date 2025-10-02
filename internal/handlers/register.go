package handlers

import (
	"auth-service/internal/repository"
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

func isJSONContentType(contentType string) bool {
	return contentType == "application/json"
}

func dataValidator(Login string, Password string) error {
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

func parseRequestBody(body io.Reader) (*RegisterRequest, error) {
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

func passwordHash(password string) ([]byte, error) {
	bufPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return bufPass, err
}

func Register(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !isJSONContentType(contentType) {
		SendError(w, http.StatusUnsupportedMediaType, "требуется заголовок Content-Type: application/json")
		return
	}

	requestData, bodyError := parseRequestBody(r.Body)

	if bodyError != nil {

		SendError(w, http.StatusBadRequest, bodyError.Error())
		return
	}

	errData := dataValidator(requestData.Login, requestData.Password)

	if errData != nil {

		SendError(w, http.StatusBadRequest, errData.Error())
		return
	}

	hash, hashErr := passwordHash((requestData.Password))

	if hashErr != nil {

		SendError(w, http.StatusBadRequest, hashErr.Error())
		return
	}

	err := repository.RegisterUser(requestData.Login, hash)
	if err != nil {
		SendError(w, http.StatusConflict, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte("register OK"))
}
