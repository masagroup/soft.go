package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEObjectIDManagerClear(t *testing.T) {
	rm := &MockEObjectIDManager{}
	rm.On("Clear").Once()
	rm.Clear()
	mock.AssertExpectationsForObjects(t, rm)
}

func TestMockEObjectIDManagerRegister(t *testing.T) {
	rm := &MockEObjectIDManager{}
	o := &MockEObject{}
	rm.On("Register", o).Once()
	rm.Register(o)
	mock.AssertExpectationsForObjects(t, rm, o)
}

func TestMockEObjectIDManagerUnRegister(t *testing.T) {
	rm := &MockEObjectIDManager{}
	o := &MockEObject{}
	rm.On("UnRegister", o).Once()
	rm.UnRegister(o)
	mock.AssertExpectationsForObjects(t, rm, o)
}

func TestMockEObjectIDManagerGetID(t *testing.T) {
	rm := &MockEObjectIDManager{}
	o := &MockEObject{}
	rm.On("GetID", o).Return("id1").Once()
	rm.On("GetID", o).Return(func(EObject) interface{} {
		return "id2"
	}).Once()
	assert.Equal(t, "id1", rm.GetID(o))
	assert.Equal(t, "id2", rm.GetID(o))
	mock.AssertExpectationsForObjects(t, rm, o)
}

func TestMockEObjectIDManagerGetDetachedID(t *testing.T) {
	rm := &MockEObjectIDManager{}
	o := &MockEObject{}
	rm.On("GetDetachedID", o).Return("id1").Once()
	rm.On("GetDetachedID", o).Return(func(EObject) interface{} {
		return "id2"
	}).Once()
	assert.Equal(t, "id1", rm.GetDetachedID(o))
	assert.Equal(t, "id2", rm.GetDetachedID(o))
	mock.AssertExpectationsForObjects(t, rm, o)
}

func TestMockEObjectIDManagerGetEObject(t *testing.T) {
	rm := &MockEObjectIDManager{}
	o := &MockEObject{}
	rm.On("GetEObject", "id1").Return(o).Once()
	rm.On("GetEObject", "id2").Return(func(interface{}) EObject {
		return o
	}).Once()
	assert.Equal(t, o, rm.GetEObject("id1"))
	assert.Equal(t, o, rm.GetEObject("id2"))
	mock.AssertExpectationsForObjects(t, rm, o)
}
