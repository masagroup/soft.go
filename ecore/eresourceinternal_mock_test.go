package ecore

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEResourceInternalDoLoad(t *testing.T) {
	r := &MockEResourceInternal{}
	uri, _ := url.Parse("test://file.t")
	f, _ := os.Open(uri.String())
	r.On("DoLoad", f).Once()
	r.DoLoad(f)
	mock.AssertExpectationsForObjects(t, r)
}

func TestMockEResourceInternalDoSave(t *testing.T) {
	r := &MockEResourceInternal{}
	uri, _ := url.Parse("test://file.t")
	f, _ := os.Create(uri.String())
	r.On("DoSave", f).Once()
	r.DoSave(f)
	mock.AssertExpectationsForObjects(t, r)
}

func TestMockEResourceInternalDoUnLoad(t *testing.T) {
	r := &MockEResourceInternal{}
	r.On("DoUnload").Once()
	r.DoUnload()
	mock.AssertExpectationsForObjects(t, r)
}

func TestMockEResourceInternalBasicSetLoaded(t *testing.T) {
	r := &MockEResourceInternal{}
	n1 := &MockENotificationChain{}
	n2 := &MockENotificationChain{}
	r.On("basicSetLoaded", false, n1).Return(n2).Once()
	r.On("basicSetLoaded", false, n1).Return(func(bool, ENotificationChain) ENotificationChain {
		return n2
	}).Once()
	assert.Equal(t, n2, r.basicSetLoaded(false, n1))
	assert.Equal(t, n2, r.basicSetLoaded(false, n1))
	mock.AssertExpectationsForObjects(t, r, n1, n2)
}

func TestMockEResourceInternalBasicSetResourceSet(t *testing.T) {
	r := &MockEResourceInternal{}
	rs := &MockEResourceSet{}
	n1 := &MockENotificationChain{}
	n2 := &MockENotificationChain{}
	r.On("basicSetResourceSet", rs, n1).Return(n2).Once()
	r.On("basicSetResourceSet", rs, n1).Return(func(EResourceSet, ENotificationChain) ENotificationChain {
		return n2
	}).Once()
	assert.Equal(t, n2, r.basicSetResourceSet(rs, n1))
	assert.Equal(t, n2, r.basicSetResourceSet(rs, n1))
	mock.AssertExpectationsForObjects(t, r, rs, n1, n2)
}
