package main

import (
	"auth-service/internal/handlers"
	"auth-service/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	r.Post("/login", handlers.LoginHandler)

	http.ListenAndServe(":3333", r)
}
