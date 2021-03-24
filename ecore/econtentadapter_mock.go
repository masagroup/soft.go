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
