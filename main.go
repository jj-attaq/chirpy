package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	srvMux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: srvMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
