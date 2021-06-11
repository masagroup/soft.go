package library

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEContentAdapterEObjectTopChanged(t *testing.T) {

	library := GetFactory().CreateLibrary()
	employee := GetFactory().CreateEmployee()

	// create and add a new content adapter
	adapter := ecore.NewMockContentAdapter()
	library.EAdapters().Add(adapter)

	// check that a notification is raised when employee is added
	adapter.On("NotifyChanged", mock.MatchedBy(func(n ecore.ENotification) bool {
		return n.GetNotifier() == library &&
			n.GetFeatureID() == LIBRARY__EMPLOYEES &&
			n.GetNewValue() == employee &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ecore.ADD &&
			n.GetPosition() == 0
	})).Once()

	// add employee
	library.GetEmployees().Add(employee)

	// check expectations
	adapter.AssertExpectations(t)
}

func TestEContentAdapterEObjectChildChanged(t *testing.T) {

	library := GetFactory().CreateLibrary()
	employee := GetFactory().CreateEmployee()
	employee.SetFirstName("oldName")
	library.GetEmployees().Add(employee)

	// create and add a new content adapter
	adapter := ecore.NewMockContentAdapter()
	library.EAdapters().Add(adapter)

	// check that a notification is raised when employee is added
	adapter.On("NotifyChanged", mock.MatchedBy(func(n ecore.ENotification) bool {
		return n.GetNotifier() == employee &&
			n.GetFeatureID() == EMPLOYEE__FIRST_NAME &&
			n.GetNewValue() == "newName" &&
			n.GetOldValue() == "oldName" &&
			n.GetEventType() == ecore.SET &&
			n.GetPosition() == -1
	})).Once()

	// change employee name
	employee.SetFirstName("newName")

	// check expectations
	adapter.AssertExpectations(t)
}

func TestEContentAdapterResource(t *testing.T) {
	ecore.GetPackageRegistry().RegisterPackage(GetPackage())

	fileURI := ecore.CreateFileURI("testdata/library.complex.xml")
	resource := ecore.NewEResourceImpl()
	resource.SetURI(fileURI)
	resource.Load()
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty())
	assert.True(t, resource.GetWarnings().Empty())

	adapter := ecore.NewMockContentAdapter()
	resource.EAdapters().Add(adapter)
}
