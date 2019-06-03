package main

import (
	"net/http"

	"github.com/golang/glog"
)

func main() {
	router := NewRouter()
	if err := http.ListenAndServe(":9090", loggingMiddleware(router)); err != nil {
		glog.Fatal(err)
	}
}
