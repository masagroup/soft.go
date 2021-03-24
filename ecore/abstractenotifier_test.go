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

	"github.com/stretchr/testify/assert"
)

func TestAbstractNotifier(t *testing.T) {
	notifier := &AbstractENotifier{}
	notifier.SetInterfaces(notifier)
	assert.Equal(t, notifier, notifier.GetInterfaces())
	assert.Equal(t, notifier, notifier.AsENotifier())
	assert.Equal(t, notifier, notifier.AsENotifierInternal())
	assert.Nil(t, notifier.EBasicAdapters())
	assert.False(t, notifier.EBasicHasAdapters())
	assert.NotNil(t, notifier.EAdapters())
	assert.True(t, notifier.EAdapters().Empty())
	assert.False(t, notifier.EDeliver())
	assert.Panics(t, func() {
		notifier.ESetDeliver(false)
	})
}

func TestNotifierNotification(t *testing.T) {
	notifier := &AbstractENotifier{}
	notifier.SetInterfaces(notifier)
	notification := &notifierNotification{notifier: notifier}

	assert.Equal(t, notifier, notification.GetNotifier())
	assert.Equal(t, nil, notification.GetFeature())
	assert.Equal(t, -1, notification.GetFeatureID())
}
