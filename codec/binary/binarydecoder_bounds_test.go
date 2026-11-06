// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package binary

import (
	"bytes"
	"io"
	"runtime"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

// opaqueReader hides the Len method of its reader: the decoder does not know the input length
type opaqueReader struct {
	io.Reader
}

// header returns the signature and the version of a binary resource followed by the given ints
func header(t *testing.T, ints ...int64) []byte {
	var buffer bytes.Buffer
	e := msgpack.NewEncoder(&buffer)
	require.NoError(t, e.EncodeBytes(binarySignature))
	require.NoError(t, e.EncodeInt(int64(binaryVersion)))
	for _, i := range ints {
		require.NoError(t, e.EncodeInt(i))
	}
	return buffer.Bytes()
}

// decodeErrors decodes data in a new resource and returns its errors
func decodeErrors(r io.Reader, options map[string]any) ecore.EList {
	eResource := ecore.NewEResourceSetImpl().CreateResource(ecore.NewURI("file:///bounds.bin"))
	NewDecoder(eResource, r, options).DecodeResource()
	return eResource.GetErrors()
}

// a few bytes announcing 2^31 objects are refused without allocating for them, whether the
// decoder knows the length of the input or not
func TestBinaryDecoder_HugeSize(t *testing.T) {
	data := header(t, 1<<31)
	for name, r := range map[string]func() io.Reader{
		"known length":   func() io.Reader { return bytes.NewReader(data) },
		"unknown length": func() io.Reader { return opaqueReader{bytes.NewReader(data)} },
	} {
		t.Run(name, func(t *testing.T) {
			var errors ecore.EList
			assertSmallAllocation(t, func() { errors = decodeErrors(r(), nil) })
			require.False(t, errors.Empty())
		})
	}
}

// assertSmallAllocation checks that f allocates less than 1 MiB
func assertSmallAllocation(t *testing.T, f func()) {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	f()
	runtime.ReadMemStats(&after)
	assert.Less(t, after.TotalAlloc-before.TotalAlloc, uint64(1<<20))
}

func TestBinaryDecoder_NegativeSize(t *testing.T) {
	errors := decodeErrors(bytes.NewReader(header(t, -5)), nil)
	require.False(t, errors.Empty())
	assert.Contains(t, errors.Get(0).(ecore.EDiagnostic).GetMessage(), "invalid size -5")
}

// an object id is a scalar or an extension: an array, which msgpack allocates from its announced
// length, is refused before it is read
func TestBinaryDecoder_ObjectIDArray(t *testing.T) {
	data := []byte{msgpcode.Array32, 0x7f, 0xff, 0xff, 0xff}
	d := NewDecoder(ecore.NewEResourceImpl(), opaqueReader{bytes.NewReader(data)}, nil)
	var err error
	assertSmallAllocation(t, func() { _, err = d.decodeInterface() })
	require.Error(t, err)
}

// the uuid and ulid extensions have 16 bytes: a larger announced length is refused before
// allocating it
func TestBinaryDecoder_ExtensionLength(t *testing.T) {
	for _, extID := range []int8{ulidExtID, uuidExtID} {
		data := []byte{msgpcode.Ext32, 0xff, 0xff, 0xff, 0xf0, byte(extID)}
		var err error
		assertSmallAllocation(t, func() {
			_, err = msgpack.NewDecoder(opaqueReader{bytes.NewReader(data)}).DecodeInterface()
		})
		require.Error(t, err)
	}
}

// newChain returns a resource of depth objects, each one containing the next
func newChain(t *testing.T, depth int) []byte {
	f := ecore.GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetName("chain")
	ePackage.SetNsURI("http://chain")
	ePackage.SetEFactoryInstance(f.CreateEFactory())
	eClass := f.CreateEClass()
	eClass.SetName("Node")
	ePackage.GetEClassifiers().Add(eClass)
	eReference := f.CreateEReference()
	eReference.SetName("child")
	eReference.SetContainment(true)
	eReference.SetEType(eClass)
	eClass.GetEStructuralFeatures().Add(eReference)

	// the decoder finds the package in the global registry
	ecore.GetPackageRegistry().RegisterPackage(ePackage)
	t.Cleanup(func() { ecore.GetPackageRegistry().UnregisterPackage(ePackage) })

	// the bytes the encoder writes for that chain, written directly: building and encoding a deep
	// chain is quadratic (EResource walks up the containers)
	var buffer bytes.Buffer
	e := msgpack.NewEncoder(&buffer)
	encode := func(values ...any) {
		for _, value := range values {
			require.NoError(t, e.Encode(value))
		}
	}
	encode(binarySignature, binaryVersion, 1)
	// the first object: id, package (id, ns uri, uri id, uri, fragment), class (id, name), feature
	// (id + 1, name)
	encode(0, 0, "http://chain", 0, "", "//", 0, "Node")
	for i := 1; i < depth; i++ {
		// the child feature, then the child: id, package id, class id
		encode(1)
		if i == 1 {
			encode("child")
		}
		encode(i, 0, 0)
	}
	// the end of the features of every object
	for range depth {
		encode(0)
	}
	return buffer.Bytes()
}

func TestBinaryDecoder_MaxDepthOption(t *testing.T) {
	data := newChain(t, 10)
	assert.True(t, decodeErrors(bytes.NewReader(data), map[string]any{OPTION_MAX_DEPTH: 10}).Empty())
	errors := decodeErrors(bytes.NewReader(data), map[string]any{OPTION_MAX_DEPTH: 9})
	require.False(t, errors.Empty())
	assert.Contains(t, errors.Get(0).(ecore.EDiagnostic).GetMessage(), "nested deeper than 9 objects")
}

// a nesting deeper than DefaultMaxDepth is refused with an error instead of growing the stack
func TestBinaryDecoder_DefaultMaxDepth(t *testing.T) {
	data := newChain(t, DefaultMaxDepth+1)
	errors := decodeErrors(bytes.NewReader(data), nil)
	require.False(t, errors.Empty())
	assert.Contains(t, errors.Get(0).(ecore.EDiagnostic).GetMessage(), "nested deeper than")
}
