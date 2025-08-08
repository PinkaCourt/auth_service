package handlers

// handlers/register.go:
// Прими JSON с login и password
// Проверь, что логин не пустой
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)


func main() {
	r := chi.NewRouter()
//	r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)


	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("register"))




// 		Как:
// Проверить, что заголовок Content-Type равен application/json.
// Создать структуру-контейнер (например, RegisterRequest), куда будут парситься данные.
// Использовать стандартный пакет encoding/json для преобразования тела запроса в эту структуру.
// Если JSON невалидный (например, поля не совпадают), вернуть ошибку 400 Bad Request.



		//Получить тело HTTP-запроса (JSON) и преобразовать его в объект с полями login и password.


	})



	http.ListenAndServe(":3333", r)


	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == "GET" {
	// 		w.Write([]byte("root."))
	// 	}
	// })


}
