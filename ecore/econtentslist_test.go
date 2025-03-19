package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestEContentsListNoFeatures(t *testing.T) {
	mockObject := NewMockEObject(t)
	value := struct{}{}
	mockList := NewEmptyImmutableEList()
	l := newEContentsList(mockObject, mockList, false)
	assert.Equal(t, 0, l.Size())
	assert.Equal(t, -1, l.IndexOf(value))
	assert.False(t, l.Contains(value))

	it := l.Iterator()
	assert.False(t, it.HasNext())
	assert.Panics(t, func() {
		it.Next()
	})
}

type EContentsListTestSuite struct {
	suite.Suite
	mockObject  *MockEObject
	mockFeature *MockEStructuralFeature
	l           *eContentsList
}

func (suite *EContentsListTestSuite) SetupTest() {
	t := suite.T()
	suite.mockObject = NewMockEObject(t)
	suite.mockFeature = NewMockEStructuralFeature(t)
	suite.l = newEContentsList(suite.mockObject, NewImmutableEList([]any{suite.mockFeature}), false)
}

func (suite *EContentsListTestSuite) AfterTest() {
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)
}

func (suite *EContentsListTestSuite) TestEmpty() {
	// single not set
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Return(false).Once()
	assert.True(suite.T(), suite.l.Empty())
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)

	// single set
	value := struct{}{}
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Return(true).Once()
	suite.mockObject.EXPECT().EGetResolve(suite.mockFeature, false).Return(value).Once()
	suite.mockFeature.EXPECT().IsMany().Return(false).Once()
	assert.False(suite.T(), suite.l.Empty())

	// many
	mockList := NewMockEList(suite.T())
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Return(true).Once()
	suite.mockObject.EXPECT().EGetResolve(suite.mockFeature, false).Return(mockList).Once()
	suite.mockFeature.EXPECT().IsMany().Return(true).Once()
	mockList.EXPECT().Empty().Return(false).Once()
	assert.False(suite.T(), suite.l.Empty())
}

func (suite *EContentsListTestSuite) TestSize() {
	// single not set
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Return(false).Once()
	assert.Equal(suite.T(), 0, suite.l.Size())
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)

	// single set
	value := struct{}{}
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Return(true).Once()
	suite.mockObject.EXPECT().EGetResolve(suite.mockFeature, false).Return(value).Once()
	suite.mockFeature.EXPECT().IsMany().Return(false).Once()
	assert.Equal(suite.T(), 1, suite.l.Size())
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)

	// many
	l := NewImmutableEList([]any{struct{}{}})
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Return(true).Once()
	suite.mockObject.EXPECT().EGetResolve(suite.mockFeature, false).Return(l).Once()
	suite.mockFeature.EXPECT().IsMany().Return(true).Once()
	assert.Equal(suite.T(), 1, suite.l.Size())
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)
}

func (suite *EContentsListTestSuite) TestContains() {

}

func (suite *EContentsListTestSuite) TestIndexOf() {

}

func (suite *EContentsListTestSuite) TestGet() {

}

func (suite *EContentsListTestSuite) TestAll() {

}

func (suite *EContentsListTestSuite) TestGetUnResolved() {
	assert.Equal(suite.T(), suite.l, suite.l.GetUnResolvedList())
}

func TestEContentsList(t *testing.T) {
	suite.Run(t, &EContentsListTestSuite{})
}

type EContentsListIteratorTestSuite struct {
	EContentsListTestSuite
	it EIterator
}

func (suite *EContentsListIteratorTestSuite) SetupTest() {
	suite.EContentsListTestSuite.SetupTest()
	suite.it = suite.l.Iterator()
}

func (suite *EContentsListIteratorTestSuite) TestIteratorEmpty() {
	t := suite.T()
	it := suite.it
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Once().Return(false)
	assert.False(t, it.HasNext())
	assert.Panics(t, func() {
		it.Next()
	})
}

func (suite *EContentsListIteratorTestSuite) TestIteratorSingle() {
	t := suite.T()
	it := suite.it
	value := struct{}{}
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Once().Return(true)
	suite.mockObject.EXPECT().EGetResolve(suite.mockFeature, false).Once().Return(value)
	suite.mockFeature.EXPECT().IsMany().Once().Return(false)
	assert.True(t, it.HasNext())
	assert.Equal(t, value, it.Next())
}

func (suite *EContentsListIteratorTestSuite) TestIteratorManyEmpty() {
	t := suite.T()
	it := suite.it
	mockList := NewMockEList(t)
	mockIterator := &MockEIterator{}
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Once().Return(true)
	suite.mockObject.EXPECT().EGetResolve(suite.mockFeature, false).Once().Return(mockList)
	suite.mockFeature.EXPECT().IsMany().Once().Return(true)
	mockList.EXPECT().Iterator().Return(mockIterator).Once()
	mockIterator.EXPECT().HasNext().Return(false).Once()
	assert.False(t, it.HasNext())
	assert.False(t, it.HasNext())
	assert.Panics(t, func() {
		it.Next()
	})
	mock.AssertExpectationsForObjects(t, mockList, mockIterator)
}

func (suite *EContentsListIteratorTestSuite) TestIteratorManyFilled() {
	t := suite.T()
	it := suite.it
	mockList := NewMockEList(t)
	mockIterator := &MockEIterator{}
	mockResult := NewMockEObject(t)
	suite.mockObject.EXPECT().EIsSet(suite.mockFeature).Once().Return(true)
	suite.mockObject.EXPECT().EGetResolve(suite.mockFeature, false).Once().Return(mockList)
	suite.mockFeature.EXPECT().IsMany().Once().Return(true)
	mockList.EXPECT().Iterator().Return(mockIterator).Once()
	mockIterator.EXPECT().HasNext().Return(true).Once()
	mockIterator.EXPECT().Next().Return(mockResult).Once()
	assert.True(t, it.HasNext())
	assert.Equal(t, mockResult, it.Next())
	mock.AssertExpectationsForObjects(t, mockList, mockIterator, mockResult)
}

func TestEContentsListIterator(t *testing.T) {
	suite.Run(t, &EContentsListIteratorTestSuite{})
}
