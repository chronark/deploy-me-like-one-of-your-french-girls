package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// mysqlConfig converts the mysql://user:password@host/database URL that Unkey
// sets for a PlanetScale connection into a go-sql-driver config. The URL's ssl
// query parameter is in Node format, so the config always verifies TLS instead.
func mysqlConfig(raw string) (*mysql.Config, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	if u.Scheme != "mysql" {
		return nil, fmt.Errorf("scheme is %q, want mysql", u.Scheme)
	}

	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = u.Host
	if u.Port() == "" {
		cfg.Addr = u.Host + ":3306"
	}
	cfg.User = u.User.Username()
	cfg.Passwd, _ = u.User.Password()
	cfg.DBName = strings.TrimPrefix(u.Path, "/")
	cfg.TLSConfig = "true"
	cfg.ParseTime = true

	return cfg, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	dbURL := os.Getenv("ANDREAS_DB_URL")
	if dbURL == "" {
		log.Fatal("ANDREAS_DB_URL is not set")
	}
	dbConfig, err := mysqlConfig(dbURL)
	if err != nil {
		log.Fatalf("invalid ANDREAS_DB_URL: %v", err)
	}
	connector, err := mysql.NewConnector(dbConfig)
	if err != nil {
		log.Fatalf("invalid ANDREAS_DB_URL: %v", err)
	}
	db := sql.OpenDB(connector)
	defer db.Close()

	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		var version string
		if err := db.QueryRowContext(r.Context(), "SELECT VERSION()").Scan(&version); err != nil {
			log.Println("error querying andreas-db: ", err)
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintf(w, "connected to andreas-db, version=%s\n", version)
	})

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(map[string]any{
			"hello": "world",
		})
		if err != nil {
			log.Println("error marshaling headers: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("json:")
		fmt.Println("my coool message " + string(b))
		w.Write([]byte("I logged stuff"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received edited")
		w.Write([]byte("hello deployed after rolled back " + "host=" + r.URL.Host + time.Now().Format(time.RFC3339)))
	})
	// Echo the raw query exactly as received, to check what the gateway forwards.
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("query rawQuery=%q", r.URL.RawQuery)
		fmt.Fprintf(w, "rawQuery=%q principal=%q\n", r.URL.RawQuery, r.Header.Get("X-Unkey-Principal"))
	})
	http.HandleFunc("/compute", func(w http.ResponseWriter, r *http.Request) {
		result := 0.0
		for i := 0; i < 1_000; i++ {
			log.Printf("i=%d", i)
			for j := 0; j < 1_000; j++ {
				result += math.Sqrt(float64(i*j + 1))
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
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://jsonplaceholder.typicode.com/todos/1", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	})

	http.HandleFunc("/v1/abc", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.Write([]byte("hello from /v1/abc"))
	})
	http.HandleFunc("/v1/def", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.Write([]byte("hello from /v1/def"))
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.Write([]byte("hello from /v1/def"))
	})

	http.HandleFunc("/region", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(os.Getenv("UNKEY_REGION")))
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
