package main

import (
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/raymondvooo/doggy-date-app/server/api"
	"github.com/raymondvooo/doggy-date-app/server/gql"
	"github.com/raymondvooo/doggy-date-app/server/postgres"
)

func main() {
	port, exists := os.LookupEnv("PORT")

	if !exists {
		port = "8081"
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

	gqlSchema, err := getSchema("./gql/schema.graphql")
	if err != nil {
		panic(err)
	}

	schema := graphql.MustParseSchema(gqlSchema, &gql.Resolver{Db: db})

	router := chi.NewRouter()
	// Add some middleware to our router

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)

	router.Use(
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		// middleware.Logger,          // log api request calls
		middleware.DefaultCompress, // compress results, mostly gzipping assets and json
		// middleware.StripSlashes,    // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer, // recover from panics without crashing server
	)

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	router.Route("/test", func(router chi.Router) {
		router.Get("/{name}", getName)
	})

	// Create the graphql route with a Server method to handle it
	router.Handle("/graphql", &relay.Handler{Schema: schema})

	router.Route("/emailExists", func(router chi.Router) {
		router.Post("/", func(w http.ResponseWriter, req *http.Request) {
			api.CheckEmailExists(w, req, db)
		})
	})
	defer db.Close()
	http.ListenAndServe(":"+port, router)

}

func getName(w http.ResponseWriter, req *http.Request) {
	//chi.UrlParam to read http url parameter
	name := chi.URLParam(req, "name")
	//string print
	w.Write([]byte(fmt.Sprintf("My name is %s", name)))
}

func getSchema(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

var page = []byte(`
<!DOCTYPE html>
<html>

<head>
  <meta charset=utf-8/>
  <meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
  <title>GraphQL Playground</title>
  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
  <link rel="shortcut icon" href="//cdn.jsdelivr.net/npm/graphql-playground-react/build/favicon.png" />
  <script src="//cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
</head>

<body>
  <div id="root">
    <style>
      body {
        background-color: rgb(23, 42, 58);
        font-family: Open Sans, sans-serif;
        height: 90vh;
      }
      #root {
        height: 100%;
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      .loading {
        font-size: 32px;
        font-weight: 200;
        color: rgba(255, 255, 255, .6);
        margin-left: 20px;
      }
      img {
        width: 78px;
        height: 78px;
      }
      .title {
        font-weight: 400;
      }
    </style>
    <img src='//cdn.jsdelivr.net/npm/graphql-playground-react/build/logo.png' alt=''>
    <div class="loading"> Loading
      <span class="title">GraphQL Playground</span>
    </div>
  </div>
  <script>window.addEventListener('load', function (event) {
      GraphQLPlayground.init(document.getElementById('root'), {endpoint: '/graphql'})
    })</script>
</body>

</html>
`)
