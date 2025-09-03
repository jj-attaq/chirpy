package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	cfg := apiConfig{}
	const port = "8080"
	const filepathRoot = "."

	srvMux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	wrapped := cfg.middlewareMetricsInc(handler)

	srvMux.Handle("/app/", wrapped)
	srvMux.HandleFunc("/healthz", handlerReadiness)
	srvMux.HandleFunc("/metrics", cfg.handlerCountReqs)
	srvMux.HandleFunc("/reset", cfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: srvMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerCountReqs(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rw.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Add(-cfg.fileserverHits.Load())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
