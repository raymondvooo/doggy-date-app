package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	port, exists := os.LookupEnv("PORT")

	if !exists {
		port = "8080"
	}
	log.Printf("Starting server on port %s\n", port)

	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("hi"))
	})
	router.Route("/test", func(router chi.Router) {
		router.Get("/{name}", getName)
	})
	http.ListenAndServe(":"+port, router)

}

func getName(w http.ResponseWriter, req *http.Request) {
	//chi.UrlParam to read http url parameter
	name := chi.URLParam(req, "name")
	//string print
	w.Write([]byte(fmt.Sprintf("My name is %s", name)))
}
