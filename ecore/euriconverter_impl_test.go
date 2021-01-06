package ecore

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEURIConverterGetURIHandlers(t *testing.T) {
	c := NewEURIConverterImpl()
	assert.NotNil(t, c.GetURIHandlers())
}

func TestEURIConverterGetURIHandler(t *testing.T) {
	c := NewEURIConverterImpl()
	assert.Nil(t, c.GetURIHandler(nil))

	{
		uri, _ := url.Parse("test://file.ext")
		assert.Nil(t, c.GetURIHandler(uri))
	}

	{
		uri, _ := url.Parse("file://file.ext")
		assert.NotNil(t, c.GetURIHandler(uri))
	}
}

func TestEURIConverterCreateReader(t *testing.T) {
	mockHandler := &MockEURIHandler{}
	c := NewEURIConverterImpl()
	c.uriHandlers = NewImmutableEList([]interface{}{mockHandler})

	uri, _ := url.Parse("test://file.ext")
	mockFile, _ := os.Open(uri.String())
	mockHandler.On("CanHandle", uri).Return(true).Once()
	mockHandler.On("CreateReader", uri).Return(mockFile).Once()
	assert.Equal(t, mockFile, c.CreateReader(uri))
	mock.AssertExpectationsForObjects(t, mockHandler)

	mockHandler.On("CanHandle", uri).Return(false).Once()
	assert.Equal(t, nil, c.CreateReader(uri))
	mock.AssertExpectationsForObjects(t, mockHandler)
}

func TestEURIConverterCreateWriter(t *testing.T) {
	mockHandler := &MockEURIHandler{}
	c := NewEURIConverterImpl()
	c.uriHandlers = NewImmutableEList([]interface{}{mockHandler})

	uri, _ := url.Parse("test://file.ext")
	mockFile, _ := os.Create(uri.String())

	mockHandler.On("CanHandle", uri).Return(true).Once()
	mockHandler.On("CreateWriter", uri).Return(mockFile).Once()
	assert.Equal(t, mockFile, c.CreateWriter(uri))
	mock.AssertExpectationsForObjects(t, mockHandler)

	mockHandler.On("CanHandle", uri).Return(false).Once()
	assert.Equal(t, nil, c.CreateWriter(uri))
	mock.AssertExpectationsForObjects(t, mockHandler)
}

func TestEURIConverterNormalize(t *testing.T) {
	c := NewEURIConverterImpl()
	uri, _ := url.Parse("test://file.ext")
	assert.Equal(t, uri, c.Normalize(uri))
}
