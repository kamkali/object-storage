package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var _ domain.StorageService = (*StorageService)(nil)

func TestStorageService_PutObject(t *testing.T) {
	ctx := context.Background()
	type mockDeps struct {
		manager *mocks.StorageManager
		node    *mocks.StorageNode
	}
	type args struct {
		ctx context.Context
		o   *domain.Object
	}
	testCases := []struct {
		name     string
		args     args
		prepFunc func(d *mockDeps, a args)
		wantErr  bool
	}{
		{
			name: "puts object successfully",
			args: args{
				ctx: ctx,
				o: &domain.Object{
					ID:      uuid.New(),
					Content: []byte("content"),
				},
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("SelectNode", a.ctx, a.o.ID).
					Return(d.node, nil)

				d.node.
					On("PutObject", a.ctx, a.o).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "cannot select node",
			args: args{
				ctx: ctx,
				o: &domain.Object{
					ID:      uuid.New(),
					Content: []byte("content"),
				},
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("SelectNode", a.ctx, a.o.ID).
					Return(nil, errors.New("err"))
			},
			wantErr: true,
		},
		{
			name: "cannot put to selected node",
			args: args{
				ctx: ctx,
				o: &domain.Object{
					ID:      uuid.New(),
					Content: []byte("content"),
				},
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("SelectNode", a.ctx, a.o.ID).
					Return(d.node, nil)

				d.node.
					On("PutObject", a.ctx, a.o).
					Return(errors.New("err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			deps := &mockDeps{
				manager: mocks.NewStorageManager(t),
				node:    mocks.NewStorageNode(t),
			}
			if tt.prepFunc != nil {
				tt.prepFunc(deps, tt.args)
			}

			s := NewStorage(deps.manager)
			if err := s.PutObject(tt.args.ctx, tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("PutObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorageService_GetObject(t *testing.T) {
	ctx := context.Background()
	expected := &domain.Object{
		ID:      uuid.New(),
		Content: []byte("object"),
	}
	type mockDeps struct {
		manager *mocks.StorageManager
		node    *mocks.StorageNode
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	testCases := []struct {
		name     string
		args     args
		prepFunc func(d *mockDeps, a args)
		expected *domain.Object
		wantErr  bool
	}{
		{
			name: "gets object successfully",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("SelectNode", a.ctx, a.id).
					Return(d.node, nil)

				d.node.
					On("GetObject", a.ctx, a.id).
					Return(expected, nil)
			},
			expected: expected,
			wantErr:  false,
		},
		{
			name: "cannot select node",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("SelectNode", a.ctx, a.id).
					Return(nil, errors.New("err"))
			},
			wantErr: true,
		},
		{
			name: "cannot put to selected node",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("SelectNode", a.ctx, a.id).
					Return(d.node, nil)

				d.node.
					On("GetObject", a.ctx, a.id).
					Return(nil, errors.New("err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			deps := &mockDeps{
				manager: mocks.NewStorageManager(t),
				node:    mocks.NewStorageNode(t),
			}
			if tt.prepFunc != nil {
				tt.prepFunc(deps, tt.args)
			}

			s := NewStorage(deps.manager)
			got, err := s.GetObject(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}
