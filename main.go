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
		w.Write([]byte("deployed through redeploy2: " + time.Now().Format(time.RFC3339)))
	})

	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		
		w.Write([]byte(os.Environ()))
	})
	
	http.ListenAndServe(":"+port, nil)
}
