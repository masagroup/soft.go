package ecore

import (
	"testing"

	"github.com/OneOfOne/go-utils/memory"
)

func BenchmarkSizes(b *testing.B) {
	n := NewEObjectImpl()
	b.Log(memory.Sizeof(n))
}
