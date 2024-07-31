package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func maxAllowedMiddleware(n uint) mux.MiddlewareFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acquire()
			defer release()
			next.ServeHTTP(w, r)
		})
	}
}

// enableCORS sets the Vary: Origin and Access-Control-Allow-Origin response headers in order to
// enabled CORS for trusted origins.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		println("to day")
		// On run this if there's an Origin request header pres
		// If there is a match, then set an "Access-Control-Allow-Origin" response
		// header with the request origin as the value and break out of the loop.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
