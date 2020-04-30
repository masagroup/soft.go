// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

// *****************************************************************************
//
// Warning: This file was generated by soft.generator.go Generator
//
// *****************************************************************************

package ecore

type MockEFactory struct {
	MockEModelElement
}

// GetEPackage get the value of ePackage
func (eFactory *MockEFactory) GetEPackage() EPackage {
	ret := eFactory.Called()

	var r EPackage
	if rf, ok := ret.Get(0).(func() EPackage); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(EPackage)
		}
	}

	return r
}

// SetEPackage provides mock implementation for setting the value of ePackage
func (eFactory *MockEFactory) SetEPackage(newEPackage EPackage) {
	eFactory.Called(newEPackage)
}

// Create provides mock implementation
func (eFactory *MockEFactory) Create(eClass EClass) EObject {
	ret := eFactory.Called(eClass)

	var r EObject
	if rf, ok := ret.Get(0).(func() EObject); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(EObject)
		}
	}

	return r
}

// CreateFromString provides mock implementation
func (eFactory *MockEFactory) CreateFromString(eDataType EDataType, literalValue string) interface{} {
	ret := eFactory.Called(eDataType, literalValue)

	var r interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(interface{})
		}
	}

	return r
}

// ConvertToString provides mock implementation
func (eFactory *MockEFactory) ConvertToString(eDataType EDataType, instanceValue interface{}) string {
	ret := eFactory.Called(eDataType, instanceValue)

	var r string
	if rf, ok := ret.Get(0).(func() string); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(string)
		}
	}

	return r
}
