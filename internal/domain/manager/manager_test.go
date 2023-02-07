package manager

import (
	"errors"
	"testing"

	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestStorageManager_RefreshNodes(t *testing.T) {
	ctx := context.Background()
	lb := mocks.NewLoadBalancer(t)
	dv := mocks.NewNodeDiscoverer(t)
	s := NewStorageManager(lb, dv)

	nodes := []domain.StorageNode{
		&mocks.StorageNode{},
		&mocks.StorageNode{},
	}
	t.Run("happy path", func(t *testing.T) {
		dv.On("DiscoverNodes", ctx).Return(nodes, nil).Once()
		lb.On("ReBalance", ctx, nodes).Return(nil).Once()
		err := s.RefreshNodes(ctx)
		assert.NoError(t, err)
	})
	t.Run("discoverer err wrapped", func(t *testing.T) {
		mockErr := errors.New("mockerr")
		dv.On("DiscoverNodes", ctx).Return(nil, mockErr).Once()
		err := s.RefreshNodes(ctx)
		assert.ErrorIs(t, err, mockErr)
	})
	t.Run("load balancer err wrapped", func(t *testing.T) {
		mockErr := errors.New("mockerr")
		dv.On("DiscoverNodes", ctx).Return(nodes, nil).Once()
		lb.On("ReBalance", ctx, nodes).Return(mockErr).Once()

		err := s.RefreshNodes(ctx)
		assert.ErrorIs(t, err, mockErr)
	})
}
