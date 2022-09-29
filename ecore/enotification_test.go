package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventTypeString(t *testing.T) {
	assert.Equal(t, "ADD", ADD.String())
}
