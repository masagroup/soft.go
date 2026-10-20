package sqlite

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/require"
)

func TestMockObject_GetSQLID(t *testing.T) {
	m := NewMockObject(t)
	r := ecore.NewMockRun(t)
	m.EXPECT().GetSQLID().Return(1).Run(func() { r.Run() }).Once()
	m.EXPECT().GetSQLID().RunAndReturn(func() int64 { return 1 }).Once()
	m.EXPECT().GetSQLID()
	require.Equal(t, int64(1), m.GetSQLID())
	require.Equal(t, int64(1), m.GetSQLID())
	require.Panics(t, func() {
		m.GetSQLID()
	})
}

func TestMockObject_SetSQLID(t *testing.T) {
	v := int64(1)
	m := NewMockObject(t)
	r := ecore.NewMockRun(t, v)
	m.EXPECT().SetSQLID(v).Return().Run(func(a int64) { r.Run(a) }).Once()
	m.EXPECT().SetSQLID(v).RunAndReturn(func(a int64) {}).Once()
	m.SetSQLID(v)
	m.SetSQLID(v)
}
