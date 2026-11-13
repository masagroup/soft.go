// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"fmt"
	"math/bits"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// HybridRoot is the root of the hybrid classes, as soft.gen generates them with the property
// hybridFields: store-backed classes whose getters read typed fields, a cache of the values of
// the store. The descriptors of the fields (Value, List) and the Layout of a class ask it what
// depends on how it keeps its values. HybridEObjectImpl is the root over EStoreEObjectImpl; other
// roots, such as the objects of a view, implement the same methods.
type HybridRoot interface {
	// EHybridBase returns the address of the root in its object: the levels of the fields of the
	// classes find their blocks at their offsets from it.
	EHybridBase() unsafe.Pointer
	// ELoadFeature returns the value of the single-valued feature featureID, feature, for its
	// typed getter, which publishes it: the value of the store, the default value of the feature
	// when it has none, a proxy resolved.
	ELoadFeature(featureID int, feature EStructuralFeature) any
	// ELoadList returns the list of the many-valued feature featureID for its typed getter.
	ELoadList(featureID int) EList
	// ESetField sets the plain feature featureID to value for its typed setter.
	ESetField(featureID int, value any)
	// ESetsUnchanged reports whether setting a feature to its current value changes nothing, so
	// that the typed setter skips it.
	ESetsUnchanged() bool
	// EBaseGetFromID, EBaseSetFromID and EBaseIsSetFromID are the reflective methods of the root
	// without the typed fields: those of the features the fields do not hold.
	EBaseGetFromID(featureID int, resolve bool) any
	EBaseSetFromID(featureID int, value any)
	EBaseIsSetFromID(featureID int) bool
}

// LayoutProvider is implemented by the hybrid classes: ELayout returns the Layout of the typed
// fields of the class, nil for a class without field.
type LayoutProvider interface {
	ELayout() *Layout
}

// FieldsHeader is the header of a block of typed fields, its first field: a bit per field,
// claimed by the one writing it, then published.
type FieldsHeader struct {
	claimed atomic.Uint64
	mask    atomic.Uint64
}

func (h *FieldsHeader) has(bit uint64) bool { return h.mask.Load()&bit != 0 }

// claim claims the field of bit for its caller, which writes it then publishes it. It returns
// false when the field is claimed already.
func (h *FieldsHeader) claim(bit uint64) bool { return h.claimed.Or(bit)&bit == 0 }

func (h *FieldsHeader) publish(bit uint64) { h.mask.Or(bit) }

func (h *FieldsHeader) published() uint64 { return h.mask.Load() }

// init publishes the fields of mask in a block nobody else holds yet.
func (h *FieldsHeader) init(mask uint64) {
	h.claimed.Store(mask)
	h.mask.Store(mask)
}

// await returns once the field of bit, claimed by another caller, is published.
func (h *FieldsHeader) await(bit uint64) {
	for h.mask.Load()&bit == 0 {
		runtime.Gosched()
	}
}

// BlockRef holds the block of the typed fields of a class of a hybrid object, a struct whose first
// field is a FieldsHeader, by its header. A published field is never written again: a change
// replaces the block, so that a reader holding it reads the values that were current when it took
// it.
type BlockRef struct {
	p atomic.Pointer[FieldsHeader]
}

// block returns the block a load publishes in: taken before the read, so that a value read
// before the fields are dropped is not published after.
func (r *BlockRef) block(l *Level) *FieldsHeader {
	if b := r.p.Load(); b != nil {
		return b
	}
	b := l.newBlock()
	if !r.p.CompareAndSwap(nil, b) {
		if current := r.p.Load(); current != nil {
			return current
		}
	}
	return b
}

// clear drops the block.
func (r *BlockRef) clear() {
	if r.p.Load() != nil {
		r.p.Store(nil)
	}
}

// replace replaces the block by a copy without the fields of clear, set by set when not nil: a
// load that read the store before the change publishes in the replaced block.
func (r *BlockRef) replace(l *Level, clear uint64, set func(*FieldsHeader)) {
	for {
		b := r.p.Load()
		if b == nil && set == nil {
			return
		}
		next := l.without(b, clear)
		if set != nil {
			set(next)
		} else if next.published() == 0 {
			next = nil
		}
		if r.p.CompareAndSwap(b, next) {
			return
		}
	}
}

// Fields holds the block B of the typed fields of a class of a hybrid object: the typed view of
// a BlockRef, which the generated getters read.
type Fields[B any] struct {
	ref BlockRef
}

// Published returns the block if the field of bit is published in it, nil otherwise.
func (f *Fields[B]) Published(bit uint64) *B {
	if h := f.ref.p.Load(); h != nil && h.has(bit) {
		return (*B)(unsafe.Pointer(h))
	}
	return nil
}

// Ref returns the block reference of f, which the descriptors of its fields use.
func (f *Fields[B]) Ref() *BlockRef { return &f.ref }

// FieldFlags describe a field.
type FieldFlags uint8

const (
	// Reference marks the field of a reference: the reflective read resolving proxies reads it,
	// and a proxy set to it is not published.
	Reference FieldFlags = 1 << iota
	// Plain marks the field of a single-valued feature that is neither a containment nor
	// bidirectional: its typed setter, and ESetFromID, set it without the reflective lookups.
	Plain
	// Nillable marks the field of a pointer type: its nil value is read as an untyped nil.
	Nillable
	// list marks the field of a many-valued feature.
	list
)

// field is a typed field of a class, whatever its type: what the reflective methods of
// the hybrid root use.
type field interface {
	info() *fieldInfo
	// valueOf returns the value of the field, loaded if it is not published.
	valueOf(o HybridRoot, r *BlockRef) any
	// setOf sets the field to value with its typed setter, false if value is not of its type.
	setOf(o HybridRoot, r *BlockRef, value any) bool
	// replaceOf publishes value, just set to the feature, in the field.
	replaceOf(r *BlockRef, value any)
	// isSetOf reports whether the feature is set, and whether the field answers it.
	isSetOf(o HybridRoot, r *BlockRef) (isSet bool, answered bool)
	// copy copies the field from src to dst.
	copy(dst, src *FieldsHeader)
}

// fieldInfo is what all fields have.
type fieldInfo struct {
	level  *Level
	id     int     // the feature ID in the class
	bit    uint64  // the bit of the field in its block
	index  int     // the index of the field in its level
	offset uintptr // the offset of the field in the block
	flags  FieldFlags
}

func (i *fieldInfo) info() *fieldInfo { return i }

// at returns the field of the block of header h.
func at[T any](h *FieldsHeader, offset uintptr) *T {
	return (*T)(unsafe.Add(unsafe.Pointer(h), offset))
}

// Level is the typed fields of one class of a hybrid object, in a block: a struct whose first
// field is a FieldsHeader, held by a Fields at offset from the hybrid root of the object (HybridRoot.EHybridBase).
type Level struct {
	class    func() EClass
	offset   uintptr
	newBlock func() *FieldsHeader
	fields   []field
	features func() []EStructuralFeature
}

// NewLevel returns the level of the typed fields of class, in a block B held by a Fields at offset
// from the hybrid root of the objects of the class. B must have FieldsHeader as its first
// field.
func NewLevel[B any](class func() EClass, offset uintptr) *Level {
	if t := reflect.TypeFor[B](); t.Kind() != reflect.Struct || t.NumField() == 0 || t.Field(0).Type != reflect.TypeFor[FieldsHeader]() {
		panic(fmt.Sprintf("ecore: %v does not have FieldsHeader as its first field", t))
	}
	l := &Level{class: class, offset: offset, newBlock: func() *FieldsHeader { return (*FieldsHeader)(unsafe.Pointer(new(B))) }}
	l.features = sync.OnceValue(func() []EStructuralFeature {
		c := l.class()
		features := make([]EStructuralFeature, len(l.fields))
		for i, f := range l.fields {
			features[i] = c.GetEStructuralFeature(f.info().id)
		}
		return features
	})
	return l
}

// refOf returns the block reference of the level in the object of o.
func (l *Level) refOf(o HybridRoot) *BlockRef {
	return (*BlockRef)(unsafe.Add(o.EHybridBase(), l.offset))
}

func (l *Level) add(f field, i *fieldInfo) {
	if i.index != len(l.fields) || i.index >= 64 {
		panic(fmt.Sprintf("ecore: field %d of a level registered as bit %d", len(l.fields), i.index))
	}
	l.fields = append(l.fields, f)
}

// without returns a new block holding the published fields of b but those of clear: a field is
// read only once published, never written again then. An empty block for a nil one.
func (l *Level) without(b *FieldsHeader, clear uint64) *FieldsHeader {
	next := l.newBlock()
	if b == nil {
		return next
	}
	mask := b.published() &^ clear
	for m := mask; m != 0; m &= m - 1 {
		l.fields[bits.TrailingZeros64(m)].copy(next, b)
	}
	next.init(mask)
	return next
}

// ValueField is the typed field of a single-valued feature of a comparable type T.
type ValueField[T comparable] struct {
	fieldInfo
}

// Value registers the field of the single-valued feature id at bit in l, at offset in the block.
func Value[T comparable](l *Level, id int, bit int, offset uintptr, flags FieldFlags) *ValueField[T] {
	f := &ValueField[T]{fieldInfo{level: l, id: id, bit: 1 << bit, index: bit, offset: offset, flags: flags}}
	l.add(f, &f.fieldInfo)
	return f
}

func (d *ValueField[T]) copy(dst, src *FieldsHeader) {
	*at[T](dst, d.offset) = *at[T](src, d.offset)
}

// isProxy reports whether value is a proxy object.
func isProxy(value any) bool {
	object, _ := value.(EObject)
	return object != nil && object.EIsProxy()
}

// Load reads the value of the field from the store and publishes it in r: the cold path of the
// typed getter, and of the getter of a reference whose published object is a proxy.
//
// A reference is resolved at each read, as by the reflective read of EMF: a proxy, which did not
// resolve, is not published, and is resolved again by the next read. A published object may
// become a proxy after it was published, as an object of an unloaded resource, or one a decoder
// gives its proxy URI after setting it: its field is replaced by the object resolved.
func (d *ValueField[T]) Load(o HybridRoot, r *BlockRef) T {
	b := r.block(d.level)
	loaded := o.ELoadFeature(d.id, d.level.features()[d.index])
	value, _ := loaded.(T)
	reference := d.flags&Reference != 0
	if reference && isProxy(loaded) {
		return value
	}
	if b.claim(d.bit) {
		*at[T](b, d.offset) = value
		b.publish(d.bit)
	} else if reference && b.has(d.bit) && isProxy(*at[T](b, d.offset)) {
		r.replace(d.level, d.bit, func(next *FieldsHeader) {
			*at[T](next, d.offset) = value
			next.init(next.published() | d.bit)
		})
	}
	return value
}

// Get returns the value of the field, loaded if it is not published, or if it is a reference
// whose object became a proxy.
func (d *ValueField[T]) Get(o HybridRoot, r *BlockRef) T {
	if b := r.p.Load(); b != nil && b.has(d.bit) {
		value := *at[T](b, d.offset)
		if d.flags&Reference == 0 || !isProxy(value) {
			return value
		}
	}
	return d.Load(o, r)
}

// Set sets the feature to value: nothing when the published value is value, the object has a
// store and nothing is notified; else through the store without the reflective lookups.
func (d *ValueField[T]) Set(o HybridRoot, r *BlockRef, value T) {
	if b := r.p.Load(); b != nil && b.has(d.bit) && *at[T](b, d.offset) == value && o.ESetsUnchanged() {
		return
	}
	o.ESetField(d.id, value)
}

func (d *ValueField[T]) valueOf(o HybridRoot, r *BlockRef) any {
	value := d.Get(o, r)
	if d.flags&Nillable != 0 {
		var zero T
		if value == zero {
			return nil
		}
	}
	return value
}

func (d *ValueField[T]) setOf(o HybridRoot, r *BlockRef, value any) bool {
	typed, ok := value.(T)
	if !ok {
		return false
	}
	d.Set(o, r, typed)
	return true
}

func (d *ValueField[T]) replaceOf(r *BlockRef, value any) {
	typed, ok := value.(T)
	if ok && d.flags&Reference != 0 {
		object, _ := value.(EObject)
		ok = object != nil && !object.EIsProxy()
	}
	if !ok {
		r.replace(d.level, d.bit, nil)
		return
	}
	// a field nobody claimed is written in place: a load that read the store before the set
	// fails to claim it, and returns its value without publishing it
	if b := r.p.Load(); b != nil && b.claim(d.bit) {
		*at[T](b, d.offset) = typed
		b.publish(d.bit)
		return
	}
	r.replace(d.level, d.bit, func(b *FieldsHeader) {
		*at[T](b, d.offset) = typed
		b.init(b.published() | d.bit)
	})
}

func (d *ValueField[T]) isSetOf(HybridRoot, *BlockRef) (bool, bool) { return false, false }

// ListField is the typed field of a many-valued feature: its list, kept by the changes of the
// values, as it reads the current revision of the object.
type ListField struct {
	fieldInfo
}

// List registers the field of the many-valued feature id at bit in l, at offset in the block.
func List(l *Level, id int, bit int, offset uintptr) *ListField {
	f := &ListField{fieldInfo{level: l, id: id, bit: 1 << bit, index: bit, offset: offset, flags: list}}
	l.add(f, &f.fieldInfo)
	return f
}

func (d *ListField) copy(dst, src *FieldsHeader) {
	*at[EList](dst, d.offset) = *at[EList](src, d.offset)
}

// Load returns the list of the feature, made once per block: the cold path of the typed getter.
// Every reader gets the same list.
func (d *ListField) Load(o HybridRoot, r *BlockRef) EList {
	b := r.block(d.level)
	if b.claim(d.bit) {
		*at[EList](b, d.offset) = o.ELoadList(d.id)
		b.publish(d.bit)
	} else {
		b.await(d.bit)
	}
	return *at[EList](b, d.offset)
}

// Get returns the list of the feature.
func (d *ListField) Get(o HybridRoot, r *BlockRef) EList {
	if b := r.p.Load(); b != nil && b.has(d.bit) {
		return *at[EList](b, d.offset)
	}
	return d.Load(o, r)
}

func (d *ListField) valueOf(o HybridRoot, r *BlockRef) any { return d.Get(o, r) }

func (d *ListField) setOf(HybridRoot, *BlockRef, any) bool { return false }

func (d *ListField) replaceOf(*BlockRef, any) {}

// isSetOf reports a list set when it is not empty, without asking the store.
func (d *ListField) isSetOf(o HybridRoot, r *BlockRef) (bool, bool) {
	return !d.Get(o, r).Empty(), true
}

// Layout is the typed fields of a hybrid class: its levels, its own and those of the classes it
// extends, and its fields by feature ID.
type Layout struct {
	levels []*Level
	mu     sync.Mutex
	byID   atomic.Pointer[[]field]
}

// NewLayout returns the layout of a class whose own fields are level, extending the class of
// parent (nil for none).
func NewLayout(parent *Layout, level *Level) *Layout {
	l := &Layout{}
	if parent != nil {
		l.levels = append(l.levels, parent.levels...)
	}
	l.levels = append(l.levels, level)
	return l
}

// index builds the fields by feature ID, at first use: the fields register themselves in the
// levels after the layout is made.
func (l *Layout) index() []field {
	l.mu.Lock()
	defer l.mu.Unlock()
	if byID := l.byID.Load(); byID != nil {
		return *byID
	}
	var byID []field
	for _, level := range l.levels {
		for _, f := range level.fields {
			id := f.info().id
			if id >= len(byID) {
				byID = append(byID, make([]field, id+1-len(byID))...)
			}
			byID[id] = f
		}
	}
	l.byID.Store(&byID)
	return byID
}

// field returns the field of featureID, nil for none.
func (l *Layout) field(featureID int) field {
	if l == nil {
		return nil
	}
	var byID []field
	if p := l.byID.Load(); p != nil {
		byID = *p
	} else {
		byID = l.index()
	}
	if featureID < 0 || featureID >= len(byID) {
		return nil
	}
	return byID[featureID]
}

// clearAll drops the blocks of o.
func (l *Layout) clearAll(o HybridRoot) {
	if l == nil {
		return
	}
	for _, level := range l.levels {
		level.refOf(o).clear()
	}
}

// clearFeatures drops the fields of the single-valued features of featureIDs of o.
func (l *Layout) clearFeatures(o HybridRoot, featureIDs []int) {
	if l == nil {
		return
	}
	for _, level := range l.levels {
		var bits uint64
		for _, featureID := range featureIDs {
			if f := l.field(featureID); f != nil && f.info().level == level && f.info().flags&list == 0 {
				bits |= f.info().bit
			}
		}
		if bits != 0 {
			level.refOf(o).replace(level, bits, nil)
		}
	}
}

// fieldOf returns the typed field of featureID and its block reference in o, nil for none.
func (l *Layout) fieldOf(o HybridRoot, featureID int) (field, *BlockRef) {
	if f := l.field(featureID); f != nil {
		return f, f.info().level.refOf(o)
	}
	return nil, nil
}

// EGetFromID returns the value of the feature featureID of o: the one of its typed field, for a
// feature held by one (a reference only when resolving), else the reflective read of soft.go.
func (l *Layout) EGetFromID(o HybridRoot, featureID int, resolve bool) any {
	if f, r := l.fieldOf(o, featureID); f != nil && (resolve || f.info().flags&Reference == 0) {
		return f.valueOf(o, r)
	}
	return o.EBaseGetFromID(featureID, resolve)
}

// ESetFromID sets the feature featureID of o to value with its typed setter, for a plain feature
// given a value of its type, else with the reflective write of soft.go.
func (l *Layout) ESetFromID(o HybridRoot, featureID int, value any) {
	if f, r := l.fieldOf(o, featureID); f != nil && f.info().flags&Plain != 0 && f.setOf(o, r, value) {
		return
	}
	o.EBaseSetFromID(featureID, value)
}

// EIsSetFromID reports whether the feature featureID of o is set: the list of a typed field when
// it is not empty, else what soft.go answers.
func (l *Layout) EIsSetFromID(o HybridRoot, featureID int) bool {
	if f, r := l.fieldOf(o, featureID); f != nil {
		if isSet, answered := f.isSetOf(o, r); answered {
			return isSet
		}
	}
	return o.EBaseIsSetFromID(featureID)
}

// Value returns the value of the feature featureID of o held by a typed field, loaded if it is not
// published, and false when no field holds it.
func (l *Layout) Value(o HybridRoot, featureID int) (any, bool) {
	if f, r := l.fieldOf(o, featureID); f != nil {
		return f.valueOf(o, r), true
	}
	return nil, false
}

// Replace publishes value, just set to the feature featureID of o, in its typed field: the next
// read does not go back to the store.
func (l *Layout) Replace(o HybridRoot, featureID int, value any) {
	if f, r := l.fieldOf(o, featureID); f != nil {
		f.replaceOf(r, value)
	}
}

// ClearAll drops the typed fields of o.
func (l *Layout) ClearAll(o HybridRoot) { l.clearAll(o) }

// ClearFeatures drops the typed fields of the single-valued features of featureIDs of o.
func (l *Layout) ClearFeatures(o HybridRoot, featureIDs []int) { l.clearFeatures(o, featureIDs) }
