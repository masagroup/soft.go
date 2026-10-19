package integration

import (
	"testing"

	"github.com/masagroup/soft.go/test/library"
	"github.com/stretchr/testify/assert"
)

// TestWriterExt checks the hand-written extension of the class Writer (writer_ext.go in
// test/library), which the generated factory creates.
func TestWriterExt(t *testing.T) {
	w := library.GetFactory().CreateWriter()
	w.SetFirstName("First 1")
	w.SetLastName("Last 1")
	assert.Equal(t, "First 1--Last 1", w.GetName())

	w.SetName("First 2--Last 2")
	assert.Equal(t, "First 2", w.GetFirstName())
	assert.Equal(t, "Last 2", w.GetLastName())
}
