package middleware

import "net/http"

type HandleFunc func(w http.ResponseWriter, r *http.Request)

type Middleware func(HandleFunc) HandleFunc

func Handle(finalHandler HandleFunc, middlewares ...Middleware) HandleFunc {
	if finalHandler == nil {
		panic("No final handler")
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		finalHandler = middlewares[i](finalHandler)
	}
	return finalHandler
}
