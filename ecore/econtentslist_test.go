package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestEContentsListNoFeatures(t *testing.T) {
	mockObject := &MockEObject{}
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
	suite.mockObject = &MockEObject{}
	suite.mockFeature = &MockEStructuralFeature{}
	suite.l = newEContentsList(suite.mockObject, NewImmutableEList([]interface{}{suite.mockFeature}), false)
}

func (suite *EContentsListTestSuite) AfterTest() {
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)
}

func (suite *EContentsListTestSuite) TestEmpty() {
	// single not set
	suite.mockObject.On("EIsSet", suite.mockFeature).Return(false).Once()
	assert.True(suite.T(), suite.l.Empty())
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)

	// single set
	value := struct{}{}
	suite.mockObject.On("EIsSet", suite.mockFeature).Return(true).Once()
	suite.mockObject.On("EGetResolve", suite.mockFeature, false).Return(value).Once()
	suite.mockFeature.On("IsMany").Return(false).Once()
	assert.False(suite.T(), suite.l.Empty())

	// many
	mockList := &MockEList{}
	suite.mockObject.On("EIsSet", suite.mockFeature).Return(true).Once()
	suite.mockObject.On("EGetResolve", suite.mockFeature, false).Return(mockList).Once()
	suite.mockFeature.On("IsMany").Return(true).Once()
	mockList.On("Empty").Return(false).Once()
	assert.False(suite.T(), suite.l.Empty())
}

func (suite *EContentsListTestSuite) TestSize() {
	// single not set
	suite.mockObject.On("EIsSet", suite.mockFeature).Return(false).Once()
	assert.Equal(suite.T(), 0, suite.l.Size())
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)

	// single set
	value := struct{}{}
	suite.mockObject.On("EIsSet", suite.mockFeature).Return(true).Once()
	suite.mockObject.On("EGetResolve", suite.mockFeature, false).Return(value).Once()
	suite.mockFeature.On("IsMany").Return(false).Once()
	assert.Equal(suite.T(), 1, suite.l.Size())
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)

	// many
	l := NewImmutableEList([]interface{}{struct{}{}})
	suite.mockObject.On("EIsSet", suite.mockFeature).Return(true).Once()
	suite.mockObject.On("EGetResolve", suite.mockFeature, false).Return(l).Once()
	suite.mockFeature.On("IsMany").Return(true).Once()
	assert.Equal(suite.T(), 1, suite.l.Size())
	mock.AssertExpectationsForObjects(suite.T(), suite.mockObject, suite.mockFeature)
}

func (suite *EContentsListTestSuite) TestContains() {

}

func (suite *EContentsListTestSuite) TestIndexOf() {

}

func (suite *EContentsListTestSuite) TestGet() {

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
	suite.mockObject.On("EIsSet", suite.mockFeature).Once().Return(false)
	assert.False(t, it.HasNext())
	assert.Panics(t, func() {
		it.Next()
	})
}

func (suite *EContentsListIteratorTestSuite) TestIteratorSingle() {
	t := suite.T()
	it := suite.it
	value := struct{}{}
	suite.mockObject.On("EIsSet", suite.mockFeature).Once().Return(true)
	suite.mockObject.On("EGetResolve", suite.mockFeature, false).Once().Return(value)
	suite.mockFeature.On("IsMany").Once().Return(false)
	assert.True(t, it.HasNext())
	assert.Equal(t, value, it.Next())
}

func (suite *EContentsListIteratorTestSuite) TestIteratorManyEmpty() {
	t := suite.T()
	it := suite.it
	mockList := &MockEList{}
	mockIterator := &MockEIterator{}
	suite.mockObject.On("EIsSet", suite.mockFeature).Once().Return(true)
	suite.mockObject.On("EGetResolve", suite.mockFeature, false).Once().Return(mockList)
	suite.mockFeature.On("IsMany").Once().Return(true)
	mockList.On("Iterator").Return(mockIterator).Once()
	mockIterator.On("HasNext").Return(false).Once()
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
	mockList := &MockEList{}
	mockIterator := &MockEIterator{}
	mockResult := &MockEObject{}
	suite.mockObject.On("EIsSet", suite.mockFeature).Once().Return(true)
	suite.mockObject.On("EGetResolve", suite.mockFeature, false).Once().Return(mockList)
	suite.mockFeature.On("IsMany").Once().Return(true)
	mockList.On("Iterator").Return(mockIterator).Once()
	mockIterator.On("HasNext").Return(true).Once()
	mockIterator.On("Next").Return(mockResult).Once()
	assert.True(t, it.HasNext())
	assert.Equal(t, mockResult, it.Next())
	mock.AssertExpectationsForObjects(t, mockList, mockIterator, mockResult)
}

func TestEContentsListIterator(t *testing.T) {
	suite.Run(t, &EContentsListIteratorTestSuite{})
}
