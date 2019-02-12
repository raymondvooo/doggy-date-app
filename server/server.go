package main

import (
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
	"github.com/minio/minio-go"
	"github.com/raymondvooo/doggy-date-app/server/api"
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

	// creds := credentials.NewStaticCredentials(os.Getenv("AWSAccessKeyId"), os.Getenv("AWSSecretKey"), "")
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region:      aws.String("us-west-1"),
	// 	Credentials: creds}))
	// cfg := aws.NewConfig().WithRegion("us-west-1").WithCredentials(creds).WithLogLevel(aws.LogDebug)
	// s3b := s3.New(sess, cfg)

	// creates a new AWS S3 client instance
	minioClient, err := minio.NewWithRegion("s3.amazonaws.com", os.Getenv("AWSAccessKeyId"), os.Getenv("AWSSecretKey"), true, "us-west-1")
	if err != nil {
		log.Fatalln(err)
	}

	//Create a new connection to our pg database
	db, err := postgres.NewConnection(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	//Load graphql Schema
	gqlSchema, err := getSchema("./gql/schema.graphql")
	if err != nil {
		panic(err)
	}

	//Parses graphql schema string into Schema object
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
		middleware.StripSlashes,    // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,       // recover from panics without crashing server
	)

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(graphQLPlayground)
	}))

	// Create the graphql route with a Server method to handle it
	router.Route("/graphql", func(router chi.Router) {
		router.Handle("/", &relay.Handler{Schema: schema})
		// router.Handle("/date", &relay.Handler{Schema: schema})
	})

	router.Route("/emailExists", func(router chi.Router) {
		router.Post("/", func(w http.ResponseWriter, req *http.Request) {
			api.CheckEmailExists(w, req, db)
		})
	})

	router.Route("/user", func(router chi.Router) {
		router.Route("/{uid}", func(router chi.Router) {
			pb := api.ProfileBuilder{}
			router.Route("/upload", func(router chi.Router) {
				router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					uid := chi.URLParam(req, "uid")
					pb = api.ProfileBuilder{ID: graphql.ID(uid)}
					w.Write(uploadTest)
				}))
				router.Handle("/send", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					pb.UploadAnyS3(w, req, minioClient, db, "users", pb.ID)
				}))
			})
		})
	})

	router.Route("/dog", func(router chi.Router) {
		router.Route("/{dogId}", func(router chi.Router) {
			pb := api.ProfileBuilder{}
			router.Route("/upload", func(router chi.Router) {
				router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					did := chi.URLParam(req, "dogId")
					pb = api.ProfileBuilder{ID: graphql.ID(did)}
					w.Write(uploadTest)
				}))
				router.Handle("/send", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					pb.UploadAnyS3(w, req, minioClient, db, "dogs", pb.ID)
				}))
			})
		})
	})

	defer db.Close()
	http.ListenAndServe(":"+port, router)

}

func getSchema(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

var graphQLPlayground = []byte(`
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

var uploadTest = []byte(`	
	<html>
<head>
       <title>Upload file</title>
</head>
<body>
<form enctype="multipart/form-data" action="https://doggy-date-go.herokuapp.com/upload/send" method="post">
    <input type="file" name="uploadfile" />
    <input type="hidden" name="token" value="{{.}}"/>
    <input type="submit" value="upload" />
</form>
</body>
</html>	
`)
