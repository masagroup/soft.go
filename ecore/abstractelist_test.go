package ecore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbstractEList_Size(t *testing.T) {
	require.Panics(t, func() {
		l := &AbstractEList{}
		l.Size()
	})
}

func TestAbstractEList_ToArray(t *testing.T) {
	require.Panics(t, func() {
		l := &AbstractEList{}
		l.ToArray()
	})
}
