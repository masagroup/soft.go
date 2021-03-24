// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

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
