package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoCodec_Decoder(t *testing.T) {
	mockResource := NewMockEResource(t)
	cd := &NoCodec{}
	d := cd.NewDecoder(mockResource, nil, nil)
	d.Decode()
	o, err := d.DecodeObject()
	assert.Nil(t, o)
	assert.Nil(t, err)
}

func TestNoCodec_Encoder(t *testing.T) {
	mockResource := NewMockEResource(t)
	mockObject := NewMockEObject(t)
	cd := &NoCodec{}
	d := cd.NewEncoder(mockResource, nil, nil)
	d.Encode()
	err := d.EncodeObject(mockObject)
	assert.Nil(t, err)
}
