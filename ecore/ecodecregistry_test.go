// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResoureCodecRegistrySingleton(t *testing.T) {
	r := GetCodecRegistry()
	require.NotNil(t, r)
	assert.Same(t, r, GetCodecRegistry())
	// the codecs of the extensions are registered by the codec packages
	assert.Nil(t, r.GetExtensionToCodecMap()["ecore"])
	assert.Nil(t, r.GetExtensionToCodecMap()["xml"])
}

func TestResoureCodecRegistrySingletonGetCodec(t *testing.T) {
	r := GetCodecRegistry()
	assert.IsType(t, &NoCodec{}, r.GetCodec(NewURI("memory://model.xml")))
	assert.Nil(t, r.GetCodec(NewURI("*.xml")))
}

func TestRegisterCodec(t *testing.T) {
	r := GetCodecRegistry()
	codec := NewMockECodec(t)
	RegisterCodec("test", codec)
	t.Cleanup(func() { delete(r.GetExtensionToCodecMap(), "test") })
	assert.Same(t, codec, r.GetCodec(NewURI("model.test")))

	RegisterProtocolCodec("testprotocol", codec)
	t.Cleanup(func() { delete(r.GetProtocolToCodecMap(), "testprotocol") })
	assert.Same(t, codec, r.GetCodec(NewURI("testprotocol://model.xml")))
}

// Concurrent first callers must all get the one registry, built once.
func TestGetCodecRegistryConcurrentInitialization(t *testing.T) {
	previous := resourceCodecRegistryInstance
	t.Cleanup(func() {
		resourceCodecRegistryInstance = previous
		resourceCodecRegistryOnce = sync.Once{}
		resourceCodecRegistryOnce.Do(func() {})
	})
	resourceCodecRegistryInstance = nil
	resourceCodecRegistryOnce = sync.Once{}
	registries := make(chan ECodecRegistry, 100)
	var group sync.WaitGroup
	for range cap(registries) {
		group.Add(1)
		go func() {
			defer group.Done()
			registries <- GetCodecRegistry()
		}()
	}
	group.Wait()
	close(registries)

	want := <-registries
	assert.NotNil(t, want)
	for registry := range registries {
		assert.Same(t, want, registry)
	}
}
