package main

import (
	"auth-service/internal/handlers"
	"auth-service/internal/repository"
	//"log" // Используем log для более информативного вывода
	"net/http"

	"github.com/go-chi/chi/v5"
	//	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type RequestData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	println("Server OPEN")

	repository.ConnectDB()

	defer repository.DB.Close()

	repository.CreateUsersTable()

	defer println("Server CLOSE")

	r := chi.NewRouter()

	r.Get("/", handlers.RootHandler)
	r.Post("/register", handlers.Register)
	r.Get("/login", handlers.LoginHandler)

	http.ListenAndServe(":3333", r)
}

// todo! Репозиторий для работы с БД

func AddUser(login string, pass []byte) {
	checkLogin(login) // todo! если не проходит то не регистрируем
	addUserToBD(login, pass)

	println(login, "AddUser-OK")
}

func addUserToBD(login string, pass []byte) {
	println("repository.addUserToBD: Имитация записи в БД для пользователя", login)
}

func checkLogin(login string) {
	println("repository.checkLogin: Имитация проверки логина", login, "в базе.")

	// CheckUsersLoginInTable(requestData.Login) //todo!
}
