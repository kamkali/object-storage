// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kamkalis/object-storage/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// StorageNode is an autogenerated mock type for the StorageNode type
type StorageNode struct {
	mock.Mock
}

// GetObject provides a mock function with given fields: ctx, key
func (_m *StorageNode) GetObject(ctx context.Context, key string) (*domain.Object, error) {
	ret := _m.Called(ctx, key)

	var r0 *domain.Object
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Object); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Object)
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

// ID provides a mock function with given fields:
func (_m *StorageNode) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsAlive provides a mock function with given fields: ctx
func (_m *StorageNode) IsAlive(ctx context.Context) bool {
	ret := _m.Called(ctx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// PutObject provides a mock function with given fields: ctx, o
func (_m *StorageNode) PutObject(ctx context.Context, o *domain.Object) error {
	ret := _m.Called(ctx, o)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Object) error); ok {
		r0 = rf(ctx, o)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStorageNode interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorageNode creates a new instance of StorageNode. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorageNode(t mockConstructorTestingTNewStorageNode) *StorageNode {
	mock := &StorageNode{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
