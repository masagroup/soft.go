package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBasicEStoreList_Constructors(t *testing.T) {
	{
		mockOwner := &MockEObject{}
		mockFeature := &MockEStructuralFeature{}
		mockStore := &MockEStore{}
		list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)
		assert.NotNil(t, list)
	}
	{
		mockOwner := &MockEObject{}
		mockReference := &MockEReference{}
		mockStore := &MockEStore{}
		mockReference.On("IsContainment").Return(true).Once()
		mockReference.On("IsResolveProxies").Return(false).Once()
		mockReference.On("IsUnsettable").Return(false).Once()
		mockReference.On("GetEOpposite").Return(nil).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
	}
	{
		mockOwner := &MockEObject{}
		mockReference := &MockEReference{}
		mockOpposite := &MockEReference{}
		mockStore := &MockEStore{}
		mockReference.On("IsContainment").Return(true).Once()
		mockReference.On("IsResolveProxies").Return(false).Once()
		mockReference.On("IsUnsettable").Return(false).Once()
		mockReference.On("GetEOpposite").Return(mockOpposite).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
	}
	{
		mockOwner := &MockEObject{}
		mockReference := &MockEReference{}
		mockStore := &MockEStore{}
		mockReference.On("IsContainment").Return(false).Once()
		mockReference.On("IsResolveProxies").Return(false).Once()
		mockReference.On("IsUnsettable").Return(false).Once()
		mockReference.On("GetEOpposite").Return(nil).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
	}
	{
		mockOwner := &MockEObject{}
		mockReference := &MockEReference{}
		mockOpposite := &MockEReference{}
		mockStore := &MockEStore{}
		mockReference.On("IsContainment").Return(false).Once()
		mockReference.On("IsResolveProxies").Return(false).Once()
		mockReference.On("IsUnsettable").Return(false).Once()
		mockReference.On("GetEOpposite").Return(mockOpposite).Once()
		list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
		assert.NotNil(t, list)
	}
}

func TestBasicEStoreList_Accessors(t *testing.T) {
	mockOwner := &MockEObject{}
	mockFeature := &MockEStructuralFeature{}
	mockStore := &MockEStore{}
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	assert.Equal(t, mockOwner, list.GetOwner())
	assert.Equal(t, mockFeature, list.GetFeature())
	assert.Equal(t, mockStore, list.GetStore())

	mockClass := &MockEClass{}
	mockClass.On("GetFeatureID", mockFeature).Return(0).Once()
	mockOwner.On("EClass").Return(mockClass).Once()
	assert.Equal(t, 0, list.GetFeatureID())
	mock.AssertExpectationsForObjects(t, mockOwner, mockClass, mockFeature, mockStore)
}

func TestBasicEStoreList_Add(t *testing.T) {
	mockOwner := &MockEObject{}
	mockFeature := &MockEStructuralFeature{}
	mockStore := &MockEStore{}
	mockAdapter := new(MockEAdapter)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	// already present
	mockStore.On("Size", mockOwner, mockFeature).Return(0).Twice()
	mockStore.On("Contains", mockOwner, mockFeature, 1).Return(true).Once()
	assert.False(t, list.Add(1))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// add 1 to the list
	mockStore.On("Size", mockOwner, mockFeature).Return(0).Twice()
	mockStore.On("Contains", mockOwner, mockFeature, 1).Return(false).Once()
	mockStore.On("Add", mockOwner, mockFeature, 0, 1).Once()
	mockOwner.On("EDeliver").Return(false).Once()
	assert.True(t, list.Add(1))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// add 2 to the list
	mockStore.On("Size", mockOwner, mockFeature).Return(1).Twice()
	mockStore.On("Contains", mockOwner, mockFeature, 2).Return(false).Once()
	mockStore.On("Add", mockOwner, mockFeature, 1, 2).Once()
	mockOwner.On("EDeliver").Return(true).Once()
	mockOwner.On("EAdapters").Return(NewImmutableEList([]interface{}{mockAdapter})).Once()
	mockOwner.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == ADD && n.GetNewValue() == 2
	})).Once()
	assert.True(t, list.Add(2))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

func TestBasicEStoreList_AddReferenceContainmentNoOpposite(t *testing.T) {
	mockOwner := &MockEObject{}
	mockReference := &MockEReference{}
	mockStore := &MockEStore{}
	mockReference.On("IsContainment").Return(true).Once()
	mockReference.On("IsResolveProxies").Return(false).Once()
	mockReference.On("IsUnsettable").Return(false).Once()
	mockReference.On("GetEOpposite").Return(nil).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)

	mockObject := &MockEObjectInternal{}
	mockStore.On("Size", mockOwner, mockReference).Return(0).Twice()
	mockStore.On("Contains", mockOwner, mockReference, mockObject).Return(false).Once()
	mockStore.On("Add", mockOwner, mockReference, 0, mockObject).Once()
	mockReference.On("GetFeatureID").Return(0).Once()
	mockObject.On("EInverseAdd", mockOwner, EOPPOSITE_FEATURE_BASE-0, nil).Return(nil).Once()
	mockOwner.On("EDeliver").Return(false).Once()
	assert.True(t, list.Add(mockObject))
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject)
}

func TestBasicEStoreList_AddReferenceContainmentOpposite(t *testing.T) {
	mockOwner := &MockEObject{}
	mockReference := &MockEReference{}
	mockOpposite := &MockEReference{}
	mockStore := &MockEStore{}
	mockReference.On("IsContainment").Return(true).Once()
	mockReference.On("IsResolveProxies").Return(false).Once()
	mockReference.On("IsUnsettable").Return(false).Once()
	mockReference.On("GetEOpposite").Return(mockOpposite).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)

	mockObject := &MockEObjectInternal{}
	mockClass := &MockEClass{}
	mockStore.On("Size", mockOwner, mockReference).Return(0).Twice()
	mockStore.On("Contains", mockOwner, mockReference, mockObject).Return(false).Once()
	mockStore.On("Add", mockOwner, mockReference, 0, mockObject).Once()
	mockReference.On("GetEOpposite").Return(mockOpposite).Once()
	mockObject.On("EClass").Return(mockClass).Once()
	mockClass.On("GetFeatureID", mockOpposite).Return(1).Once()
	mockObject.On("EInverseAdd", mockOwner, 1, nil).Return(nil).Once()
	mockOwner.On("EDeliver").Return(false).Once()
	assert.True(t, list.Add(mockObject))
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject, mockClass, mockOpposite)
}

func TestBasicEStoreList_AddWithNotification(t *testing.T) {
	mockOwner := &MockEObject{}
	mockFeature := &MockEStructuralFeature{}
	mockStore := &MockEStore{}
	mockNotifications := &MockENotificationChain{}
	mockAdapter := new(MockEAdapter)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	// add 1
	mockStore.On("Size", mockOwner, mockFeature).Return(1).Once()
	mockStore.On("Add", mockOwner, mockFeature, 1, 2).Once()
	mockOwner.On("EDeliver").Return(true).Once()
	mockOwner.On("EAdapters").Return(NewImmutableEList([]interface{}{mockAdapter})).Once()
	mockNotifications.On("Add", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == ADD && n.GetNewValue() == 2
	})).Return(true).Once()
	assert.Equal(t, mockNotifications, list.AddWithNotification(2, mockNotifications))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockAdapter, mockNotifications)
}

func TestBasicEStoreList_Insert(t *testing.T) {
	mockOwner := &MockEObject{}
	mockFeature := &MockEStructuralFeature{}
	mockStore := &MockEStore{}
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	assert.Panics(t, func() {
		list.Insert(-1, 0)
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockStore.On("Size", mockOwner, mockFeature).Return(1).Once()
	assert.Panics(t, func() {
		list.Insert(2, 0)
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

func TestBasicEStoreList_AddAll(t *testing.T) {
	mockOwner := &MockEObject{}
	mockFeature := &MockEStructuralFeature{}
	mockStore := &MockEStore{}
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	mockStore.On("Size", mockOwner, mockFeature).Return(0).Twice()
	mockStore.On("Contains", mockOwner, mockFeature, 1).Return(false).Once()
	mockStore.On("Add", mockOwner, mockFeature, 0, 1).Once()
	mockOwner.On("EDeliver").Return(false).Once()
	assert.True(t, list.AddAll(NewImmutableEList([]interface{}{1})))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)
}

func TestBasicEStoreList_InsertAll(t *testing.T) {
	mockOwner := &MockEObject{}
	mockFeature := &MockEStructuralFeature{}
	mockStore := &MockEStore{}
	mockAdapter := new(MockEAdapter)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	// invalid index
	assert.Panics(t, func() {
		list.InsertAll(-1, NewImmutableEList([]interface{}{}))
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// already present element
	mockStore.On("Size", mockOwner, mockFeature).Return(0).Once()
	mockStore.On("Contains", mockOwner, mockFeature, 1).Return(true).Once()
	assert.False(t, list.InsertAll(0, NewImmutableEList([]interface{}{1})))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	// single element inserted
	mockStore.On("Size", mockOwner, mockFeature).Return(0).Once()
	mockStore.On("Contains", mockOwner, mockFeature, 1).Return(false).Once()
	mockStore.On("Add", mockOwner, mockFeature, 0, 1).Once()
	mockOwner.On("EDeliver").Return(true).Once()
	mockOwner.On("EAdapters").Return(NewImmutableEList([]interface{}{mockAdapter})).Once()
	mockOwner.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == ADD && n.GetNewValue() == 1
	})).Once()
	assert.True(t, list.InsertAll(0, NewImmutableEList([]interface{}{1})))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

	mockStore.On("Size", mockOwner, mockFeature).Return(0).Once()
	mockStore.On("Contains", mockOwner, mockFeature, 1).Return(false).Once()
	mockStore.On("Contains", mockOwner, mockFeature, 2).Return(false).Once()
	mockStore.On("Add", mockOwner, mockFeature, 0, 1).Once()
	mockStore.On("Add", mockOwner, mockFeature, 1, 2).Once()
	mockOwner.On("EDeliver").Return(true).Once()
	mockOwner.On("EAdapters").Return(NewImmutableEList([]interface{}{mockAdapter})).Once()
	mockOwner.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == ADD_MANY
	})).Once()
	assert.True(t, list.InsertAll(0, NewImmutableEList([]interface{}{1, 2})))
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore)

}

func TestBasicEStoreList_MoveObject(t *testing.T) {
	mockOwner := &MockEObject{}
	mockFeature := &MockEStructuralFeature{}
	mockStore := &MockEStore{}
	mockObject := &MockEObjectInternal{}
	mockAdapter := new(MockEAdapter)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	mockStore.On("IndexOf", mockOwner, mockFeature, 1).Return(-1).Once()
	assert.Panics(t, func() {
		list.MoveObject(1, 1)
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)

	mockStore.On("IndexOf", mockOwner, mockFeature, 1).Return(0).Once()
	mockStore.On("Size", mockOwner, mockFeature).Return(1).Twice()
	mockStore.On("Move", mockOwner, mockFeature, 1, 0).Return(mockObject).Once()
	mockOwner.On("EDeliver").Return(true).Once()
	mockOwner.On("EAdapters").Return(NewImmutableEList([]interface{}{mockAdapter})).Once()
	mockOwner.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == MOVE
	})).Once()
	list.MoveObject(1, 1)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)
}

func TestBasicEStoreList_Move(t *testing.T) {
	mockOwner := &MockEObject{}
	mockFeature := &MockEStructuralFeature{}
	mockStore := &MockEStore{}
	mockObject := &MockEObjectInternal{}
	mockAdapter := new(MockEAdapter)
	list := NewBasicEStoreList(mockOwner, mockFeature, mockStore)

	assert.Panics(t, func() {
		list.MoveObject(-1, 1)
	})
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)

	mockStore.On("Size", mockOwner, mockFeature).Return(1).Twice()
	mockStore.On("Move", mockOwner, mockFeature, 1, 0).Return(mockObject).Once()
	mockOwner.On("EDeliver").Return(true).Once()
	mockOwner.On("EAdapters").Return(NewImmutableEList([]interface{}{mockAdapter})).Once()
	mockOwner.On("ENotify", mock.MatchedBy(func(n ENotification) bool {
		return n.GetNotifier() == mockOwner && n.GetFeature() == mockFeature && n.GetEventType() == MOVE
	})).Once()
	list.Move(0, 1)
	mock.AssertExpectationsForObjects(t, mockOwner, mockFeature, mockStore, mockObject)
}

func TestBasicEStoreList_Get_NoProxy(t *testing.T) {
	mockOwner := &MockEObject{}
	mockReference := &MockEReference{}
	mockOpposite := &MockEReference{}
	mockStore := &MockEStore{}
	mockObject := &MockEObject{}
	mockReference.On("IsContainment").Return(false).Once()
	mockReference.On("IsResolveProxies").Return(false).Once()
	mockReference.On("IsUnsettable").Return(false).Once()
	mockReference.On("GetEOpposite").Return(mockOpposite).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)

	mockStore.On("Get", list.owner, list.feature, 0).Return(mockObject).Once()
	assert.Equal(t, mockObject, list.Get(0))
}

func TestBasicEStoreList_Get_Proxy(t *testing.T) {
	mockOwner := &MockEObjectInternal{}
	mockReference := &MockEReference{}
	mockOpposite := &MockEReference{}
	mockStore := &MockEStore{}
	mockObject := &MockEObjectInternal{}
	mockResolved := &MockEObjectInternal{}
	mockReference.On("IsContainment").Return(true).Once()
	mockReference.On("IsResolveProxies").Return(true).Once()
	mockReference.On("IsUnsettable").Return(false).Once()
	mockReference.On("GetEOpposite").Return(mockOpposite).Once()
	list := NewBasicEStoreList(mockOwner, mockReference, mockStore)
	assert.NotNil(t, list)
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject, mockOpposite)

	mockClass := &MockEClass{}
	mockStore.On("Get", list.owner, list.feature, 0).Return(mockObject).Once()
	mockObject.On("EIsProxy").Return(true).Once()
	mockOwner.On("EResolveProxy", mockObject).Return(mockResolved).Once()
	mockStore.On("Set", list.owner, list.feature, 0, mockResolved).Return(mockResolved).Return(mockObject).Once()
	mockReference.On("GetEOpposite").Return(mockOpposite).Once()
	mockClass.On("GetFeatureID", mockOpposite).Return(0).Once()
	mockObject.On("EClass").Return(mockClass).Once()
	mockObject.On("EInverseRemove", mockOwner, 0, nil).Return(nil).Once()
	mockResolved.On("EInternalContainer").Return(nil).Once()
	mockReference.On("GetEOpposite").Return(mockOpposite).Once()
	mockClass.On("GetFeatureID", mockOpposite).Return(0).Once()
	mockResolved.On("EClass").Return(mockClass).Once()
	mockResolved.On("EInverseAdd", mockOwner, 0, nil).Return(nil).Once()
	mockOwner.On("EDeliver").Return(false).Once()
	assert.Equal(t, mockResolved, list.Get(0))
	mock.AssertExpectationsForObjects(t, mockOwner, mockReference, mockStore, mockObject, mockOpposite, mockResolved, mockClass)
}
