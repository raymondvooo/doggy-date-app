package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/raymondvooo/doggy-date-app/server/postgres"

	"github.com/go-chi/chi"
)

func main() {
	port, exists := os.LookupEnv("PORT")

	if !exists {
		port = "8080"
	}
	log.Printf("Starting server on port %s\n", port)

	localEnvError := godotenv.Load()
	if localEnvError != nil {
		log.Fatal("Error loading .env file")
	}

	pgDb := os.Getenv("DATABASE_URL")
	if len(pgDb) == 0 {
		log.Fatal("Invalid database url")
		return
	}

	//Create a new connection to our pg database
	db, err := postgres.NewConnection(pgDb)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to db", db)

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
