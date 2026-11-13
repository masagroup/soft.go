// *****************************************************************************
// Copyright(c) 2026 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// compactTestFields is a block of 32 fields, the most a FieldsHeader holds.
type compactTestFields struct {
	FieldsHeader
	v [32]int
}

// wideTestFields is a block of 64 fields, the most a WideFieldsHeader holds.
type wideTestFields struct {
	WideFieldsHeader
	v [64]int
}

// newTestLevel returns a level of block B with n int fields, the field i at v[i].
func newTestLevel[B any](n int, offsetOfV uintptr) *Level {
	l := NewLevel[B](hybridTestClass, 0)
	for i := range n {
		Value[int](l, i, i, offsetOfV+uintptr(i)*unsafe.Sizeof(int(0)), Plain)
	}
	return l
}

func TestFieldsHeader_Sizes(t *testing.T) {
	assert.Equal(t, uintptr(8), unsafe.Sizeof(FieldsHeader{}))
	assert.Equal(t, uintptr(16), unsafe.Sizeof(WideFieldsHeader{}))
	assert.False(t, NewLevel[compactTestFields](hybridTestClass, 0).wide)
	assert.True(t, NewLevel[wideTestFields](hybridTestClass, 0).wide)
}

// A level with a FieldsHeader holds 32 fields; the 33rd needs a WideFieldsHeader.
func TestLevel_CompactHeaderLimit(t *testing.T) {
	type blockOf33 struct {
		FieldsHeader
		v [33]int
	}
	assert.NotPanics(t, func() { newTestLevel[compactTestFields](32, unsafe.Offsetof(compactTestFields{}.v)) })
	assert.Panics(t, func() { newTestLevel[blockOf33](33, unsafe.Offsetof(blockOf33{}.v)) })
	assert.NotPanics(t, func() { newTestLevel[wideTestFields](64, unsafe.Offsetof(wideTestFields{}.v)) })
}

// testHeaderBits checks claim, publish, published and init on a block of l, whose last field has
// bit last.
func testHeaderBits(t *testing.T, l *Level, last uint64) {
	t.Helper()
	b := l.newBlock()
	for _, bit := range []uint64{1, last} {
		assert.True(t, l.claim(b, bit), "first claim")
		assert.False(t, l.claim(b, bit), "claimed already")
		assert.False(t, b.has(bit), "claimed, not published")
		assert.Zero(t, l.published(b)&bit, "a claim is not published")
		b.publish(bit)
		assert.True(t, b.has(bit))
	}
	assert.Equal(t, 1|last, l.published(b), "the published bits only")

	other := l.newBlock()
	l.init(other, 1|last)
	assert.Equal(t, 1|last, l.published(other))
	assert.True(t, other.has(last))
	assert.False(t, l.claim(other, last), "init claims")
	assert.True(t, l.claim(other, 2), "init claims only its mask")
}

func TestLevel_CompactHeaderBits(t *testing.T) {
	l := newTestLevel[compactTestFields](32, unsafe.Offsetof(compactTestFields{}.v))
	testHeaderBits(t, l, 1<<31)
}

func TestLevel_WideHeaderBits(t *testing.T) {
	l := newTestLevel[wideTestFields](64, unsafe.Offsetof(wideTestFields{}.v))
	testHeaderBits(t, l, 1<<63)
}

// without copies the published fields but those cleared, whatever the header.
func TestLevel_WithoutKeepsTheFields(t *testing.T) {
	for _, tc := range []struct {
		name string
		l    *Level
		last int
	}{
		{"compact", newTestLevel[compactTestFields](32, unsafe.Offsetof(compactTestFields{}.v)), 31},
		{"wide", newTestLevel[wideTestFields](64, unsafe.Offsetof(wideTestFields{}.v)), 63},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := tc.l
			value := func(b *FieldsHeader, i int) *int {
				return at[int](b, l.fields[i].info().offset)
			}
			b := l.newBlock()
			for _, i := range []int{0, 5, tc.last} {
				require.True(t, l.claim(b, 1<<i))
				*value(b, i) = 100 + i
				b.publish(1 << i)
			}
			next := l.without(b, 1<<5)
			assert.Equal(t, uint64(1|1<<tc.last), l.published(next))
			assert.Equal(t, 100, *value(next, 0))
			assert.Equal(t, 100+tc.last, *value(next, tc.last))
			assert.Zero(t, *value(next, 5), "cleared")
			assert.True(t, l.claim(next, 1<<5), "a cleared field can be claimed again")
		})
	}
}

// Of the goroutines claiming a field at once, one only gets it.
func TestLevel_ConcurrentClaims(t *testing.T) {
	for _, tc := range []struct {
		name string
		l    *Level
		bit  uint64
	}{
		{"compact", newTestLevel[compactTestFields](32, unsafe.Offsetof(compactTestFields{}.v)), 1 << 31},
		{"wide", newTestLevel[wideTestFields](64, unsafe.Offsetof(wideTestFields{}.v)), 1 << 63},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for range 200 {
				b := tc.l.newBlock()
				var won atomic.Int32
				var wg sync.WaitGroup
				for range 8 {
					wg.Go(func() {
						if tc.l.claim(b, tc.bit) {
							won.Add(1)
							b.publish(tc.bit)
						} else {
							b.await(tc.bit)
						}
					})
				}
				wg.Wait()
				assert.Equal(t, int32(1), won.Load())
				assert.True(t, b.has(tc.bit))
			}
		})
	}
}
