package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestNotifierWithAdapter(t *testing.T) {
	notifier := NewBasicNotifier()
	mockEAdapter := new(MockEAdapter)
	mockEAdapter.On("SetTarget", notifier).Once()
	notifier.EAdapters().Add(mockEAdapter)

	mockEAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == notifier &&
			n.GetFeatureID() == -1 &&
			n.GetNewValue() == nil &&
			n.GetOldValue() == mockEAdapter &&
			n.GetEventType() == REMOVING_ADAPTER &&
			n.GetPosition() == 0
	})).Once()
	mockEAdapter.On("GetTarget").Return(notifier).Once()
	mockEAdapter.On("SetTarget", nil).Once()
	notifier.EAdapters().Remove(mockEAdapter)
	mockEAdapter.AssertExpectations(t)
}
