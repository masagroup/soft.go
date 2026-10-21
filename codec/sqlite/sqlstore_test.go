package sqlite

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

const libraryNSURI = "http:///org/eclipse/emf/examples/library/library.ecore/1.0.0"

// newSQLMockEPackageRegistry returns a mock package registry without factory registered by nsURI:
// the decoder creates the objects with the factory instances of the packages.
func newSQLMockEPackageRegistry(t interface {
	mock.TestingT
	Cleanup(func())
}) *ecore.MockEPackageRegistry {
	r := ecore.NewMockEPackageRegistry(t)
	r.EXPECT().GetFactory(mock.Anything).Return(nil).Maybe()
	return r
}

type dynamicSQLEObjectImpl struct {
	ecore.DynamicEObjectImpl
	sqlID int64
	store ecore.EStore
}

func newDynamicSQLEObjectImpl() *dynamicSQLEObjectImpl {
	o := new(dynamicSQLEObjectImpl)
	o.SetInterfaces(o)
	o.Initialize()
	return o
}

func (o *dynamicSQLEObjectImpl) SetSQLID(sqlID int64) {
	o.sqlID = sqlID
}

func (o *dynamicSQLEObjectImpl) GetSQLID() int64 {
	return o.sqlID
}

func (o *dynamicSQLEObjectImpl) SetEStore(store ecore.EStore) {
	o.store = store
}

func (o *dynamicSQLEObjectImpl) GetEStore() ecore.EStore {
	return o.store
}

type dynamicSQLFactory struct {
	ecore.EFactoryExt
}

func newDynamicSQLFactory() *dynamicSQLFactory {
	eFactory := new(dynamicSQLFactory)
	eFactory.SetInterfaces(eFactory)
	eFactory.Initialize()
	return eFactory
}

func (eFactory *dynamicSQLFactory) Create(eClass ecore.EClass) ecore.EObject {
	if eFactory.GetEPackage() != eClass.GetEPackage() || eClass.IsAbstract() {
		panic("The class '" + eClass.GetName() + "' is not a valid classifier")
	}
	if ecore.IsMapEntry(eClass) {
		eEntry := ecore.NewDynamicEMapEntryImpl()
		eEntry.SetEClass(eClass)
		return eEntry
	} else {
		eObject := newDynamicSQLEObjectImpl()
		eObject.SetEClass(eClass)
		return eObject
	}
}

func closeFile(f *os.File, reported *error) {
	if err := f.Close(); *reported == nil {
		*reported = err
	}
}

func copyFile(src, dest string) (err error) {
	if err = os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return
	}

	d, err := os.Create(dest)
	if err != nil {
		return
	}
	defer closeFile(d, &err)

	s, err := os.Open(src)
	if err != nil {
		return
	}
	defer closeFile(s, &err)

	_, err = io.Copy(d, s)
	return
}

func TestSQLStore_Constructor_Create(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	require.NoError(t, err)
	require.NoError(t, s.Close())
	require.FileExists(t, dbPath)
}

func TestSQLStore_Constructor_Existing(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	require.NoError(t, err)
	require.NoError(t, s.Close())
}

func TestSQLStore_Memory_Create(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "library.base.sqlite")
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, map[string]any{OPTION_IN_MEMORY_DATABASE: true})
	require.Nil(t, err)
	require.NotNil(t, s)
	require.NoError(t, err)
	require.NoError(t, s.Close())
	require.FileExists(t, dbPath)
}

func TestSQLStore_Memory_Existing(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, map[string]any{OPTION_IN_MEMORY_DATABASE: true})
	require.Nil(t, err)
	require.NotNil(t, s)
	require.NoError(t, err)
	require.NoError(t, s.Close())
}

func TestSQLStore_Get_Int(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	objectID := int64(7)
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(objectID).Once()
	mockObject.EXPECT().SetSQLID(objectID).Return().Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	v := s.Get(mockObject, eFeature, -1)
	assert.Equal(t, 3, v)
}

func TestSQLStore_Get_Enum(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("category")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// Mystery == 0
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(7)).Once()
	mockObject.EXPECT().SetSQLID(int64(7)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	mockPackageRegitry.EXPECT().GetPackage(libraryNSURI).Return(ePackage).Once()
	v := s.Get(mockObject, eFeature, -1)
	assert.Equal(t, 0, v)

	// Biography == 2
	mockObject.EXPECT().GetSQLID().Return(int64(6)).Once()
	mockObject.EXPECT().SetSQLID(int64(6)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	v = s.Get(mockObject, eFeature, -1)
	assert.Equal(t, 2, v)
}

func TestSQLStore_Get_String_Null(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("name")
	require.NotNil(t, eFeature)

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(1)).Once()
	mockObject.EXPECT().SetSQLID(int64(1)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	// database
	dbPath := filepath.Join(t.TempDir(), "library.owner.sqlite")
	err := copyFile("testdata/library.owner.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	v := s.Get(mockObject, eFeature, -1)
	assert.Equal(t, "", v)
}

func TestSQLStore_Get_Object_Nil(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("ownerPdg")
	require.NotNil(t, eFeature)

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	v := s.Get(mockObject, eFeature, -1)
	assert.Nil(t, v)
}

func TestSQLStore_Get_Object(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	ePackage.SetEFactoryInstance(newDynamicSQLFactory())

	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eLibraryClass)

	eLibraryOwnerFeature := eLibraryClass.GetEStructuralFeatureFromName("ownerPdg")
	require.NotNil(t, eLibraryOwnerFeature)

	ePersonClass, _ := ePackage.GetEClassifier("Person").(ecore.EClass)
	require.NotNil(t, ePersonClass)

	ePersonAdressAttribute, _ := ePersonClass.GetEStructuralFeatureFromName("address").(ecore.EAttribute)
	require.NotNil(t, ePersonAdressAttribute)

	ePersonFirstNameAttribute, _ := ePersonClass.GetEStructuralFeatureFromName("firstName").(ecore.EAttribute)
	require.NotNil(t, ePersonFirstNameAttribute)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.owner.sqlite")
	err := copyFile("testdata/library.owner.sqlite", dbPath)
	require.Nil(t, err)

	mockObject := NewMockObject(t)
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject.EXPECT().GetSQLID().Return(int64(1)).Once()
	mockObject.EXPECT().SetSQLID(int64(1)).Once()
	mockObject.EXPECT().EClass().Return(eLibraryClass).Once()
	mockPackageRegitry.EXPECT().GetPackage(libraryNSURI).Return(ePackage).Once()
	v, _ := s.Get(mockObject, eLibraryOwnerFeature, -1).(ecore.EObject)
	require.NotNil(t, v)
	assert.Equal(t, ePersonClass, v.EClass())

	sqlObject, _ := v.(Object)
	require.NotNil(t, sqlObject)
	assert.Equal(t, int64(2), sqlObject.GetSQLID())

	storeObject, _ := v.(ecore.EStoreProvider)
	require.NotNil(t, storeObject)
	assert.Equal(t, s, storeObject.GetEStore())
}

func TestSQLStore_Get_Reference_WithFragment(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("author")
	require.NotNil(t, eFeature)

	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	mockObject := NewMockObject(t)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.ref.sqlite")
	err := copyFile("testdata/library.ref.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject.EXPECT().GetSQLID().Return(int64(6)).Once()
	mockObject.EXPECT().SetSQLID(int64(6)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	v, _ := s.Get(mockObject, eFeature, -1).(ecore.EObject)
	require.NotNil(t, v)
	assert.True(t, v.EIsProxy())
	assert.Equal(t, "#//@library/@writers.0", v.(ecore.EObjectInternal).EProxyURI().String())
}

func TestSQLStore_Get_Reference_WithSQLID(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eBookClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eBookClass)

	eWriterClass, _ := ePackage.GetEClassifier("Writer").(ecore.EClass)
	require.NotNil(t, eWriterClass)

	eWriterFirstNameAttribute, _ := eWriterClass.GetEStructuralFeatureFromName("firstName").(ecore.EAttribute)
	require.NotNil(t, eWriterFirstNameAttribute)

	eFeature := eBookClass.GetEStructuralFeatureFromName("author")
	require.NotNil(t, eFeature)

	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	mockObject := NewMockObject(t)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.ref.sqlite")
	err := copyFile("testdata/library.ref.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject.EXPECT().GetSQLID().Return(int64(7)).Once()
	mockObject.EXPECT().SetSQLID(int64(7)).Once()
	mockObject.EXPECT().EClass().Return(eBookClass).Once()
	mockPackageRegitry.EXPECT().GetPackage("http:///org/eclipse/emf/examples/library/library.ecore/1.0.0").Return(ePackage).Once()
	v, _ := s.Get(mockObject, eFeature, -1).(ecore.EObject)
	require.NotNil(t, v)
	assert.False(t, v.EIsProxy())
	require.Equal(t, "", v.EGet(eWriterFirstNameAttribute))
}

func TestSQLStore_Set_Int(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// set
	objectID := int64(6)
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(objectID).Once()
	mockObject.EXPECT().SetSQLID(objectID).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	oldValue := s.Set(mockObject, eFeature, -1, 5, true)
	assert.Equal(t, 4, oldValue)

	// check result by querying store directly
	copies := -1
	err = s.ExecuteQuery(context.Background(), "SELECT copies FROM Lendable WHERE lendableID=6",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			copies = stmt.ColumnInt(0)
			return nil
		}})
	require.NoError(t, err)
	assert.Equal(t, 5, copies)
}

func TestSQLStore_Set_Reference_Nil(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("author")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.ref.sqlite")
	err := copyFile("testdata/library.ref.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// set
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(6)).Once()
	mockObject.EXPECT().SetSQLID(int64(6)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	oldValue, _ := s.Set(mockObject, eFeature, -1, nil, true).(ecore.EObject)
	require.NotNil(t, oldValue)
	assert.True(t, oldValue.EIsProxy())
	assert.Equal(t, "#//@library/@writers.0", oldValue.(ecore.EObjectInternal).EProxyURI().String())

	// check result by querying store directly
	var author string
	err = s.ExecuteQuery(context.Background(), "SELECT author FROM book WHERE bookID=6",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			author = stmt.ColumnText(0)
			return nil
		}})
	require.NoError(t, err)
	require.Empty(t, author)
}

func TestSQLStore_Set_List_Primitive(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	oldValue := s.Set(mockObject, eFeature, 1, "c4", true)
	require.Equal(t, "c32", oldValue)

	// check result by querying store directly
	content := ""
	err = s.ExecuteQuery(context.Background(), "SELECT contents FROM book_contents WHERE bookID=5 AND idx=2.0",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			content = stmt.ColumnText(0)
			return nil
		}})
	require.NoError(t, err)
	require.Equal(t, "c4", content)
}

func TestSQLStore_Get_List_String(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	v := s.Get(mockObject, eFeature, 1)
	require.Equal(t, "c32", v)
}

func TestSQLStore_IsSet_SingleValue_NotSet(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("ownerPdg")
	require.NotNil(t, eFeature)

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	isSet := s.IsSet(mockObject, eFeature)
	assert.False(t, isSet)
}

func TestSQLStore_IsSet_SingleValue(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("ownerPdg")
	require.NotNil(t, eFeature)

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(1)).Once()
	mockObject.EXPECT().SetSQLID(int64(1)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	// database
	dbPath := filepath.Join(t.TempDir(), "library.owner.sqlite")
	err := copyFile("testdata/library.owner.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	isSet := s.IsSet(mockObject, eFeature)
	assert.True(t, isSet)
}

func TestSQLStore_IsSet_ManyValue_NotSet(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eFeature)

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(1)).Once()
	mockObject.EXPECT().SetSQLID(int64(1)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	// database
	dbPath := filepath.Join(t.TempDir(), "library.owner.sqlite")
	err := copyFile("testdata/library.owner.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	isSet := s.IsSet(mockObject, eFeature)
	assert.False(t, isSet)
}

func TestSQLStore_IsSet_ManyValue_Set(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eFeature)

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	// database
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	isSet := s.IsSet(mockObject, eFeature)
	assert.True(t, isSet)
}

func TestSQLStore_UnSet_Single(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// unset
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(6)).Once()
	mockObject.EXPECT().SetSQLID(int64(6)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.UnSet(mockObject, eFeature)

	// check result by querying store directly
	copies := -1
	err = s.ExecuteQuery(context.Background(), "SELECT copies FROM Lendable WHERE lendableID=6",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			copies = stmt.ColumnInt(0)
			return nil
		}})
	require.NoError(t, err)
	require.Equal(t, 0, copies)
}

func TestSQLStore_UnSet_Many(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// unset
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.UnSet(mockObject, eFeature)

	// check result by querying store directly
	count := -1
	err = s.ExecuteQuery(context.Background(), "SELECT COUNT(contents) FROM book_contents WHERE bookID=5",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			count = stmt.ColumnInt(0)
			return nil
		}})
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestSQLStore_IsEmpty_False(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.False(t, s.IsEmpty(mockObject, eFeature))
}

func TestSQLStore_IsEmpty_True(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("borrowers")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.True(t, s.IsEmpty(mockObject, eFeature))
}

func TestSQLStore_IsEmpty_NonExisting(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("borrowers")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.True(t, s.IsEmpty(mockObject, eFeature))
}

func TestSQLStore_Size(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.Equal(t, 3, s.Size(mockObject, eFeature))
}

func TestSQLStore_Size_Empty(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("borrowers")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.Equal(t, 0, s.Size(mockObject, eFeature))
}

func TestSQLStore_Size_NonExisting(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("borrowers")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.Equal(t, 0, s.Size(mockObject, eFeature))
}

func TestSQLStore_Contains_Primitive(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.True(t, s.Contains(mockObject, eFeature, "c31"))
}

func TestSQLStore_Contains_Reference(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	mockRef := NewMockObject(t)
	mockRef.EXPECT().GetSQLID().Return(int64(6)).Once()
	mockRef.EXPECT().SetSQLID(int64(6)).Once()
	assert.True(t, s.Contains(mockObject, eFeature, mockRef))
}

func TestSQLStore_Contains_NoTable(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("borrowers")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(3)).Once()
	mockObject.EXPECT().SetSQLID(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.False(t, s.Contains(mockObject, eFeature, nil))
}

func TestSQLStore_IndexOf_Existing(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	mockRef := NewMockObject(t)
	mockRef.EXPECT().GetSQLID().Return(int64(6)).Once()
	mockRef.EXPECT().SetSQLID(int64(6)).Once()
	assert.Equal(t, 0, s.IndexOf(mockObject, eFeature, mockRef))
}

func TestSQLStore_IndexOf_NonExisting(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	mockRef := NewMockObject(t)
	mockRef.EXPECT().GetSQLID().Return(int64(3)).Once()
	mockRef.EXPECT().SetSQLID(int64(3)).Once()
	assert.Equal(t, -1, s.IndexOf(mockObject, eFeature, mockRef))
}

func TestSQLStore_IndexOf_Multiple(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.indexof.sqlite")
	err := copyFile("testdata/library.indexof.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.Equal(t, 1, s.IndexOf(mockObject, eFeature, "c2"))
}

func TestSQLStore_LastIndexOf_Existing(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	mockRef := NewMockObject(t)
	mockRef.EXPECT().GetSQLID().Return(int64(7)).Once()
	mockRef.EXPECT().SetSQLID(int64(7)).Once()
	assert.Equal(t, 1, s.LastIndexOf(mockObject, eFeature, mockRef))
}

func TestSQLStore_LastIndexOf_NonExisting(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	mockRef := NewMockObject(t)
	mockRef.EXPECT().GetSQLID().Return(int64(3)).Once()
	mockRef.EXPECT().SetSQLID(int64(3)).Once()
	assert.Equal(t, -1, s.LastIndexOf(mockObject, eFeature, mockRef))
}

func TestSQLStore_LastIndexOf_Multiple(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.indexof.sqlite")
	err := copyFile("testdata/library.indexof.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.Equal(t, 2, s.LastIndexOf(mockObject, eFeature, "c2"))
}

func TestSQLStore_Remove_Object(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	ePackage.SetEFactoryInstance(newDynamicSQLFactory())

	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eLibraryClass)

	eBookClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eBookClass)

	eBooksFeature := eLibraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eBooksFeature)

	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eLibraryClass).Once()
	mockPackageRegitry.EXPECT().GetPackage(libraryNSURI).Return(ePackage).Once()
	book, _ := s.Remove(mockObject, eBooksFeature, 0, true).(Object)
	require.NotNil(t, book)
	assert.Equal(t, eBookClass, book.EClass())
	assert.Equal(t, int64(6), book.GetSQLID())
}

func TestSQLStore_Remove_NonExisting(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	ePackage.SetEFactoryInstance(newDynamicSQLFactory())

	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eLibraryClass)

	eBookClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eBookClass)

	eBooksFeature := eLibraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eBooksFeature)

	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eLibraryClass).Once()
	previous := s.Remove(mockObject, eBooksFeature, 2, true)
	assert.Nil(t, previous)
}

func TestSQLStore_Clear(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// clear
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Clear(mockObject, eFeature)

	// check result by querying store directly
	count := -1
	err = s.ExecuteQuery(context.Background(), "SELECT COUNT(contents) FROM book_contents WHERE bookID=5",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			count = stmt.ColumnInt(0)
			return nil
		}})
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestSQLStore_Add_First_Empty(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.emptylist.sqlite")
	err := copyFile("testdata/library.emptylist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// add
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Add(mockObject, eFeature, 0, "c0")

	// check result by querying store directly
	content := ""
	err = s.ExecuteQuery(context.Background(), "SELECT contents FROM book_contents WHERE bookID=5 ORDER BY idx ASC LIMIT 1",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			content = stmt.ColumnText(0)
			return nil
		}})
	require.NoError(t, err)
	require.Equal(t, "c0", content)

	// check result by querying store directly
	count := -1
	err = s.ExecuteQuery(context.Background(), "SELECT COUNT(contents) FROM book_contents WHERE bookID=5",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			count = stmt.ColumnInt(0)
			return nil
		}})
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestSQLStore_Add_First_NonEmpty(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Add(mockObject, eFeature, 0, "c0")

	// check result by querying store directly
	content := ""
	err = s.ExecuteQuery(context.Background(), "SELECT contents FROM book_contents WHERE bookID=5 ORDER BY idx ASC LIMIT 1",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			content = stmt.ColumnText(0)
			return nil
		}})
	require.NoError(t, err)
	require.Equal(t, "c0", content)
}

func TestSQLStore_Add_Last(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Add(mockObject, eFeature, 3, "c5")

	// check result by querying store directly
	content := ""
	err = s.ExecuteQuery(context.Background(), "SELECT contents FROM book_contents WHERE bookID=5 ORDER BY idx DESC LIMIT 1",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			content = stmt.ColumnText(0)
			return nil
		}})
	assert.NoError(t, err)
	assert.Equal(t, "c5", content)
}

func TestSQLStore_Add_Middle(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Add(mockObject, eFeature, 1, "c12")

	// check result by querying store directly
	content := ""
	err = s.ExecuteQuery(context.Background(), "SELECT contents FROM book_contents WHERE bookID=5 ORDER BY idx ASC LIMIT 1 OFFSET 1",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			content = stmt.ColumnText(0)
			return nil
		}})
	assert.NoError(t, err)
	assert.Equal(t, "c12", content)
}

func TestSQLStore_Add_Invalid(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	logger := newTestLogger(t.Name(), t)
	options := map[string]any{OPTION_LOGGER: logger}
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, options)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass)
	s.Add(mockObject, eFeature, 6, "c")
}

func TestSQLStore_AddAll(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	s.AddAll(mockObject, eFeature, 3, ecore.NewImmutableEList([]any{"c34", "c35"}))

	// check result by querying store directly
	count := 0
	err = s.ExecuteQuery(context.Background(), "SELECT COUNT(contents) FROM book_contents WHERE bookID=5",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			count = stmt.ColumnInt(0)
			return nil
		}})
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestSQLStore_Move_End(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// move operation
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Move(mockObject, eFeature, 0, 2, true)

	// check result by querying store directly
	idx := -1.0
	err = s.ExecuteQuery(context.Background(), "SELECT idx FROM book_contents WHERE bookID=5 ORDER BY idx DESC LIMIT 1",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			idx = stmt.ColumnFloat(0)
			return nil
		}})
	assert.NoError(t, err)
	assert.Equal(t, 4.0, idx)
}

func TestSQLStore_Move_Begin(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// move operation
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Move(mockObject, eFeature, 2, 0, true)

	// check result by querying store directly
	idx := -1.0
	err = s.ExecuteQuery(context.Background(), "SELECT idx FROM book_contents WHERE bookID=5 ORDER BY idx ASC LIMIT 1",
		&sqlitex.ExecOptions{ResultFunc: func(stmt *sqlite.Stmt) error {
			idx = stmt.ColumnFloat(0)
			return nil
		}})
	assert.NoError(t, err)
	assert.Equal(t, 0.5, idx)
}

func TestSQLStore_ToArray_Primitive(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// create store
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5)).Once()
	mockObject.EXPECT().SetSQLID(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.Equal(t, []any{"c31", "c32", "c33"}, s.ToArray(mockObject, eFeature))
}

func TestSQLStore_ToArray_Objects(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	ePackage.SetEFactoryInstance(newDynamicSQLFactory())

	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eLibraryClass)

	eBookClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eBookClass)

	eBooksFeature := eLibraryClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eBooksFeature)

	// create store
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(2)).Once()
	mockObject.EXPECT().SetSQLID(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eLibraryClass).Once()
	mockPackageRegitry.EXPECT().GetPackage(libraryNSURI).Return(ePackage).Once()
	a := s.ToArray(mockObject, eBooksFeature)
	require.Equal(t, 2, len(a))
	b, _ := a[0].(ecore.EObject)
	require.NotNil(t, b)
	require.Equal(t, eBookClass, b.EClass())
}

func TestSQLStore_GetContainer(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	ePackage.SetEFactoryInstance(newDynamicSQLFactory())

	eLibraryClass, _ := ePackage.GetEClassifier("Library").(ecore.EClass)
	require.NotNil(t, eLibraryClass)

	eBookClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eBookClass)

	// create store
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	objectID := int64(6)
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(objectID).Once()
	mockObject.EXPECT().SetSQLID(objectID).Once()
	mockPackageRegitry.EXPECT().GetPackage(libraryNSURI).Return(ePackage).Once()

	container, feature := s.GetContainer(mockObject)
	require.NotNil(t, container)
	require.Equal(t, "Library", container.EClass().GetName())
	require.NotNil(t, feature)
	require.Equal(t, "books", feature.GetName())

}

func TestSQLStore_Serialize(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	//
	bytes, err := s.Serialize(context.Background()).Await(context.Background())
	require.NoError(t, err)
	require.NotNil(t, bytes)
	requireSameDB(t, "testdata/library.store.sqlite", bytes)
}

func TestSQLStore_SerializeBig(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, map[string]any{OPTION_MAX_ALLOC_SIZE: 0})
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	//
	bytes, err := s.Serialize(context.Background()).Await(context.Background())
	require.NoError(t, err)
	require.NotNil(t, bytes)
	requireSameDB(t, "testdata/library.store.sqlite", bytes)
}

func TestSQLStore_GetRoots(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	mockPackageRegitry.EXPECT().GetPackage(libraryNSURI).Return(ePackage).Once()
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	roots := s.GetRoots()
	require.Equal(t, 1, len(roots))
}

func TestSQLStore_RemoveRoot(t *testing.T) {
	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	objectID := int64(1)
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(objectID)
	mockObject.EXPECT().SetSQLID(objectID).Return()

	s.RemoveRoot(mockObject)
	roots := s.GetRoots()
	require.Equal(t, 0, len(roots))
}

func TestSQLStore_AddRoot(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eBookClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eBookClass)

	eBookName := eBookClass.GetEStructuralFeatureFromName("title")
	require.NotNil(t, eBookName)

	// create a book
	eFactory := ePackage.GetEFactoryInstance()
	eBook := eFactory.Create(eBookClass)
	eBook.ESet(eBookName, "MyBook")

	// database
	dbPath := filepath.Join(t.TempDir(), "library.add.sqlite")

	// store
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()
	roots := s.GetRoots()
	require.Equal(t, 0, len(roots))

	s.AddRoot(eBook)

	roots = s.GetRoots()
	require.Equal(t, 1, len(roots))
}

func TestSQLStore_ParallelRead(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// object
	objectID := int64(7)
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(objectID)
	mockObject.EXPECT().SetSQLID(objectID).Return()
	mockObject.EXPECT().EClass().Return(eClass)

	// readers
	numReaders := 10
	done := &sync.WaitGroup{}
	done.Add(numReaders)
	for range numReaders {
		go func() {
			v := s.Get(mockObject, eFeature, -1)
			assert.Equal(t, 3, v)
			done.Done()
		}()
	}
	done.Wait()
}

func TestSQLStore_Parallel_GetSet(t *testing.T) {
	logger := newDebugLogger()
	require.NotNil(t, logger)

	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, nil, map[string]any{OPTION_LOGGER: logger})
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// object
	objectID := int64(7)
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(objectID)
	mockObject.EXPECT().SetSQLID(objectID).Return()
	mockObject.EXPECT().EClass().Return(eClass)

	// readers
	numOperations := 20
	done := &sync.WaitGroup{}
	done.Add(numOperations)
	for i := range numOperations {
		if i%4 == 0 {
			go func() {
				v := s.Set(mockObject, eFeature, -1, 3, true)
				assert.Equal(t, 3, v)
				done.Done()
			}()
		} else {
			go func() {
				v := s.Get(mockObject, eFeature, -1)
				assert.Equal(t, 3, v)
				done.Done()
			}()
		}
	}
	done.Wait()
}

func TestSQLStore_Parallel_AddSize(t *testing.T) {
	logger := newTestLogger(t.Name(), t)
	require.NotNil(t, logger)

	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/legacy/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, map[string]any{OPTION_LOGGER: logger})
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(5))
	mockObject.EXPECT().SetSQLID(int64(5))
	mockObject.EXPECT().EClass().Return(eClass)

	// readers
	numOperations := 20
	index := s.Size(mockObject, eFeature)
	for i := range numOperations {
		if i%4 == 0 {
			s.Add(mockObject, eFeature, index, fmt.Sprintf("c4%v", index))
			index++
		} else {
			size := s.Size(mockObject, eFeature)
			assert.NotEqual(t, 0, size)
		}
	}
}

func TestSQLStore_UnRegister(t *testing.T) {
	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	objectID := int64(1)
	mockObject := NewMockObject(t)
	mockObject2 := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(objectID)
	mockObject2.EXPECT().GetSQLID().Return(2)
	mockObject.EXPECT().EContents().Return(ecore.NewImmutableEList([]any{mockObject2})).Once()
	mockObject2.EXPECT().EContents().Return(ecore.NewImmutableEList([]any{})).Once()
	_, err = s.UnRegister(mockObject, true).Await(context.Background())
	require.NoError(t, err)
}

func TestSQLStore_Panics(t *testing.T) {
	logger := newTestLogger(t.Name(), t)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	options := map[string]any{OPTION_LOGGER: logger}
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, options)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// a panic in an operation executed by the writer is an error of the operation
	_, err = s.execute(newOperation("test", operationWrite, nil, nil, 0, nil, func() (any, error) {
		panic("invalid test")
	}, withWait()))
	require.EqualError(t, err, "invalid test")

	_, err = s.execute(newOperation("test", operationWrite, nil, nil, 0, nil, func() (any, error) {
		panic(errors.New("invalid test"))
	}, withWait()))
	require.EqualError(t, err, "invalid test")

	// a panic in a read executed on a pool connection too
	_, err = s.execute(newOperation("test", operationRead, nil, nil, 0, nil, func() (any, error) {
		panic("invalid read")
	}, withWait()))
	require.EqualError(t, err, "invalid read")

	// the store still works : the errors are reported once by the flush
	size := 0
	require.NoError(t, s.ExecuteQuery(context.Background(), "SELECT COUNT(*) from library_books", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			size = stmt.ColumnInt(0)
			return nil
		},
	}))
	require.Equal(t, 2, size)
	require.Error(t, s.WaitOperations(context.Background(), nil))
	require.NoError(t, s.WaitOperations(context.Background(), nil))
}

func TestSQLStore_UnlockOperation_Simple(t *testing.T) {
	//t.Setenv("TestLogger", "true")
	logger := newTestLogger(t.Name(), t)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	options := map[string]any{OPTION_LOGGER: logger, OPTION_OPERATION_TIMEOUT: 25 * time.Millisecond}
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, options)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	size := 0
	err = s.ExecuteQuery(context.Background(), "SELECT COUNT(*) from library_books", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			time.Sleep(50 * time.Millisecond)
			size = stmt.ColumnInt(0)
			return nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, 2, size)
}
