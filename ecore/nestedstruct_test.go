package ecore

import (
	"testing"
)

type A struct {
	a1 string
	a2 string
	a3 int
	a4 int
}

func makeA() A {
	return A{
		a1: "a1",
		a2: "a2",
		a3: 1,
		a4: 2,
	}
}

type B struct {
	A
	b1 string
	b2 string
	b3 int
	b4 int
}

func makeB() B {
	return B{
		A:  makeA(),
		b1: "b1",
		b2: "b2",
		b3: 1,
		b4: 2,
	}
}

type C struct {
	B
	c1 string
	c2 string
	c3 int
	c4 int
}

func makeC() C {
	return C{
		B:  makeB(),
		c1: "c1",
		c2: "c2",
		c3: 1,
		c4: 2,
	}
}

func BenchmarkStructNested(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := C{
			B: B{
				A: A{a1: "a1",
					a2: "a2",
					a3: 1,
					a4: 2},
				b1: "b1",
				b2: "b2",
				b3: 1,
				b4: 2,
			},
			c1: "c1",
			c2: "c2",
			c3: 1,
			c4: 2,
		}
		c.c4 += 1
	}
}

func BenchmarkStructInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := C{}
		c.a1 = "a1"
		c.a2 = "a2"
		c.a3 = 1
		c.a4 = 2
		c.b1 = "b1"
		c.b2 = "b2"
		c.b3 = 1
		c.b4 = 2
		c.c1 = "c1"
		c.c2 = "c2"
		c.c3 = 1
		c.c4 = 2
		c.c4 += 1
	}
}

func BenchmarkStructMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := makeC()
		c.c4 += 1
	}
}
