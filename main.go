package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/runner989/RSS_Blog_Aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL is not found in the environment")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database")
	}

	apiCfg := apiConfig{
		DB: database.New(db),
	}
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})

	mux.HandleFunc("GET /v1/healthz", readinessHandler)
	mux.HandleFunc("POST /v1/users", apiCfg.createUsersHandler)
	mux.HandleFunc("GET /v1/err", errorHandler)

	log.Printf("Starting RSS Feed Aggregator on port %s", port)
	log.Fatal(server.ListenAndServe())
}
