package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEFactoryExtCreate(t *testing.T) {
	f := NewEFactoryExt()
	mockClass := &MockEClass{}
	mockPackage := &MockEPackage{}
	mockClass.On("GetEPackage").Return(mockPackage).Once()
	assert.Panics(t, func() { f.Create(mockClass) })
	mock.AssertExpectationsForObjects(t, mockClass, mockPackage)
}
