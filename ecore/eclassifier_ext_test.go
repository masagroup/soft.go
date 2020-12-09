package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEClassifierExtClassifierID(t *testing.T) {
	c := newEClassifierExt()
	assert.Equal(t, -1, c.GetClassifierID())

	mockPackage := &MockEPackage{}
	mockClassifiers := &MockEList{}
	c.ESetInternalContainer(mockPackage, ECLASSIFIER__EPACKAGE)
	mockPackage.On("GetEClassifiers").Return(mockClassifiers).Once()
	mockPackage.On("EIsProxy").Return(false).Once()
	mockClassifiers.On("IndexOf", c).Return(0).Once()
	assert.Equal(t, 0, c.GetClassifierID())
	mock.AssertExpectationsForObjects(t, mockPackage, mockClassifiers)
}

func TestEClassifierExtGetDefaultValue(t *testing.T) {
	c := newEClassifierExt()
	assert.Nil(t, c.GetDefaultValue())
}
