package ecore

import (
	"net/url"
	"testing"

	"github.com/OneOfOne/go-utils/memory"
)

type StructWithDefinedMembers struct {
	// interfaces interface{}
	// b          bool
	// l1         EList
	// l2         EList
	proxyURI *url.URL
}

type StructWithUndefined struct {
	interfaces interface{}
	arr        interface{}
	flags      int
}

func BenchmarkSizes(b *testing.B) {
	s := &StructWithUndefined{}
	s.interfaces = s
	s.arr = []interface{}{}
	b.Log(memory.Sizeof(s))

	n := NewEObjectImpl()
	b.Log(memory.Sizeof(n))
}
