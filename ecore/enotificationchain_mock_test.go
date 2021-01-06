package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockENotificationChainAdd(t *testing.T) {
	nc := &MockENotificationChain{}
	n := &MockENotification{}
	nc.On("Add", n).Return(true).Once()
	nc.On("Add", n).Return(func(ENotification) bool {
		return false
	}).Once()
	assert.True(t, nc.Add(n))
	assert.False(t, nc.Add(n))
	mock.AssertExpectationsForObjects(t, n, nc)
}

func TestMockENotificationChainDispatch(t *testing.T) {
	nc := &MockENotificationChain{}
	nc.On("Dispatch").Once()
	nc.Dispatch()
	mock.AssertExpectationsForObjects(t, nc)
}
