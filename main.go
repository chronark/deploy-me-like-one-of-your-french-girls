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
		log.Println("received rseques2t ")
		w.Write([]byte("hello 44222 I hate regions!!! did flo break krane? " + "host="+r.URL.Host + time.Now().Format(time.RFC3339)))
	})



	http.HandleFunc("/v1/headers", func(w http.ResponseWriter, r *http.Request) {

		b, err := json.Marshal(r.Header)
		if err != nil {
			log.Println("error marshaling headers: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	})

	http.HandleFunc("/v1/abc", func(w http.ResponseWriter, r *http.Request) {
		log.Println( r.URL.Path)
		w.Write([]byte("hello from /v1/abc"))
	})
	http.HandleFunc("/v1/def", func(w http.ResponseWriter, r *http.Request) {
		log.Println( r.URL.Path)
		w.Write([]byte("hello from /v1/def"))
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Println( r.URL.Path)
		w.Write([]byte("hello from /v1/def"))
	})
	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {

		b, _ := json.Marshal(os.Environ())


		w.Write(b)
	})

	http.ListenAndServe(":"+port, nil)
}
