package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	router := NewRouter()
	if err := http.ListenAndServe(":9090", loggingMiddleware(router)); err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		log.Println(r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}
