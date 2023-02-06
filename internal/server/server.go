package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kamkalis/object-storage/internal/domain"

	"github.com/gorilla/mux"
	"github.com/kamkalis/object-storage/internal/config"
	"github.com/kamkalis/object-storage/internal/server/schema"
	"golang.org/x/net/context"
)

type Server struct {
	router     *mux.Router
	config     *config.Config
	httpServer *http.Server

	storageService domain.StorageService
}

func New(cfg *config.Config, storageService domain.StorageService) (*Server, error) {
	r := mux.NewRouter()
	s := &Server{
		router: r,
		config: cfg,
		httpServer: &http.Server{
			Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
			Handler: r,
		},
		storageService: storageService,
	}

	s.registerRoutes()

	return s, nil
}

func (s *Server) registerRoutes() {
	s.router.Handle("/object/{id}",
		s.withTimeout(s.config.Server.Timeout, s.putObjectHandler()),
	).Methods(http.MethodPut)

	s.router.Handle("/object/{id}",
		s.withTimeout(s.config.Server.Timeout, s.getObjectHandler()),
	).Methods(http.MethodGet)
}

func (s *Server) Start() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server Started on host=%s:%s\n", s.config.Server.Host, s.config.Server.Port)

	<-done
	log.Println("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server Exited Properly")
}

func (s *Server) writeErrResponse(w http.ResponseWriter, code int, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonErr, err := json.Marshal(schema.ServerError{Error: desc})
	if err != nil {
		return
	}
	if _, err := w.Write(jsonErr); err != nil {
		log.Println(fmt.Errorf("cannot write error response: %w", err).Error())
		return
	}
}
