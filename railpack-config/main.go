package main

import (
	"net/http"
	"os"
	"runtime"
	"time"
)

// Config-file test case: no Dockerfile, but a railpack.json that pins go 1.24.
// The response includes runtime.Version() so the pinned toolchain is
// verifiable from the outside: it must report go1.24.x, not the latest.
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("railpack-config " + runtime.Version() + " " + time.Now().Format(time.RFC3339)))
	})
	http.ListenAndServe(":"+port, nil)
}
