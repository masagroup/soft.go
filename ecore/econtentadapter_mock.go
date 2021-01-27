package ecore

import (
	mock "github.com/stretchr/testify/mock"
)

type MockEContentAdapter struct {
	mock.Mock
	EContentAdapter
}

func NewMockContentAdapter() *MockEContentAdapter {
	c := new(MockEContentAdapter)
	c.SetInterfaces(c)
	return c
}

func (adapter *MockEContentAdapter) NotifyChanged(notification ENotification) {
	adapter.Called(notification)
	adapter.EContentAdapter.NotifyChanged(notification)
}
