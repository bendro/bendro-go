package util

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)

		start := time.Now()

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set(
				"Access-Control-Allow-Origin",
				origin,
			)

			w.Header().Set(
				"Access-Control-Allow-Methods",
				"OPTIONS, GET, POST, PUT, DELETE",
			)

			w.Header().Set(
				"Access-Control-Allow-Headers",
				"content-type",
			)
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
		} else {
			h.ServeHTTP(w, r)
		}

		log.Printf(
			"%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
