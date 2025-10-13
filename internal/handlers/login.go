package handlers

import (
	"auth-service/internal/repository"
	"auth-service/internal/utils"
	"errors"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	requestData, bodyError := utils.CheckRequest(r)

	if bodyError != nil {
		utils.SendError(w, http.StatusBadRequest, bodyError.Error())
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
	//todo!!!!!
	if !isEqual {
		err := errors.New("неверный пароль")
		utils.SendError(w, http.StatusUnauthorized, err.Error()) //todo!
	} else {
		w.Write([]byte("login"))
	}

}
