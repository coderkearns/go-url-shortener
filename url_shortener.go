package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func uuid(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

type urlRecord struct {
	Url   string `json:"url"`
	Short string `json:"s"`
}

var records = []urlRecord{
	{"http://localhost:8080/", "ah4d"},
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	record := urlRecord{
		r.URL.Query().Get("url"),
		uuid(4),
	}

	records = append(records, record)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	for _, record := range records {
		if r.URL.Query().Get("s") == record.Short {
			http.Redirect(w, r, record.Url, 301)
			return
		}
	}

	fmt.Fprintf(w, "404 Not Found\n")
}

func allHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func main() {
	host := flag.String("host", "localhost:8080", "Where to host the server")
	flag.Parse()

	http.HandleFunc("/s", shortenHandler)
	http.HandleFunc("/s/", shortenHandler)

	http.HandleFunc("/r", redirectHandler)
	http.HandleFunc("/r/", redirectHandler)

	http.HandleFunc("/", allHandler)

	log.Fatal(http.ListenAndServe(*host, nil))
}
