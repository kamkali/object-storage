package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/domain/mocks"
	"github.com/kamkalis/object-storage/internal/server/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPutObjectHandler(t *testing.T) {
	type args struct {
		id   string
		body io.Reader
	}
	testCases := []struct {
		name     string
		args     args
		prepFunc func(s *mocks.StorageService, a args)
		wantCode int
	}{
		{
			name: "happy path",
			args: args{
				id:   uuid.NewString(),
				body: strings.NewReader("object body"),
			},
			prepFunc: func(s *mocks.StorageService, a args) {
				s.On("PutObject", mock.Anything, &domain.Object{
					ID:      uuid.MustParse(a.id),
					Content: bytes.NewBuffer([]byte("object body")),
					Size:    11,
				}).Return(nil)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "invalid object id",
			args: args{
				id: "weird_id",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "no body",
			args: args{
				id:   uuid.NewString(),
				body: nil,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			args: args{
				id:   uuid.NewString(),
				body: strings.NewReader("object body"),
			},
			prepFunc: func(s *mocks.StorageService, a args) {
				s.On("PutObject", mock.Anything, mock.Anything).
					Return(errors.New("something went wrong"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mocks.NewStorageService(t)
			if tt.prepFunc != nil {
				tt.prepFunc(mockStorage, tt.args)
			}

			s := Server{storageService: mockStorage}
			objectID := tt.args.id
			req := httptest.NewRequest("PUT", "/object/"+objectID, tt.args.body)
			req = mux.SetURLVars(req, map[string]string{
				"id": objectID,
			})

			res := httptest.NewRecorder()
			s.putObjectHandler().ServeHTTP(res, req)
			assert.Equal(t, tt.wantCode, res.Code)
		})
	}
}

func TestGetObjectHandler(t *testing.T) {
	type args struct {
		id string
	}
	type resp struct {
		body []byte
		code int
	}
	testCases := []struct {
		name     string
		args     args
		prepFunc func(s *mocks.StorageService, a args)
		want     resp
	}{
		{
			name: "happy path",
			args: args{
				id: uuid.NewString(),
			},
			prepFunc: func(s *mocks.StorageService, a args) {
				s.On("GetObject", mock.Anything, uuid.MustParse(a.id)).
					Return(&domain.Object{
						ID:      uuid.MustParse(a.id),
						Content: bytes.NewBuffer([]byte("object body")),
						Size:    11,
					}, nil)
			},
			want: resp{
				body: []byte("object body"),
				code: http.StatusOK,
			},
		},
		{
			name: "invalid object id",
			args: args{
				id: "weird_id",
			},
			want: resp{
				body: mustMarshal(t, schema.ServerError{Error: schema.ErrBadRequest}),
				code: http.StatusBadRequest,
			},
		},
		{
			name: "not found",
			args: args{
				id: uuid.NewString(),
			},
			prepFunc: func(s *mocks.StorageService, a args) {
				s.On("GetObject", mock.Anything, mock.Anything).
					Return(nil, domain.ErrObjNotFound)
			},
			want: resp{
				body: mustMarshal(t, schema.ServerError{Error: schema.ErrNotFound}),
				code: http.StatusNotFound,
			},
		},
		{
			name: "service error",
			args: args{
				id: uuid.NewString(),
			},
			prepFunc: func(s *mocks.StorageService, a args) {
				s.On("GetObject", mock.Anything, mock.Anything).
					Return(nil, errors.New("something went wrong"))
			},
			want: resp{
				body: mustMarshal(t, schema.ServerError{Error: schema.ErrInternal}),
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mocks.NewStorageService(t)
			if tt.prepFunc != nil {
				tt.prepFunc(mockStorage, tt.args)
			}

			s := Server{storageService: mockStorage}
			objectID := tt.args.id
			req := httptest.NewRequest("GET", "/object/"+objectID, nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": objectID,
			})

			res := httptest.NewRecorder()
			s.getObjectHandler().ServeHTTP(res, req)
			assert.Equal(t, tt.want.code, res.Code)
			assert.Equal(t, tt.want.body, res.Body.Bytes())
		})
	}
}

func mustMarshal(t *testing.T, v any) []byte {
	marshalled, err := json.Marshal(v)
	require.NoError(t, err)
	return marshalled
}
