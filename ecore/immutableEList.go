package ecore

import "strconv"

type immutableEList struct {
	data []interface{}
}

// NewImmutableEList return a new ImmutableEList
func NewImmutableEList(data []interface{}) *immutableEList {
	return &immutableEList{data: data}
}

func (l *immutableEList) Add(elem interface{}) bool {
	panic("Immutable list can't be modified")
}

func (l *immutableEList) AddAll(list EList) bool {
	panic("Immutable list can't be modified")
}

func (l *immutableEList) Insert(index int, elem interface{}) bool {
	panic("Immutable list can't be modified")
}

func (l *immutableEList) InsertAll(index int, list EList) bool {
	panic("Immutable list can't be modified")
}

func (l *immutableEList) MoveObject(newIndex int, elem interface{}) {
	panic("Immutable list can't be modified")
}

func (l *immutableEList) Move(oldIndex, newIndex int) interface{} {
	panic("Immutable list can't be modified")
}

// Get an element of the array
func (l *immutableEList) Get(index int) interface{} {
	if index < 0 || index >= l.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(l.Size()))
	}
	return l.data[index]
}

func (l *immutableEList) Set(index int, elem interface{}) interface{} {
	panic("Immutable list can't be modified")
}

func (l *immutableEList) RemoveAt(index int) interface{} {
	panic("Immutable list can't be modified")
}

func (l *immutableEList) Remove(elem interface{}) bool {
	panic("Immutable list can't be modified")
}

// Size count the number of element in the array
func (l *immutableEList) Size() int {
	return len(l.data)
}

func (l *immutableEList) Clear() {
	panic("Immutable list can't be modified")
}

// Empty return true if the array contains 0 element
func (l *immutableEList) Empty() bool {
	return l.Size() == 0
}

// Contains return if an array contains or not an element
func (l *immutableEList) Contains(elem interface{}) bool {
	return l.IndexOf(elem) != -1
}

// IndexOf return the index on an element in an array, else return -1
func (l *immutableEList) IndexOf(elem interface{}) int {
	for i, value := range l.data {
		if value == elem {
			return i
		}
	}
	return -1
}

// Iterator through the array
func (l *immutableEList) Iterator() EIterator {
	return &listIterator{list: l}
}

// ToArray convert to array
func (l *immutableEList) ToArray() []interface{} {
	return l.data
}

func (l *immutableEList) GetUnResolvedList() EList {
	return l
}
