package integration

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/masagroup/soft.go/test/library"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEContentAdapterEObjectTopChanged(t *testing.T) {

	lib := library.GetFactory().CreateLibrary()
	employee := library.GetFactory().CreateEmployee()

	// create and add a new content adapter
	adapter := ecore.NewMockEContentAdapter(t)
	lib.EAdapters().Add(adapter)

	// check that a notification is raised when employee is added
	adapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ecore.ENotification) bool {
		return n.GetNotifier() == lib &&
			n.GetFeatureID() == library.LIBRARY__EMPLOYEES &&
			n.GetNewValue() == employee &&
			n.GetOldValue() == nil &&
			n.GetEventType() == ecore.ADD &&
			n.GetPosition() == 0
	})).Once()

	// add employee
	lib.GetEmployees().Add(employee)

	// check expectations
	adapter.AssertExpectations(t)
}

func TestEContentAdapterEObjectChildChanged(t *testing.T) {

	lib := library.GetFactory().CreateLibrary()
	employee := library.GetFactory().CreateEmployee()
	employee.SetFirstName("oldName")
	lib.GetEmployees().Add(employee)

	// create and add a new content adapter
	adapter := ecore.NewMockEContentAdapter(t)
	lib.EAdapters().Add(adapter)

	// check that a notification is raised when employee is added
	adapter.EXPECT().NotifyChanged(mock.MatchedBy(func(n ecore.ENotification) bool {
		return n.GetNotifier() == employee &&
			n.GetFeatureID() == library.EMPLOYEE__FIRST_NAME &&
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
	ecore.GetPackageRegistry().RegisterPackage(library.GetPackage())

	fileURI := ecore.CreateFileURI("testdata/library/library.complex.xml")
	resource := ecore.NewEResourceImpl()
	resource.SetURI(fileURI)
	resource.Load()
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty())
	assert.True(t, resource.GetWarnings().Empty())

	adapter := ecore.NewMockEContentAdapter(t)
	resource.EAdapters().Add(adapter)
}
