// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kamkalis/object-storage/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// LoadBalancer is an autogenerated mock type for the LoadBalancer type
type LoadBalancer struct {
	mock.Mock
}

// GetNode provides a mock function with given fields: ctx, key
func (_m *LoadBalancer) GetNode(ctx context.Context, key string) (domain.StorageNode, error) {
	ret := _m.Called(ctx, key)

	var r0 domain.StorageNode
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.StorageNode); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.StorageNode)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReBalance provides a mock function with given fields: ctx, nodes
func (_m *LoadBalancer) ReBalance(ctx context.Context, nodes []domain.StorageNode) error {
	ret := _m.Called(ctx, nodes)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []domain.StorageNode) error); ok {
		r0 = rf(ctx, nodes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewLoadBalancer interface {
	mock.TestingT
	Cleanup(func())
}

// NewLoadBalancer creates a new instance of LoadBalancer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLoadBalancer(t mockConstructorTestingTNewLoadBalancer) *LoadBalancer {
	mock := &LoadBalancer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
