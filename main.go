package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	srvMux := http.NewServeMux()
	srvMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	srvMux.HandleFunc("/healthz", handlerReadiness)

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
