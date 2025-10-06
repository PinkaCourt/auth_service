package handlers

import (
	"auth-service/internal/repository"
	"auth-service/internal/utils"
	"errors"
	"net/http"
)

// todo! проверить все хттп статусы
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	//todo! проверки как в регистрации убрать

	contentType := r.Header.Get("Content-Type")
	if !utils.IsJSONContentType(contentType) {
		utils.SendError(w, http.StatusUnsupportedMediaType, "требуется заголовок Content-Type: application/json")
		return
	}

	requestData, bodyError := utils.ParseRequestBody(r.Body)

	if bodyError != nil {

		utils.SendError(w, http.StatusBadRequest, bodyError.Error())
		return
	}

	errData := utils.DataValidator(requestData.Login, requestData.Password)

	if errData != nil {

		utils.SendError(w, http.StatusBadRequest, errData.Error())
		return
	}

	hash, hashErr := utils.PasswordHash((requestData.Password))

	if hashErr != nil {

		utils.SendError(w, http.StatusBadRequest, hashErr.Error())
		return
	}

	exists, err := repository.CheckUserExists(requestData.Login)
	if err != nil {

		utils.SendError(w, http.StatusConflict, err.Error())
	}
	if !exists {
		println("Пользователь с таким логином НЕ существует")

		utils.SendError(w, http.StatusConflict, err.Error())
		return
	}

	isEqual, err := repository.CheckPassExists(requestData.Login, hash)

	w.Header().Set("Content-Type", "text/plain")

	if err != nil {
		utils.SendError(w, http.StatusConflict, err.Error())
	}

	if !isEqual {
		err := errors.New("неверный пароль")

		utils.SendError(w, http.StatusUnauthorized, err.Error())
	} else {
		w.Write([]byte("login"))
	}

}
