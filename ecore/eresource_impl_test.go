package ecore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResourceURI(t *testing.T) {
	u, err := url.Parse("https://example.com/foo%2fbar")
	assert.Nil(t, err)
	r := NewEResourceImpl()
	r.SetURI(u)
	assert.Equal(t, u, r.GetURI())
}

func TestResourceURINotifications(t *testing.T) {
	r := NewEResourceImpl()
	mockEAdapter := new(MockEAdapter)
	mockEAdapter.On("SetTarget", r).Once()
	r.EAdapters().Add(mockEAdapter)
	mock.AssertExpectationsForObjects(t, mockEAdapter)

	u, err := url.Parse("https://example.com/foo%2fbar")
	assert.Nil(t, err)

	mockEAdapter.On("NotifyChanged", mock.Anything).Once()
	r.SetURI(u)
	assert.Equal(t, u, r.GetURI())
	mock.AssertExpectationsForObjects(t, mockEAdapter)
}

func TestResourceContents(t *testing.T) {
	r := NewEResourceImpl()

	mockEObjectInternal := new(MockEObjectInternal)
	mockEObjectInternal.On("ESetResource", r, mock.Anything).Return(nil)
	r.GetContents().Add(mockEObjectInternal)

	mockEObjectInternal.On("ESetResource", nil, mock.Anything).Return(nil)
	r.GetContents().Remove(mockEObjectInternal)
}

func TestResourceGetEObject(t *testing.T) {
	r := NewEResourceImpl()
	assert.Nil(t, r.GetEObject("Test"))
}

func TestResourceLoadInvalid(t *testing.T) {
	r := NewEResourceImpl()
	r.SetURI(&url.URL{Path: "testdata/invalid.xml"})
	r.Load()
	assert.False(t, r.IsLoaded())
	assert.False(t, r.GetErrors().Empty())
}
