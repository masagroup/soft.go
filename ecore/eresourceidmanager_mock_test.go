package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEResourceIDManagerClear(t *testing.T) {
	rm := &MockEResourceIDManager{}
	rm.On("Clear").Once()
	rm.Clear()
	mock.AssertExpectationsForObjects(t, rm)
}

func TestMockEResourceIDManagerRegister(t *testing.T) {
	rm := &MockEResourceIDManager{}
	o := &MockEObject{}
	rm.On("Register", o).Once()
	rm.Register(o)
	mock.AssertExpectationsForObjects(t, rm, o)
}

func TestMockEResourceIDManagerUnRegister(t *testing.T) {
	rm := &MockEResourceIDManager{}
	o := &MockEObject{}
	rm.On("UnRegister", o).Once()
	rm.UnRegister(o)
	mock.AssertExpectationsForObjects(t, rm, o)
}

func TestMockEResourceIDManagerGetID(t *testing.T) {
	rm := &MockEResourceIDManager{}
	o := &MockEObject{}
	rm.On("GetID", o).Return("id1").Once()
	rm.On("GetID", o).Return(func(EObject) string {
		return "id2"
	}).Once()
	assert.Equal(t, "id1", rm.GetID(o))
	assert.Equal(t, "id2", rm.GetID(o))
	mock.AssertExpectationsForObjects(t, rm, o)
}

func TestMockEResourceIDManagerGetEObject(t *testing.T) {
	rm := &MockEResourceIDManager{}
	o := &MockEObject{}
	rm.On("GetEObject", "id1").Return(o).Once()
	rm.On("GetEObject", "id2").Return(func(string) EObject {
		return o
	}).Once()
	assert.Equal(t, o, rm.GetEObject("id1"))
	assert.Equal(t, o, rm.GetEObject("id2"))
	mock.AssertExpectationsForObjects(t, rm, o)
}
