package main

import "net/http"

func (cfg *apiConfig) handlerReset(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rw.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Add(-cfg.fileserverHits.Load())
}
