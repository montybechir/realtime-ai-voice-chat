package middleware

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func NewLogger(logLevel string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	// Apply middleware in reverse order
	// Last middleware is executed first
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
