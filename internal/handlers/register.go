package handlers

import (
	"auth-service/internal/repository"
	"auth-service/internal/utils"
	"net/http"

	_ "modernc.org/sqlite"
)

func Register(w http.ResponseWriter, r *http.Request) {
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

	err := repository.RegisterUser(requestData.Login, hash)
	if err != nil {
		utils.SendError(w, http.StatusConflict, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte("register OK"))
}
