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

type MockEModelElement struct {
	MockEObjectInternal
}

// GetEAnnotations get the value of eAnnotations
func (eModelElement *MockEModelElement) GetEAnnotations() EList {
	ret := eModelElement.Called()

	var r EList
	if rf, ok := ret.Get(0).(func() EList); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(EList)
		}
	}

	return r
}

// GetEAnnotation provides mock implementation
func (eModelElement *MockEModelElement) GetEAnnotation(source string) EAnnotation {
	ret := eModelElement.Called(source)

	var r EAnnotation
	if rf, ok := ret.Get(0).(func() EAnnotation); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(EAnnotation)
		}
	}

	return r
}
