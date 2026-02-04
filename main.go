package main

import (
	"net/http"
	"os"
	"time"
	"log"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received request")
		w.Write([]byte("deployed through githau2a2b: " + time.Now().Format(time.RFC3339)))
	})
	http.ListenAndServe(":"+port, nil)
}
