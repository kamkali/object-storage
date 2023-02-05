package server

import (
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/kamkalis/object-storage/internal/domain"
    "github.com/kamkalis/object-storage/internal/server/schema"
    "io"
    "log"
    "net/http"
)

func (s *Server) putObjectHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        vars := mux.Vars(r)

        objectID, err := uuid.Parse(vars["id"])
        if err != nil {
            s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
            return
        }

        body, err := io.ReadAll(r.Body)
        if err != nil || r.Body == http.NoBody {
            s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
            return
        }

        if err := s.storageService.PutObject(ctx, &domain.Object{
            ID:      objectID,
            Content: body,
        }); err != nil {
            s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
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
            s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
            return
        }

        object, err := s.storageService.GetObject(ctx, objectID)
        if err != nil {
            if errors.Is(err, domain.ErrObjNotFound) {
                s.writeErrResponse(w, err, http.StatusNotFound, schema.ErrNotFound)
                return
            }
            s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
            return
        }

        w.WriteHeader(http.StatusOK)
        _, err = w.Write(object.Content)
        if err != nil {
            log.Println(fmt.Errorf("cannot write object response: %w", err))
            return
        }
    }
}
