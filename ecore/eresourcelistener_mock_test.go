package ecore

import "testing"

// TestMockEResourceListenerAttached tests method Attached
func TestMockEResourceListenerAttached(t *testing.T) {
	r := NewMockEResourceListener(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o)
	r.EXPECT().Attached(o).Return().Run(func(object EObject) { m.Run(object) }).Once()
	r.Attached(o)
}

// TestMockEResourceListenerDetached tests method Detached
func TestMockEResourceListenerDetached(t *testing.T) {
	r := NewMockEResourceListener(t)
	o := NewMockEObject(t)
	m := NewMockRun(t, o)
	r.EXPECT().Detached(o).Return().Run(func(object EObject) { m.Run(object) }).Once()
	r.Detached(o)
}
