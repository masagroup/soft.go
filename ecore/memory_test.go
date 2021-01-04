package ecore

import (
	"testing"

	"github.com/OneOfOne/go-utils/memory"
)

func BenchmarkSizes(b *testing.B) {
	b.Log(memory.Sizeof(GetPackage()))
}
