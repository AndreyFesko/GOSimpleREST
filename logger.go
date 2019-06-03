package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		glog.V(2).Info(r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}
