package ecore

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestEDataTypeExtSetDefaultValue(t *testing.T) {
	d := newEDataTypeExt()
	mockValue := &MockEObject{}
	mockAdapter := &MockEAdapter{}
	mockAdapter.On("SetTarget", d).Once()
	d.EAdapters().Add(mockAdapter)
	mock.AssertExpectationsForObjects(t, mockAdapter)

	mockAdapter.On("NotifyChanged", mock.MatchedBy(func(notification ENotification) bool {
		return notification.GetEventType() == SET && notification.GetNewValue() == mockValue
	})).Once()
	d.SetDefaultValue(mockValue)
	mock.AssertExpectationsForObjects(t, mockAdapter)
}
