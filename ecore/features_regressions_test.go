package ecore

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type featuresRegressionsModel struct {
	ePackage EPackage
	eClass   EClass
	eInt     EDataType
	items    EAttribute // many, not unsettable
}

func newFeaturesRegressionsModel() *featuresRegressionsModel {
	f := GetFactory()
	m := &featuresRegressionsModel{ePackage: f.CreateEPackage()}
	m.ePackage.SetEFactoryInstance(f.CreateEFactory())
	m.eInt = f.CreateEDataType()
	m.eInt.SetName("Int")
	m.eInt.SetInstanceTypeName("int")
	m.ePackage.GetEClassifiers().Add(m.eInt)
	m.eClass = f.CreateEClass()
	m.eClass.SetName("C")
	m.ePackage.GetEClassifiers().Add(m.eClass)
	m.items = f.CreateEAttribute()
	m.items.SetName("items")
	m.items.SetEType(GetPackage().GetEString())
	m.items.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	m.eClass.GetEStructuralFeatures().Add(m.items)
	return m
}

// an invalid default literal gives the default value of the data type, at each read
func TestEStructuralFeatureExt_InvalidDefaultLiteral(t *testing.T) {
	m := newFeaturesRegressionsModel()
	count := GetFactory().CreateEAttribute()
	count.SetName("count")
	count.SetEType(m.eInt)
	count.SetDefaultValueLiteral("1.5")
	m.eClass.GetEStructuralFeatures().Add(count)
	o := m.ePackage.GetEFactoryInstance().Create(m.eClass)
	require.NotPanics(t, func() { assert.Equal(t, 0, o.EGet(count)) })
	count.SetDefaultValueLiteral("7")
	assert.Equal(t, 7, o.EGet(count))
}

// the default value is read by concurrent goroutines (checked with -race)
func TestEStructuralFeatureExt_ConcurrentDefaultValue(t *testing.T) {
	m := newFeaturesRegressionsModel()
	count := GetFactory().CreateEAttribute()
	count.SetName("count")
	count.SetEType(GetPackage().GetEInt())
	count.SetDefaultValueLiteral("7")
	m.eClass.GetEStructuralFeatures().Add(count)
	o := m.ePackage.GetEFactoryInstance().Create(m.eClass)
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.Equal(t, 7, o.EGet(count))
		}()
	}
	wg.Wait()
}

// an empty many feature is not set, even once its list was read
func TestEmptyManyFeatureIsNotSet(t *testing.T) {
	m := newFeaturesRegressionsModel()
	reflective := NewReflectiveEObjectImpl()
	reflective.SetEClass(m.eClass)
	store := NewEStoreEObjectImpl(false)
	store.SetEClass(m.eClass)
	for name, o := range map[string]EObject{
		"dynamic":    m.ePackage.GetEFactoryInstance().Create(m.eClass),
		"reflective": reflective,
		"store":      store,
	} {
		_ = o.EGet(m.items)
		assert.False(t, o.EIsSet(m.items), name)
		assert.True(t, Equals(o, m.ePackage.GetEFactoryInstance().Create(m.eClass)), name)
		o.EGet(m.items).(EList).Add("a")
		assert.True(t, o.EIsSet(m.items), name)
	}
}

// a renamed classifier is found by its new name only
func TestEPackageExt_GetEClassifierAfterRename(t *testing.T) {
	m := newFeaturesRegressionsModel()
	assert.Equal(t, m.eClass, m.ePackage.GetEClassifier("C"))
	m.eClass.SetName("D")
	assert.Equal(t, m.eClass, m.ePackage.GetEClassifier("D"))
	assert.Nil(t, m.ePackage.GetEClassifier("C"))
}

// the values of a reflective object follow the features of its new class
func TestReflectiveEObjectImpl_SetEClassResizes(t *testing.T) {
	f := GetFactory()
	m := newFeaturesRegressionsModel()
	big := f.CreateEClass()
	big.SetName("Big")
	m.ePackage.GetEClassifiers().Add(big)
	big.GetESuperTypes().Add(m.eClass)
	extra := f.CreateEAttribute()
	extra.SetName("extra")
	extra.SetEType(GetPackage().GetEString())
	big.GetEStructuralFeatures().Add(extra)

	o := NewReflectiveEObjectImpl()
	o.SetEClass(m.eClass)
	o.EGet(m.items).(EList).Add("a")
	o.SetEClass(big)
	require.NotPanics(t, func() { o.ESet(extra, "x") })
	assert.Equal(t, "x", o.EGet(extra))
	assert.Equal(t, []any{"a"}, o.EGet(m.items).(EList).ToArray())
}

// a file uri without path is an error, not a panic
func TestFileURIHandler_EmptyPath(t *testing.T) {
	r := NewEResourceImpl()
	r.SetURI(NewURI(""))
	require.NotPanics(t, func() { r.Load() })
	assert.Equal(t, 1, r.GetErrors().Size())
}
