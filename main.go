package main

import (
	"auth-service/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type RequestData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	println("Server OPEN")
	defer println("Server CLOSE")
	r := chi.NewRouter()

	connectBD()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("i'm root"))
	})

	r.Post("/register", handlers.Register)

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {

		//Пользователь логинится → сравнивайте хеш из БД с введенным паролем через bcrypt.CompareHashAndPassword.

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("login"))
	})

	http.ListenAndServe(":3333", r)

}

func connectBD() {

	db, err := sqlx.Open("sqlite", "test.db")

	if err != nil {

		println("Не удалось открыть подключение:", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {

		println("Не удалось подключиться к базе:", err)
	}
	println("Подключено к SQLite!")

}
