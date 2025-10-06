package handlers

import (
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	//просто ничего не делает
	// может редирект на авторизацию?
	// или админку повесить
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("i'm root"))
}
