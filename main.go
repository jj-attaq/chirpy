package main

import (
	"log"

	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	const port = "8080"
	const filepathRoot = "."

	srvMux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	srvMux.Handle("/app/", fsHandler)
	srvMux.HandleFunc("GET /api/healthz", handlerReadiness)
	srvMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	srvMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: srvMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
