package service

import (
	"bytes"
	"errors"
	"testing"

	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var _ domain.StorageService = (*StorageService)(nil)

func TestStorageService_PutObject(t *testing.T) {
	ctx := context.Background()
	type mockDeps struct {
		manager *mocks.NodeManager
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
					ID:      "SomeId123",
					Content: bytes.NewReader([]byte("content")),
				},
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("GetNode", a.ctx, a.o.ID).
					Return(d.node, nil)

				d.node.
					On("IsAlive", a.ctx).Return(true).
					On("PutObject", a.ctx, a.o).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "cannot select node",
			args: args{
				ctx: ctx,
				o: &domain.Object{
					ID:      "SomeId123",
					Content: bytes.NewReader([]byte("content")),
				},
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("GetNode", a.ctx, a.o.ID).
					Return(nil, errors.New("err"))
			},
			wantErr: true,
		},
		{
			name: "ID not alphanum",
			args: args{
				ctx: ctx,
				o: &domain.Object{
					ID:      "WhatNot!!!--",
					Content: bytes.NewReader([]byte("content")),
				},
			},
			wantErr: true,
		},
		{
			name: "node offline",
			args: args{
				ctx: ctx,
				o: &domain.Object{
					ID:      "SomeId123",
					Content: bytes.NewReader([]byte("content")),
				},
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("GetNode", a.ctx, a.o.ID).Return(d.node, nil)
				d.node.
					On("ID").Return("node-1").
					On("IsAlive", a.ctx).Return(false)
			},
			wantErr: true,
		},
		{
			name: "cannot put to selected node",
			args: args{
				ctx: ctx,
				o: &domain.Object{
					ID:      "SomeId123",
					Content: bytes.NewReader([]byte("content")),
				},
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("GetNode", a.ctx, a.o.ID).
					Return(d.node, nil)

				d.node.
					On("IsAlive", a.ctx).Return(true).
					On("PutObject", a.ctx, a.o).Return(errors.New("err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			deps := &mockDeps{
				manager: mocks.NewNodeManager(t),
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
		ID:      "SomeId123",
		Content: bytes.NewReader([]byte("content")),
	}
	type mockDeps struct {
		manager *mocks.NodeManager
		node    *mocks.StorageNode
	}
	type args struct {
		ctx context.Context
		id  string
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
				id:  "SomeId123",
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("GetNode", a.ctx, a.id).
					Return(d.node, nil)

				d.node.
					On("IsAlive", a.ctx).Return(true).
					On("GetObject", a.ctx, a.id).Return(expected, nil)
			},
			expected: expected,
			wantErr:  false,
		},
		{
			name: "cannot select node",
			args: args{
				ctx: ctx,
				id:  "SomeId123",
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("GetNode", a.ctx, a.id).
					Return(nil, errors.New("err"))
			},
			wantErr: true,
		},
		{
			name: "ID empty",
			args: args{
				ctx: ctx,
				id:  "WhatNot!!!--",
			},
			wantErr: true,
		},
		{
			name: "node offline",
			args: args{
				ctx: ctx,
				id:  "SomeId123",
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("GetNode", a.ctx, a.id).Return(d.node, nil)
				d.node.
					On("ID").Return("node-1").
					On("IsAlive", a.ctx).Return(false)
			},
			wantErr: true,
		},
		{
			name: "cannot put to selected node",
			args: args{
				ctx: ctx,
				id:  "SomeId123",
			},
			prepFunc: func(d *mockDeps, a args) {
				d.manager.
					On("GetNode", a.ctx, a.id).
					Return(d.node, nil)

				d.node.
					On("IsAlive", a.ctx).Return(true).
					On("GetObject", a.ctx, a.id).Return(nil, errors.New("err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			deps := &mockDeps{
				manager: mocks.NewNodeManager(t),
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
