package ecore

import "testing"

func BenchmarkIncrementalIDManager_Register(b *testing.B) {
	m := NewIncrementalIDManager()
	for i := 0; i < b.N; i++ {
		o := &MockEObject{}
		m.Register(o)
	}
}

func BenchmarkIncrementalIDManager_UnRegister(b *testing.B) {

}

func BenchmarkIncrementalIDManager_GetID(b *testing.B) {
	m := NewIncrementalIDManager()
	o := &MockEObject{}
	m.Register(o)
	for i := 0; i < b.N; i++ {
		m.GetID(o)
	}
}

func BenchmarkUniqueIDManager_Register(b *testing.B) {
	m := NewUniqueIDManager(20)
	for i := 0; i < b.N; i++ {
		o := &MockEObject{}
		m.Register(o)
	}
}

func BenchmarkUniqueIDManager_GetID(b *testing.B) {
	m := NewUniqueIDManager(20)
	o := &MockEObject{}
	m.Register(o)
	for i := 0; i < b.N; i++ {
		m.GetID(o)
	}
}

func BenchmarkULIDManager_Register(b *testing.B) {
	m := NewULIDManager()
	for i := 0; i < b.N; i++ {
		o := &MockEObject{}
		m.Register(o)
	}
}

func BenchmarkULIDManager_GetID(b *testing.B) {
	m := NewULIDManager()
	o := &MockEObject{}
	m.Register(o)
	for i := 0; i < b.N; i++ {
		m.GetID(o)
	}
}
