package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/chirpy-server-go/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	apiCfg := apiConfig{}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
	}
	apiCfg.dbQueries = database.New(db)

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerGetMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	wrappedMux := middlewareLog(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: wrappedMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
