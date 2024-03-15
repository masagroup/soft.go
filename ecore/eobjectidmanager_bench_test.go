// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "testing"

func BenchmarkIncrementalIDManager_Register(b *testing.B) {
	m := NewIncrementalIDManager()
	for i := 0; i < b.N; i++ {
		o := NewMockEObject(b)
		m.Register(o)
	}
}

func BenchmarkIncrementalIDManager_UnRegister(b *testing.B) {

}

func BenchmarkIncrementalIDManager_GetID(b *testing.B) {
	m := NewIncrementalIDManager()
	o := NewMockEObject(b)
	m.Register(o)
	for i := 0; i < b.N; i++ {
		m.GetID(o)
	}
}

func BenchmarkUniqueIDManager_Register(b *testing.B) {
	m := NewUUIDManager(20)
	for i := 0; i < b.N; i++ {
		o := NewMockEObject(b)
		m.Register(o)
	}
}

func BenchmarkUniqueIDManager_GetID(b *testing.B) {
	m := NewUUIDManager(20)
	o := NewMockEObject(b)
	m.Register(o)
	for i := 0; i < b.N; i++ {
		m.GetID(o)
	}
}
