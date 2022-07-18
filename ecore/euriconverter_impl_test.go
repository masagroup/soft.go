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
		uri, _ := ParseURI("test:/file.ext")
		assert.Nil(t, c.GetURIHandler(uri))
	}

	{
		uri, _ := ParseURI("file:/file.ext")
		assert.NotNil(t, c.GetURIHandler(uri))
	}
}

func TestEURIConverterCreateReader(t *testing.T) {
	mockHandler := &MockEURIHandler{}
	c := NewEURIConverterImpl()
	c.uriHandlers = NewImmutableEList([]interface{}{mockHandler})

	uri, _ := ParseURI("test:/file.ext")
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
	c.uriMap = map[URI]URI{*NewURI("test:"): *NewURI("file:")}
	uri, _ := ParseURI("test:/file.ext")
	normalized, _ := ParseURI("file:/file.ext")
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

	uri, _ := ParseURI("test:/file.ext")
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
		c.GetURIMap()[*NewURIBuilder(nil).SetScheme("test").URI()] = *NewURIBuilder(nil).SetScheme("test2").URI()
		assert.Equal(t, NewURIBuilder(nil).SetScheme("test2").SetPath("file.ext").URI(), c.Normalize(NewURIBuilder(nil).SetScheme("test").SetPath("file.ext").URI()))
		assert.Equal(t, NewURIBuilder(nil).SetScheme("bla").SetPath("file.ext").URI(), c.Normalize(NewURIBuilder(nil).SetScheme("bla").SetPath("file.ext").URI()))
	}
	{
		c := NewEURIConverterImpl()
		c.GetURIMap()[*NewURIBuilder(nil).SetScheme("test").URI()] = *NewURIBuilder(nil).SetScheme("test2").URI()
		c.GetURIMap()[*NewURIBuilder(nil).SetScheme("test2").URI()] = *NewURIBuilder(nil).SetScheme("test3").URI()
		assert.Equal(t, NewURIBuilder(nil).SetScheme("test3").SetPath("file.ext").URI(), c.Normalize(NewURIBuilder(nil).SetScheme("test").SetPath("file.ext").URI()))
	}
}

func TestEURIConverterURIMap(t *testing.T) {
	c := NewEURIConverterImpl()
	assert.NotNil(t, c)
}
