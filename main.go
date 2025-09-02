package main

import "net/http"

func main() {
	srvMux := http.NewServeMux()

	srv := http.Server{
		Addr:    ":8080",
		Handler: srvMux,
	}

	srv.ListenAndServe()
}
