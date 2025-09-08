package main

import (
	"encoding/json"
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
	srvMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: srvMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Body string
	}

	type errBody struct {
		Error string `json:"error"`
	}

	type validBody struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := body{}

	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("%s", err)
		// w.WriteHeader(500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	valid := validBody{
		Valid: true,
	}
	validJSON, err := json.Marshal(valid)
	if err != nil {
		newErr := errBody{
			Error: "something went wrong",
		}
		errJSON, err := json.Marshal(newErr)
		if err != nil {
			log.Printf("Error marshalling Error JSON: %s", err)
		}
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errJSON))
		return
	}

	if len(reqBody.Body) > 140 {
		newErr := errBody{
			Error: "Chirp is too long",
		}
		invalidLengthJSON, err := json.Marshal(newErr)
		if err != nil {
			log.Printf("Error marshalling Error JSON: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(invalidLengthJSON))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(validJSON)
}
