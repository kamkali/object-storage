package server

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func (s *Server) withTimeout(timeout uint, next http.HandlerFunc) http.HandlerFunc {
	// TODO: maybe just use http.TimeoutHandler
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
