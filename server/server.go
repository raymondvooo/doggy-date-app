package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
	"github.com/raymondvooo/doggy-date-app/server/gql"
	"github.com/raymondvooo/doggy-date-app/server/postgres"
)

func main() {
	port, exists := os.LookupEnv("PORT")

	if !exists {
		port = "8080"
	}
	log.Printf("Starting server on port %s\n", port)

	//if developing locally
	localEnvError := godotenv.Load()
	if localEnvError != nil {
		log.Println("Error loading .env file")
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

	// Create our root query for graphql
	rootQuery := gql.NewRoot(db)
	// Create a new graphql schema, passing in the the root query
	sc, err := graphql.NewSchema(
		graphql.SchemaConfig{Query: rootQuery.Query},
	)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
	}
	// Create a server struct that holds a pointer to our database as well
	// as the address of our graphql schema
	gqlServer := gql.Server{
		GqlSchema: &sc,
	}

	router := chi.NewRouter()
	// Add some middleware to our router
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger,          // log api request calls
		middleware.DefaultCompress, // compress results, mostly gzipping assets and json
		// middleware.StripSlashes,    // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer, // recover from panics without crashing server
	)

	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("hi"))
	})
	router.Route("/test", func(router chi.Router) {
		router.Get("/{name}", getName)
	})

	// Create the graphql route with a Server method to handle it
	router.Post("/graphql", gqlServer.GraphQL())
	defer db.Close()
	http.ListenAndServe(":"+port, router)

}

func getName(w http.ResponseWriter, req *http.Request) {
	//chi.UrlParam to read http url parameter
	name := chi.URLParam(req, "name")
	//string print
	w.Write([]byte(fmt.Sprintf("My name is %s", name)))
}
