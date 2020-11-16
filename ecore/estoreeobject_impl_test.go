package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type EStoreEObjectTest struct {
	*EStoreEObjectImpl
	mockStore *MockEStore
}

func newEStoreEObjectTest(isCaching bool) *EStoreEObjectTest {
	o := new(EStoreEObjectTest)
	o.mockStore = new(MockEStore)
	o.EStoreEObjectImpl = NewEStoreEObjectImpl(isCaching)
	o.SetInterfaces(o)
	return o
}

func (o *EStoreEObjectTest) EStore() EStore {
	return o.mockStore
}

func TestEStoreEObjectImpl_GetAttribute_Caching(t *testing.T) {
	// create object
	o := newEStoreEObjectTest(true)

	// create mocks
	mockClass := new(MockEClass)
	mockAttribute := new(MockEAttribute)
	mockStore := o.mockStore

	// initialise object with mock class
	o.setEClass(mockClass)

	// first get
	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("IsTransient").Return(false).Once()
	mockClass.On("GetFeatureCount").Return(1).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Twice()
	mockStore.On("Get", o, mockAttribute, NO_INDEX).Return(2).Once()
	assert.Equal(t, 2, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	// second - test caching
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Once()
	assert.Equal(t, 2, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_NoCaching(t *testing.T) {
	// create object
	o := newEStoreEObjectTest(false)

	// create mocks
	mockClass := new(MockEClass)
	mockAttribute := new(MockEAttribute)
	mockStore := o.mockStore

	// initialise object with mock class
	o.setEClass(mockClass)

	mockClass.On("GetFeatureCount").Return(1).Once()
	for i := 0; i < 2; i++ {
		mockAttribute.On("IsMany").Return(false).Once()
		mockAttribute.On("IsTransient").Return(false).Once()
		mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Twice()
		mockStore.On("Get", o, mockAttribute, NO_INDEX).Return(2).Once()
		assert.Equal(t, 2, o.EGetFromID(0, false))
		mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
	}
}

func TestEStoreEObjectImpl_SetAttribute(t *testing.T) {

	// create object
	o := newEStoreEObjectTest(false)

	// create mocks
	mockClass := new(MockEClass)
	mockAttribute := new(MockEAttribute)
	mockStore := o.mockStore

	// initialise object with mock class
	o.setEClass(mockClass)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("IsTransient").Return(false).Twice()
	mockClass.On("GetFeatureCount").Return(1).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Times(3)
	mockStore.On("Get", o, mockAttribute, NO_INDEX).Return(nil).Once()
	mockStore.On("Set", o, mockAttribute, NO_INDEX, 2).Return(2).Once()
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}
