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
		expected func(a args) schema.PutObjectResponse
		wantCode int
	}{
		{
			name: "happy path",
			args: args{
				id:   "SomeID123",
				body: strings.NewReader("object body"),
			},
			prepFunc: func(s *mocks.StorageService, a args) {
				s.On("PutObject", mock.Anything, &domain.Object{
					ID:      a.id,
					Content: bytes.NewBuffer([]byte("object body")),
					Size:    11,
				}).Return(nil)
			},
			expected: func(a args) schema.PutObjectResponse {
				return schema.PutObjectResponse{ID: a.id}
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "invalid object id",
			args: args{
				id:   "",
				body: strings.NewReader("object body"),
			},
			prepFunc: func(s *mocks.StorageService, a args) {
				s.On("PutObject", mock.Anything, mock.Anything).Return(domain.ErrInvalidID)
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "no body",
			args: args{
				id:   "SomeID123",
				body: nil,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			args: args{
				id:   "SomeID123",
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
			var wantResp schema.PutObjectResponse
			err := json.NewDecoder(res.Body).Decode(&wantResp)
			require.NoError(t, err)
			if tt.expected != nil {
				assert.Equal(t, tt.expected(tt.args), wantResp)
			}

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
				id: "SomeID123",
			},
			prepFunc: func(s *mocks.StorageService, a args) {
				s.On("GetObject", mock.Anything, a.id).
					Return(&domain.Object{
						ID:      a.id,
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
				id: "!!!",
			},
			prepFunc: func(s *mocks.StorageService, a args) {
				s.On("GetObject", mock.Anything, mock.Anything).
					Return(nil, domain.ErrInvalidID)
			},
			want: resp{
				body: mustMarshal(t, schema.ServerError{Error: schema.ErrInvalidID}),
				code: http.StatusBadRequest,
			},
		},
		{
			name: "not found",
			args: args{
				id: "SomeID123",
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
				id: "SomeID123",
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
