package middleware

import (
	"log"
	"net/http"
)

func AuthMiddleware(next HandleFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s\n", r.RequestURI)
		next(w, r)
	}
}
