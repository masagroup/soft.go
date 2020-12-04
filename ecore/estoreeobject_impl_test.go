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

func TestEStoreEObjectImpl_GetAttribute_Transient(t *testing.T) {
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
	mockAttribute.On("IsTransient").Return(true).Once()
	mockClass.On("GetFeatureCount").Return(1).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Twice()
	assert.Nil(t, o.EGetFromID(0, false))
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
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

func TestEStoreEObjectImpl_SetAttribute_Transient(t *testing.T) {
	// create object
	o := newEStoreEObjectTest(true)

	// create mocks
	mockClass := new(MockEClass)
	mockAttribute := new(MockEAttribute)
	mockStore := o.mockStore

	// initialise object with mock class
	o.setEClass(mockClass)

	// set
	mockAttribute.On("IsTransient").Return(true).Twice()
	mockClass.On("GetFeatureCount").Return(1).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Times(3)
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	// test get
	mockClass.On("GetFeatureCount").Return(1).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Twice()
	assert.Equal(t, 2, o.EGetFromID(0, false))
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
	mockStore.On("Set", o, mockAttribute, NO_INDEX, 2).Return(nil).Once()
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_GetAttribute_Many(t *testing.T) {

	// create object
	o := newEStoreEObjectTest(false)

	// create mocks
	mockClass := new(MockEClass)
	mockAttribute := new(MockEAttribute)
	mockStore := o.mockStore

	// initialise object with mock class
	o.setEClass(mockClass)

	mockAttribute.On("IsMany").Return(true).Once()
	mockAttribute.On("IsTransient").Return(false).Once()
	mockClass.On("GetFeatureCount").Return(1).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Twice()
	list := o.EGetFromID(0, true)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)

	eobjectlist, _ := list.(EObjectList)
	assert.NotNil(t, eobjectlist)
	enotifyinglist, _ := list.(ENotifyingList)
	assert.NotNil(t, enotifyinglist)
}

func TestEStoreEObjectImpl_SetAttribute_Caching(t *testing.T) {

	// create object
	o := newEStoreEObjectTest(true)

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
	mockStore.On("Set", o, mockAttribute, NO_INDEX, 2).Return(nil).Once()
	o.ESetFromID(0, 2)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}

func TestEStoreEObjectImpl_UnSetAttribute_Transient(t *testing.T) {
	// create object
	o := newEStoreEObjectTest(false)

	// create mocks
	mockClass := new(MockEClass)
	mockAttribute := new(MockEAttribute)

	// initialise object with mock class
	o.setEClass(mockClass)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("IsTransient").Return(true).Twice()
	mockClass.On("GetFeatureCount").Return(1).Once()
	mockClass.On("GetEStructuralFeature", 0).Return(mockAttribute).Times(3)
	o.EUnsetFromID(0)

}

func TestEStoreEObjectImpl_UnSetAttribute(t *testing.T) {

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
	mockStore.On("Get", o, mockAttribute, NO_INDEX).Return(2).Once()
	mockStore.On("UnSet", o, mockAttribute).Once()
	o.EUnsetFromID(0)
	mock.AssertExpectationsForObjects(t, mockClass, mockAttribute, mockStore)
}
