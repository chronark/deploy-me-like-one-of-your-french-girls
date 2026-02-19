package main

import (
	"net/http"
	"os"
	"time"
	"log"
	"encoding/json"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	} 
	

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received reques2t ")
		w.Write([]byte("behind auth: " + time.Now().Format(time.RFC3339)))
	})

	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {

		b, _ := json.Marshal(os.Environ())


		w.Write(b)
	})

	http.ListenAndServe(":"+port, nil)
}
