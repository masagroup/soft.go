package ecore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEResourceSetGetResources(t *testing.T) {
	rs := &MockEResourceSet{}
	l := &MockEList{}
	rs.On("GetResources").Return(l).Once()
	rs.On("GetResources").Return(func() EList {
		return l
	}).Once()
	assert.Equal(t, l, rs.GetResources())
	assert.Equal(t, l, rs.GetResources())
	mock.AssertExpectationsForObjects(t, rs, l)
}

func TestMockEResourceSetGetResource(t *testing.T) {
	rs := &MockEResourceSet{}
	r := &MockEResource{}
	uri, _ := url.Parse("test://file.t")
	rs.On("GetResource", uri, false).Return(r).Once()
	rs.On("GetResource", uri, true).Return(func(uri *url.URL, loadOnDemand bool) EResource {
		return r
	}).Once()
	assert.Equal(t, r, rs.GetResource(uri, false))
	assert.Equal(t, r, rs.GetResource(uri, true))
	mock.AssertExpectationsForObjects(t, r, rs)
}

func TestMockEResourceSetCreateResource(t *testing.T) {
	rs := &MockEResourceSet{}
	r := &MockEResource{}
	uri, _ := url.Parse("test://file.t")
	rs.On("CreateResource", uri).Return(r).Once()
	rs.On("CreateResource", uri).Return(func(uri *url.URL) EResource {
		return r
	}).Once()
	assert.Equal(t, r, rs.CreateResource(uri))
	assert.Equal(t, r, rs.CreateResource(uri))
	mock.AssertExpectationsForObjects(t, r, rs)
}

func TestMockEResourceSetGetEObject(t *testing.T) {
	rs := &MockEResourceSet{}
	o := &MockEObject{}
	uri, _ := url.Parse("test://file.t")
	rs.On("GetEObject", uri, false).Return(o).Once()
	rs.On("GetEObject", uri, true).Return(func(uri *url.URL, loadOnDemand bool) EObject {
		return o
	}).Once()
	assert.Equal(t, o, rs.GetEObject(uri, false))
	assert.Equal(t, o, rs.GetEObject(uri, true))
	mock.AssertExpectationsForObjects(t, o, rs)
}

func TestMockEResourceSetGetURIConverter(t *testing.T) {
	rs := &MockEResourceSet{}
	c := &MockEURIConverter{}
	rs.On("GetURIConverter").Return(c).Once()
	rs.On("GetURIConverter").Return(func() EURIConverter {
		return c
	}).Once()
	assert.Equal(t, c, rs.GetURIConverter())
	assert.Equal(t, c, rs.GetURIConverter())
	mock.AssertExpectationsForObjects(t, rs, c)
}

func TestMockEResourceSetSetURIConverter(t *testing.T) {
	rs := &MockEResourceSet{}
	c := &MockEURIConverter{}
	rs.On("SetURIConverter", c).Once()
	rs.SetURIConverter(c)
	mock.AssertExpectationsForObjects(t, rs, c)
}

func TestMockEResourceSetGetPackageRegistry(t *testing.T) {
	rs := &MockEResourceSet{}
	pr := &MockEPackageRegistry{}
	rs.On("GetPackageRegistry").Return(pr).Once()
	rs.On("GetPackageRegistry").Return(func() EPackageRegistry {
		return pr
	}).Once()
	assert.Equal(t, pr, rs.GetPackageRegistry())
	assert.Equal(t, pr, rs.GetPackageRegistry())
	mock.AssertExpectationsForObjects(t, rs, pr)
}

func TestMockEResourceSetSetPackageRegistry(t *testing.T) {
	rs := &MockEResourceSet{}
	pr := &MockEPackageRegistry{}
	rs.On("SetPackageRegistry", pr).Once()
	rs.SetPackageRegistry(pr)
	mock.AssertExpectationsForObjects(t, rs, pr)
}

func TestMockEResourceSetGetResourceFactoryRegistry(t *testing.T) {
	rs := &MockEResourceSet{}
	pr := &MockEResourceFactoryRegistry{}
	rs.On("GetResourceFactoryRegistry").Return(pr).Once()
	rs.On("GetResourceFactoryRegistry").Return(func() EResourceFactoryRegistry {
		return pr
	}).Once()
	assert.Equal(t, pr, rs.GetResourceFactoryRegistry())
	assert.Equal(t, pr, rs.GetResourceFactoryRegistry())
	mock.AssertExpectationsForObjects(t, rs, pr)
}

func TestMockEResourceSetSetResourceFactoryRegistry(t *testing.T) {
	rs := &MockEResourceSet{}
	pr := &MockEResourceFactoryRegistry{}
	rs.On("SetResourceFactoryRegistry", pr).Once()
	rs.SetResourceFactoryRegistry(pr)
	mock.AssertExpectationsForObjects(t, rs, pr)
}

func TestMockEResourceSetGetURIResourceMap(t *testing.T) {
	rs := &MockEResourceSet{}
	pr := make(map[*url.URL]EResource)
	rs.On("GetURIResourceMap").Return(pr).Once()
	rs.On("GetURIResourceMap").Return(func() map[*url.URL]EResource {
		return pr
	}).Once()
	assert.Equal(t, pr, rs.GetURIResourceMap())
	assert.Equal(t, pr, rs.GetURIResourceMap())
	mock.AssertExpectationsForObjects(t, rs)
}

func TestMockEResourceSetSetURIResourceMap(t *testing.T) {
	rs := &MockEResourceSet{}
	pr := make(map[*url.URL]EResource)
	rs.On("SetURIResourceMap", pr).Once()
	rs.SetURIResourceMap(pr)
	mock.AssertExpectationsForObjects(t, rs)
}
