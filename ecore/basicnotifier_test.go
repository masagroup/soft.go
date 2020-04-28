package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotifierConstructor(t *testing.T) {
	assert.NotNil(t, NewBasicNotifier())
}

func TestNotifierAccessors(t *testing.T) {
	var n ENotifier = NewBasicNotifier()
	assert.True(t, n.EDeliver())

	n.ESetDeliver(false)
	assert.False(t, n.EDeliver())

	adapters := n.EAdapters()
	assert.True(t, adapters.Empty())
}

func TestNotifierAdapters(t *testing.T) {
	n := NewBasicNotifier()
	mockEAdapter := new(MockEAdapter)
	mockEAdapter.On("SetTarget", n).Once()
	n.EAdapters().Add(mockEAdapter)
	mockEAdapter.AssertExpectations(t)
}
