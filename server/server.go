package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, router *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":8080", router)
}
