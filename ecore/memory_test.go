package ecore

import (
	"testing"

	"github.com/OneOfOne/go-utils/memory"
)

func TestMemorySizes(t *testing.T) {
	t.Log(memory.Sizeof(GetPackage()))
	t.Log(memory.Sizeof(newEClassExt()))
	t.Log(memory.Sizeof(&CompactEObjectImpl{}))
	t.Log(memory.Sizeof(&CompactEObjectContainer{}))
	t.Log(memory.Sizeof(&BasicEObjectImpl{}))
}
