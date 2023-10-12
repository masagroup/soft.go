package ecore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLStore_Constructor(t *testing.T) {
	s, err := NewSQLStore("testdata/library.complex.sqlite", NewURI(""), nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
}

func TestSQLStore_SetSingleValue(t *testing.T) {

}
