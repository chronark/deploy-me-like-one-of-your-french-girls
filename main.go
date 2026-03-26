package main

import (
	"net/http"
	"os"
	"time"
	"log"
	"encoding/json"
	"fmt"
	"math"
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
	    http.HandleFunc("/compute", func(w http.ResponseWriter, r *http.Request) {
                result := 0.0
                for i := 0; i < 1_000_000; i++ {
                        for j := 0; j < 1_000_000; j++ {
                                for k := 0; k < 1_000_000; k++ {
                                        result += math.Sqrt(float64(i*j + k + 1))
                                }
                        }
                }
                fmt.Fprintf(w, "result: %f\n", result)
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

	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
