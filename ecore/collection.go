package ecore

import "iter"

type Collection interface {

	// Add adds element at the end of the list
	// returns true if element was added
	Add(element any) bool

	// Add adds collection at the end of the list
	// returns true if collection was added to the list
	AddAll(collection Collection) bool

	// Remove removes element from the list
	// returns true if element was effectively removed
	Remove(element any) bool

	// RemoveAll removes collection from the list
	// returns true if collection was effectively removed
	RemoveAll(collection Collection) bool

	// Size returns the size of the list
	Size() int

	// Clear clears the list
	Clear()

	// Empty returns true if list is empty , false otherwise
	Empty() bool

	// Contains returns true if list constians element, false otherwise
	Contains(element any) bool

	// Iterator returns the list pull iterator
	Iterator() EIterator

	// Iterator returns the list push iterator
	All() iter.Seq[any]

	// ToArray returns the list array
	ToArray() []any
}
