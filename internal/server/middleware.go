package server

import (
	"net/http"
	"time"

	"github.com/kamkalis/object-storage/internal/server/schema"
)

func (s *Server) withTimeout(timeout uint, next http.HandlerFunc) http.Handler {
	return http.TimeoutHandler(next, time.Duration(timeout)*time.Second, schema.ErrTimedOut)
}
