package ecore

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinarySearch(t *testing.T) {
	// lexically ordered
	l := NewImmutableEList([]any{1, 10, 11, 2})
	cmp := func(a any, b any) int {
		return strings.Compare(strconv.Itoa(a.(int)), strconv.Itoa(b.(int)))
	}
	pos, found := BinarySearch(l, 2, cmp)
	require.Equal(t, 3, pos)
	require.True(t, found)
}
