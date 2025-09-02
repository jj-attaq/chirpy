package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	srvMux := http.NewServeMux()
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	srvMux.Handle("/app/", fileServer)
	srvMux.HandleFunc("/healthz", healthzHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: srvMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func healthzHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}
