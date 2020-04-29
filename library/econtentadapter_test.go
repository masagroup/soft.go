package library

import (
	"net/url"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type contentAdapter struct {
	mock.Mock
	*ecore.EContentAdapter
}

func newContentAdapter() *contentAdapter {
	c := new(contentAdapter)
	c.EContentAdapter = ecore.NewEContentAdapter()
	c.SetInterfaces(c)
	return c
}

func (adapter *contentAdapter) NotifyChanged(notification ecore.ENotification) {
	adapter.Called(notification)
	adapter.EContentAdapter.NotifyChanged(notification)
}

func TestEContentAdapterEObjectTopChanged(t *testing.T) {

	library := GetFactory().CreateLibrary()
	employee := GetFactory().CreateEmployee()

	// create and add a new content adapter
	adapter := newContentAdapter()
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
	adapter := newContentAdapter()
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

	fileURI := &url.URL{Path: "testdata/library.xml"}
	resourceFactory := ecore.GetResourceFactoryRegistry().GetFactory(fileURI)
	resource := resourceFactory.CreateResource(fileURI)
	resource.Load()
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty())
	assert.True(t, resource.GetWarnings().Empty())

	adapter := newContentAdapter()
	resource.EAdapters().Add(adapter)
}
