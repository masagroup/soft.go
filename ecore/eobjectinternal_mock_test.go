package ecore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEObjectInternal_EProperties(t *testing.T) {
	o := &MockEObjectInternal{}
	p := &MockEObjectProperties{}
	// return a value
	o.On("EProperties").Once().Return(p)
	o.On("EProperties").Once().Return(func() EDynamicProperties {
		return p
	})
	assert.Equal(t, p, o.EDynamicProperties())
	assert.Equal(t, p, o.EDynamicProperties())
	mock.AssertExpectationsForObjects(t, o, p)
}

func TestMockEObjectInternal_EStaticClass(t *testing.T) {
	o := &MockEObjectInternal{}
	c := &MockEClass{}
	// return a value
	o.On("EStaticClass").Once().Return(c)
	o.On("EStaticClass").Once().Return(func() EClass {
		return c
	})
	assert.Equal(t, c, o.EStaticClass())
	assert.Equal(t, c, o.EStaticClass())
	mock.AssertExpectationsForObjects(t, o, c)
}

func TestMockEObjectInternal_EStaticFeatureCount(t *testing.T) {
	o := &MockEObjectInternal{}
	// return a value
	o.On("EStaticFeatureCount").Once().Return(1)
	o.On("EStaticFeatureCount").Once().Return(func() int {
		return 2
	})
	assert.Equal(t, 1, o.EStaticFeatureCount())
	assert.Equal(t, 2, o.EStaticFeatureCount())
	o.AssertExpectations(t)
}

func TestMockEObjectInternal_EInternalResource(t *testing.T) {
	o := &MockEObjectInternal{}
	r := &MockEResource{}
	// return a value
	o.On("EInternalResource").Once().Return(r)
	o.On("EInternalResource").Once().Return(func() EResource {
		return r
	})
	assert.Equal(t, r, o.EInternalResource())
	assert.Equal(t, r, o.EInternalResource())
	mock.AssertExpectationsForObjects(t, o, r)
}

func TestMockEObjectInternal_EInternalContainer(t *testing.T) {
	o := &MockEObjectInternal{}
	c := &MockEObject{}
	// return a value
	o.On("EInternalContainer").Once().Return(c)
	o.On("EInternalContainer").Once().Return(func() EObject {
		return c
	})
	assert.Equal(t, c, o.EInternalContainer())
	assert.Equal(t, c, o.EInternalContainer())
	mock.AssertExpectationsForObjects(t, o, c)
}

func TestMockEObjectInternal_ESetResource(t *testing.T) {
	o := &MockEObjectInternal{}
	r := &MockEResource{}
	n := &MockENotificationChain{}

	// return a value
	o.On("ESetResource", r, n).Once().Return(n)
	o.On("ESetResource", r, n).Once().Return(func(resource EResource, notifications ENotificationChain) ENotificationChain {
		return n
	})
	assert.Equal(t, n, o.ESetResource(r, n))
	assert.Equal(t, n, o.ESetResource(r, n))
	mock.AssertExpectationsForObjects(t, r, n)
}

func TestMockEObjectInternal_EInverseAdd(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}
	n := &MockENotificationChain{}

	// return a value
	o.On("EInverseAdd", obj, 1, n).Once().Return(n)
	o.On("EInverseAdd", obj, 1, n).Once().Return(func(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
		return n
	})
	assert.Equal(t, n, o.EInverseAdd(obj, 1, n))
	assert.Equal(t, n, o.EInverseAdd(obj, 1, n))
	mock.AssertExpectationsForObjects(t, o, obj, n)
}

func TestMockEObjectInternal_EInverseRemove(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}
	n := &MockENotificationChain{}

	// return a value
	o.On("EInverseRemove", obj, 1, n).Once().Return(n)
	o.On("EInverseRemove", obj, 1, n).Once().Return(func(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
		return n
	})
	assert.Equal(t, n, o.EInverseRemove(obj, 1, n))
	assert.Equal(t, n, o.EInverseRemove(obj, 1, n))
	mock.AssertExpectationsForObjects(t, o, obj, n)
}

func TestMockEObjectInternal_EBasicInverseAdd(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}
	n := &MockENotificationChain{}

	// return a value
	o.On("EBasicInverseAdd", obj, 1, n).Once().Return(n)
	o.On("EBasicInverseAdd", obj, 1, n).Once().Return(func(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
		return n
	})
	assert.Equal(t, n, o.EBasicInverseAdd(obj, 1, n))
	assert.Equal(t, n, o.EBasicInverseAdd(obj, 1, n))
	mock.AssertExpectationsForObjects(t, o, obj, n)
}

func TestMockEObjectInternal_EBasicInverseRemove(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}
	n := &MockENotificationChain{}

	// return a value
	o.On("EBasicInverseRemove", obj, 1, n).Once().Return(n)
	o.On("EBasicInverseRemove", obj, 1, n).Once().Return(func(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
		return n
	})
	assert.Equal(t, n, o.EBasicInverseRemove(obj, 1, n))
	assert.Equal(t, n, o.EBasicInverseRemove(obj, 1, n))
	mock.AssertExpectationsForObjects(t, o, obj, n)
}

func TestMockEObjectInternal_EDerivedFeatureID(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}

	// return a value
	o.On("EDerivedFeatureID", obj, 1).Once().Return(2)
	o.On("EDerivedFeatureID", obj, 1).Once().Return(func(container EObject, featureID int) int {
		return 2
	})
	assert.Equal(t, 2, o.EDerivedFeatureID(obj, 1))
	assert.Equal(t, 2, o.EDerivedFeatureID(obj, 1))
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_EDerivedOperationID(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}

	// return a value
	o.On("EDerivedOperationID", obj, 1).Once().Return(2)
	o.On("EDerivedOperationID", obj, 1).Once().Return(func(container EObject, featureID int) int {
		return 2
	})
	assert.Equal(t, 2, o.EDerivedOperationID(obj, 1))
	assert.Equal(t, 2, o.EDerivedOperationID(obj, 1))
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_EGetFromID(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}

	// return a value
	o.On("EGetFromID", 1, false).Once().Return(obj)
	o.On("EGetFromID", 1, true).Once().Return(func(featureID int, resolve bool) interface{} {
		return obj
	})
	assert.Equal(t, obj, o.EGetFromID(1, false))
	assert.Equal(t, obj, o.EGetFromID(1, true))
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_EIsSetFromID(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}

	// return a value
	o.On("EIsSetFromID", 1).Once().Return(false)
	o.On("EIsSetFromID", 1).Once().Return(func(featureID int) bool {
		return true
	})
	assert.False(t, o.EIsSetFromID(1))
	assert.True(t, o.EIsSetFromID(1))
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_ESetFromID(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}

	// return a value
	o.On("ESetFromID", 1, obj).Once()
	o.ESetFromID(1, obj)
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_EUnsetFromID(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}

	// return a value
	o.On("EUnsetFromID", 1).Once()
	o.EUnsetFromID(1)
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_EInvokeFromID(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}
	args := &MockEList{}
	// return a value
	o.On("EInvokeFromID", 1, args).Once().Return(obj)
	o.On("EInvokeFromID", 1, args).Once().Return(func(operationID int, arguments EList) interface{} {
		return obj
	})
	assert.Equal(t, obj, o.EInvokeFromID(1, args))
	assert.Equal(t, obj, o.EInvokeFromID(1, args))
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_EObjectForFragmentSegment(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}
	// return a value
	o.On("EObjectForFragmentSegment", "uri").Once().Return(obj)
	o.On("EObjectForFragmentSegment", "uri").Once().Return(func(string) EObject {
		return obj
	})
	assert.Equal(t, obj, o.EObjectForFragmentSegment("uri"))
	assert.Equal(t, obj, o.EObjectForFragmentSegment("uri"))
	mock.AssertExpectationsForObjects(t, o, obj)
}

func TestMockEObjectInternal_EURIFragmentSegment(t *testing.T) {
	o := &MockEObjectInternal{}
	f := &MockEStructuralFeature{}
	obj := &MockEObject{}
	// return a value
	o.On("EURIFragmentSegment", f, obj).Once().Return("uri")
	o.On("EURIFragmentSegment", f, obj).Once().Return(func(EStructuralFeature, EObject) string {
		return "uri"
	})
	assert.Equal(t, "uri", o.EURIFragmentSegment(f, obj))
	assert.Equal(t, "uri", o.EURIFragmentSegment(f, obj))
	mock.AssertExpectationsForObjects(t, o, f, obj)
}

func TestMockEObjectInternal_EProxyURI(t *testing.T) {
	o := &MockEObjectInternal{}
	uri, _ := url.Parse("test://file.t")

	// return a value
	o.On("EProxyURI").Once().Return(uri)
	o.On("EProxyURI").Once().Return(func() *url.URL {
		return uri
	})
	assert.Equal(t, uri, o.EProxyURI())
	assert.Equal(t, uri, o.EProxyURI())
	mock.AssertExpectationsForObjects(t, o)
}

func TestMockEObjectInternal_ESetProxyURI(t *testing.T) {
	o := &MockEObjectInternal{}
	uri, _ := url.Parse("test://file.t")

	// return a value
	o.On("ESetProxyURI", uri).Once()
	o.ESetProxyURI(uri)
	mock.AssertExpectationsForObjects(t, o)
}

func TestMockEObjectInternal_EResolveProxy(t *testing.T) {
	o := &MockEObjectInternal{}
	obj := &MockEObject{}
	result := &MockEObject{}

	// return a value
	o.On("EResolveProxy", obj).Once().Return(result)
	o.On("EResolveProxy", obj).Once().Return(func(proxy EObject) EObject {
		return result
	})
	assert.Equal(t, result, o.EResolveProxy(obj))
	assert.Equal(t, result, o.EResolveProxy(obj))
	mock.AssertExpectationsForObjects(t, o)
}
