package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockEObjectList_GetUnResolvedList(t *testing.T) {
	o := &MockEObjectList{}
	l := &MockEList{}
	o.On("GetUnResolvedList").Once().Return(l)
	o.On("GetUnResolvedList").Once().Return(func() EList {
		return l
	})
	assert.Equal(t, l, o.GetUnResolvedList())
	assert.Equal(t, l, o.GetUnResolvedList())
	mock.AssertExpectationsForObjects(t, o, l)
}
