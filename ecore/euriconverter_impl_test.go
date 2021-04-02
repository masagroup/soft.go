package ecore

import (
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
		uri, _ := ParseURI("test://file.ext")
		assert.Nil(t, c.GetURIHandler(uri))
	}

	{
		uri, _ := ParseURI("file://file.ext")
		assert.NotNil(t, c.GetURIHandler(uri))
	}
}

func TestEURIConverterCreateReader(t *testing.T) {
	mockHandler := &MockEURIHandler{}
	c := NewEURIConverterImpl()
	c.uriHandlers = NewImmutableEList([]interface{}{mockHandler})

	uri, _ := ParseURI("test://file.ext")
	mockFile, _ := os.Open(uri.String())
	mockHandler.On("CanHandle", uri).Return(true).Once()
	mockHandler.On("CreateReader", uri).Return(mockFile, nil).Once()
	{
		r, err := c.CreateReader(uri)
		assert.Nil(t, err)
		assert.Equal(t, mockFile, r)
	}
	mock.AssertExpectationsForObjects(t, mockHandler)

	mockHandler.On("CanHandle", uri).Return(false).Once()
	{
		r, err := c.CreateReader(uri)
		assert.NotNil(t, err)
		assert.Equal(t, nil, r)
	}
	mock.AssertExpectationsForObjects(t, mockHandler)
}

func TestEURIConverterCreateReaderWithMapping(t *testing.T) {
	mockHandler := &MockEURIHandler{}
	c := NewEURIConverterImpl()
	c.uriHandlers = NewImmutableEList([]interface{}{mockHandler})
	c.uriMap = map[URI]URI{{Scheme: "test"}: {Scheme: "file"}}
	uri, _ := ParseURI("test:///file.ext")
	normalized, _ := ParseURI("file:///file.ext")
	mockFile, _ := os.Open(normalized.String())
	mockHandler.On("CanHandle", normalized).Return(true).Once()
	mockHandler.On("CreateReader", normalized).Return(mockFile, nil).Once()
	{
		r, err := c.CreateReader(uri)
		assert.Nil(t, err)
		assert.Equal(t, mockFile, r)
	}
	mock.AssertExpectationsForObjects(t, mockHandler)

	mockHandler.On("CanHandle", normalized).Return(false).Once()
	{
		r, err := c.CreateReader(uri)
		assert.NotNil(t, err)
		assert.Equal(t, nil, r)
	}
	mock.AssertExpectationsForObjects(t, mockHandler)
}

func TestEURIConverterCreateWriter(t *testing.T) {
	mockHandler := &MockEURIHandler{}
	c := NewEURIConverterImpl()
	c.uriHandlers = NewImmutableEList([]interface{}{mockHandler})

	uri, _ := ParseURI("test://file.ext")
	mockFile, _ := os.Create(uri.String())

	mockHandler.On("CanHandle", uri).Return(true).Once()
	mockHandler.On("CreateWriter", uri).Return(mockFile, nil).Once()
	{
		r, err := c.CreateWriter(uri)
		assert.Nil(t, err)
		assert.Equal(t, mockFile, r)
	}
	mock.AssertExpectationsForObjects(t, mockHandler)

	mockHandler.On("CanHandle", uri).Return(false).Once()
	{
		r, err := c.CreateWriter(uri)
		assert.NotNil(t, err)
		assert.Equal(t, nil, r)
	}
	mock.AssertExpectationsForObjects(t, mockHandler)
}

func TestEURIConverterNormalize(t *testing.T) {
	{
		c := NewEURIConverterImpl()
		uri, _ := ParseURI("test://file.ext")
		assert.Equal(t, uri, c.Normalize(uri))
	}
	{
		c := NewEURIConverterImpl()
		c.GetURIMap()[URI{Scheme: "test"}] = URI{Scheme: "test2"}
		assert.Equal(t, &URI{Scheme: "test2", Path: "file.ext"}, c.Normalize(&URI{Scheme: "test", Path: "file.ext"}))
		assert.Equal(t, &URI{Scheme: "bla", Path: "file.ext"}, c.Normalize(&URI{Scheme: "bla", Path: "file.ext"}))
	}
	{
		c := NewEURIConverterImpl()
		c.GetURIMap()[URI{Scheme: "test"}] = URI{Scheme: "test2"}
		c.GetURIMap()[URI{Scheme: "test2"}] = URI{Scheme: "test3"}
		assert.Equal(t, &URI{Scheme: "test3", Path: "file.ext"}, c.Normalize(&URI{Scheme: "test", Path: "file.ext"}))
	}
}

func TestEURIConverterURIMap(t *testing.T) {
	c := NewEURIConverterImpl()
	assert.NotNil(t, c)
}
