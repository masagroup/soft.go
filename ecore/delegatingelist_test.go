package ecore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDelegatingEList_ToArray(t *testing.T) {
	d := &AbstractDelegatingEList[*BasicEList]{}
	d.delegate = NewBasicEList([]any{1, 2})
	d.interfaces = d
	require.Equal(t, []any{1, 2}, d.ToArray())
}

func TestDelegatingEList_Get(t *testing.T) {
	d := &AbstractDelegatingEList[*BasicEList]{}
	d.delegate = NewBasicEList([]any{1, 2})
	d.interfaces = d
	require.Equal(t, 2, d.Get(1))
}

func TestDelegatingENotifyingEList_RemoveWithNotification(t *testing.T) {
	d := &AbstractDelegatingENotifyingList[*BasicENotifyingList]{}
	d.interfaces = d
	d.delegate = newBasicENotifyingListFromData([]any{1, 2})
	d.RemoveWithNotification(1, nil)
	require.Equal(t, []any{2}, d.ToArray())
}
