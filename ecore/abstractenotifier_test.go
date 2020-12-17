package ecore

import (
	"testing"
)

func TestNotifierAdapter(t *testing.T) {
	// notifier := NewBasicNotifier()
	// mockEAdapter := new(MockEAdapter)
	// mockEAdapter.On("SetTarget", notifier).Once()
	// notifier.EAdapters().Add(mockEAdapter)

	// mockEAdapter.On("NotifyChanged", mock.MatchedBy(func(n ENotification) bool {
	// 	return n.GetNotifier() == notifier &&
	// 		n.GetFeatureID() == -1 &&
	// 		n.GetFeature() == nil &&
	// 		n.GetNewValue() == nil &&
	// 		n.GetOldValue() == mockEAdapter &&
	// 		n.GetEventType() == REMOVING_ADAPTER &&
	// 		n.GetPosition() == 0
	// })).Once()
	// mockEAdapter.On("UnSetTarget", notifier).Once()
	// notifier.EAdapters().Remove(mockEAdapter)
	// mockEAdapter.AssertExpectations(t)
}
