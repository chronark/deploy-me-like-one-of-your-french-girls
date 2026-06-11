package main

import (
	"net/http"
	"os"
	"runtime"
	"time"
)

// Zero-config test case: no Dockerfile, no railpack.json. The platform must
// detect Go from go.mod and infer the build and start commands.
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("railpack-auto " + runtime.Version() + " " + time.Now().Format(time.RFC3339)))
	})
	http.ListenAndServe(":"+port, nil)
}
