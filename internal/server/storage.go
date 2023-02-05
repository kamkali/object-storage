package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/server/schema"
)

func (s *Server) putObjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)

		objectID, err := uuid.Parse(vars["id"])
		if err != nil {
			s.writeErrResponse(w, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r.Body); err != nil || r.Body == http.NoBody {
			s.writeErrResponse(w, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		if err := s.storageService.PutObject(ctx, &domain.Object{
			ID:      objectID,
			Content: &buf,
			Size:    buf.Len(),
		}); err != nil {
			s.writeErrResponse(w, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) getObjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)

		objectID, err := uuid.Parse(vars["id"])
		if err != nil {
			s.writeErrResponse(w, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		object, err := s.storageService.GetObject(ctx, objectID)
		if err != nil {
			if errors.Is(err, domain.ErrObjNotFound) {
				s.writeErrResponse(w, http.StatusNotFound, schema.ErrNotFound)
				return
			}
			s.writeErrResponse(w, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := io.Copy(w, object.Content); err != nil {
			log.Println(fmt.Errorf("cannot write object response: %w", err))
			return
		}
	}
}
