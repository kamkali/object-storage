package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/server/schema"
)

func (s *Server) putObjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		objectID := mux.Vars(r)["id"]

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r.Body); err != nil || r.Body == http.NoBody {
			s.writeErrResponse(w, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		if err := s.storageService.PutObject(ctx, &domain.Object{
			ID:          objectID,
			Content:     &buf,
			ContentType: r.Header.Get("Content-Type"),
			Size:        r.ContentLength,
		}); err != nil {
			if errors.Is(err, domain.ErrInvalidID) {
				s.writeErrResponse(w, http.StatusBadRequest, schema.ErrInvalidID)
				return
			}
			s.writeErrResponse(w, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(schema.PutObjectResponse{ID: objectID}); err != nil {
			log.Println(fmt.Errorf("cannot write object response: %w", err))
			return
		}
	}
}

func (s *Server) getObjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		objectID := mux.Vars(r)["id"]

		object, err := s.storageService.GetObject(ctx, objectID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrInvalidID):
				s.writeErrResponse(w, http.StatusBadRequest, schema.ErrInvalidID)
				return
			case errors.Is(err, domain.ErrObjNotFound):
				s.writeErrResponse(w, http.StatusNotFound, schema.ErrNotFound)
				return
			default:
				s.writeErrResponse(w, http.StatusInternalServerError, schema.ErrInternal)
				return
			}
		}

		w.Header().Set("Content-Type", object.ContentType)
		w.WriteHeader(http.StatusOK)
		if _, err := io.Copy(w, object.Content); err != nil {
			log.Println(fmt.Errorf("cannot write object response: %w", err))
			return
		}
	}
}
