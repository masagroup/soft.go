package ecore

import (
	"testing"

	"github.com/OneOfOne/go-utils/memory"
)

func BenchmarkSizes(b *testing.B) {
	b.Log(memory.Sizeof(GetPackage()))
	b.Log(memory.Sizeof(newEClassExt()))
	b.Log(memory.Sizeof(&CompactEObjectImpl{}))
	b.Log(memory.Sizeof(&CompactEObjectContainer{}))
	b.Log(memory.Sizeof(&BasicEObjectImpl{}))
}
