package ecore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockECacheProvider_IsCache(t *testing.T) {
	m := NewMockECacheProvider(t)
	r := NewMockRun(t)
	m.EXPECT().IsCache().Return(true).Run(func() { r.Run() }).Once()
	m.EXPECT().IsCache().RunAndReturn(func() bool { return true }).Once()
	m.EXPECT().IsCache()
	require.True(t, m.IsCache())
	require.True(t, m.IsCache())
	require.Panics(t, func() {
		m.IsCache()
	})
}

func TestMockECacheProvider_SetCache(t *testing.T) {
	m := NewMockECacheProvider(t)
	r := NewMockRun(t, true)
	m.EXPECT().SetCache(true).Return().Run(func(a bool) { r.Run(a) }).Once()
	m.EXPECT().SetCache(true).RunAndReturn(func(a bool) {}).Once()
	m.SetCache(true)
	m.SetCache(true)
}
