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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdapterAccessors(t *testing.T) {
	adapter := NewAbstractAdapter()
	assert.Equal(t, nil, adapter.GetTarget())

	mockNotifier := new(MockENotifier)
	adapter.SetTarget(mockNotifier)
	assert.Equal(t, mockNotifier, adapter.GetTarget())
}

type targetOnlyAdapter struct{ *AbstractEAdapter }

func (*targetOnlyAdapter) NotifyChanged(ENotification) {}

// One adapter attached to many notifiers at once, as a content adapter is during concurrent loads.
func TestAdapterAttachedConcurrently(t *testing.T) {
	adapter := &targetOnlyAdapter{NewAbstractAdapter()}
	var group sync.WaitGroup
	for range 8 {
		group.Add(1)
		go func() {
			defer group.Done()
			newEObjectImpl().EAdapters().Add(adapter)
			_ = adapter.GetTarget()
		}()
	}
	group.Wait()
	assert.NotNil(t, adapter.GetTarget())
}
