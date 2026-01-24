package main

import (
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().Format(time.RFC3339)))
	})
	http.ListenAndServe(":"+port, nil)
}
