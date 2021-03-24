package library

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriterExt(t *testing.T) {
	w := newWriterExt()
	w.SetFirstName("First 1")
	w.SetLastName("Last 1")
	assert.Equal(t, "First 1--Last 1", w.GetName())

	w.SetName("First 2--Last 2")
	assert.Equal(t, "First 2", w.GetFirstName())
	assert.Equal(t, "Last 2", w.GetLastName())
}
