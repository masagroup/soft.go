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

import "reflect"

type MockEClassifier struct {
	MockENamedElement
}

// GetClassifierID get the value of classifierID
func (eClassifier *MockEClassifier) GetClassifierID() int {
	ret := eClassifier.Called()

	var r int
	if rf, ok := ret.Get(0).(func() int); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(int)
		}
	}

	return r
}

// SetClassifierID provides mock implementation for setting the value of classifierID
func (eClassifier *MockEClassifier) SetClassifierID(newClassifierID int) {
	eClassifier.Called(newClassifierID)
}

// GetDefaultValue get the value of defaultValue
func (eClassifier *MockEClassifier) GetDefaultValue() interface{} {
	ret := eClassifier.Called()

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

// GetEPackage get the value of ePackage
func (eClassifier *MockEClassifier) GetEPackage() EPackage {
	ret := eClassifier.Called()

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

// GetInstanceClass get the value of instanceClass
func (eClassifier *MockEClassifier) GetInstanceClass() reflect.Type {
	ret := eClassifier.Called()

	var r reflect.Type
	if rf, ok := ret.Get(0).(func() reflect.Type); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(reflect.Type)
		}
	}

	return r
}

// SetInstanceClass provides mock implementation for setting the value of instanceClass
func (eClassifier *MockEClassifier) SetInstanceClass(newInstanceClass reflect.Type) {
	eClassifier.Called(newInstanceClass)
}

// GetInstanceTypeName get the value of instanceTypeName
func (eClassifier *MockEClassifier) GetInstanceTypeName() string {
	ret := eClassifier.Called()

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

// SetInstanceTypeName provides mock implementation for setting the value of instanceTypeName
func (eClassifier *MockEClassifier) SetInstanceTypeName(newInstanceTypeName string) {
	eClassifier.Called(newInstanceTypeName)
}

// IsInstance provides mock implementation
func (eClassifier *MockEClassifier) IsInstance(object interface{}) bool {
	ret := eClassifier.Called(object)

	var r bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(bool)
		}
	}

	return r
}
