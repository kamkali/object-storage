// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StorageManager is an autogenerated mock type for the StorageManager type
type StorageManager struct {
	mock.Mock
}

type mockConstructorTestingTNewStorageManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorageManager creates a new instance of StorageManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorageManager(t mockConstructorTestingTNewStorageManager) *StorageManager {
	mock := &StorageManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}