package consistent_hash

import (
	"testing"

	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestGetNode(t *testing.T) {
	var (
		ctx = context.Background()
		n1  = mocks.NewStorageNode(t)
		n2  = mocks.NewStorageNode(t)
		n3  = mocks.NewStorageNode(t)
	)

	tests := []struct {
		name     string
		nodes    []domain.StorageNode
		key      string
		prepFunc func()
		expected domain.StorageNode
		wantErr  bool
	}{
		{
			name:  "Single node in the ring",
			nodes: []domain.StorageNode{n1},
			prepFunc: func() {
				n1.On("ID").Return("node-1")
			},
			key:      "key-1",
			expected: n1,
		},
		{
			name:  "Multiple nodes in the ring",
			nodes: []domain.StorageNode{n1, n2, n3},
			prepFunc: func() {
				n1.On("ID").Return("node-3")
				n2.On("ID").Return("node-4")
				n3.On("ID").Return("node-5")
			},
			key:      "key-5",
			expected: n1,
		},
		{
			name:    "Empty ring",
			nodes:   []domain.StorageNode{},
			key:     "key-x",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepFunc != nil {
				tt.prepFunc()
			}

			lb := NewRingLoadBalancer()
			err := lb.Rebalance(ctx, tt.nodes)
			require.NoError(t, err)

			node, err := lb.GetNode(ctx, tt.key)
			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID(), node.ID())
				return
			}
			assert.Error(t, err)
		})
	}
}
