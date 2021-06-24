package ecore

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBinaryCodec_NewEncoder(t *testing.T) {
	c := &BinaryCodec{}
	mockResource := &MockEResource{}
	mockResource.On("GetURI").Return(nil).Once()
	e := c.NewEncoder(mockResource, nil, nil)
	require.NotNil(t, e)
	mock.AssertExpectationsForObjects(t, mockResource)
}

func TestBinaryCodec_NewDecoder(t *testing.T) {
	c := &BinaryCodec{}
	mockResource := &MockEResource{}
	mockResource.On("GetURI").Return(nil).Once()
	e := c.NewDecoder(mockResource, nil, nil)
	require.NotNil(t, e)
	mock.AssertExpectationsForObjects(t, mockResource)
}
