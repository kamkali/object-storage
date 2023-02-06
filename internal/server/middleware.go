package server

import (
	"net/http"
	"time"

	"github.com/kamkalis/object-storage/internal/server/schema"
)

func (s *Server) withTimeout(timeout time.Duration, next http.HandlerFunc) http.Handler {
	return http.TimeoutHandler(next, timeout, schema.ErrTimedOut)
}
