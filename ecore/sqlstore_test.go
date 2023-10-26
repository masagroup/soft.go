package ecore

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dynamicSQLEObjectImpl struct {
	DynamicEObjectImpl
	sqlID int64
}

func newDynamicSQLEObjectImpl() *dynamicSQLEObjectImpl {
	o := new(dynamicSQLEObjectImpl)
	o.SetInterfaces(o)
	o.Initialize()
	return o
}

func (o *dynamicSQLEObjectImpl) SetSqlID(sqlID int64) {
	o.sqlID = sqlID
}

func (o *dynamicSQLEObjectImpl) GetSqlID() int64 {
	return o.sqlID
}

type dynamicSQLFactory struct {
	EFactoryExt
}

func newDynamicSQLFactory() *dynamicSQLFactory {
	eFactory := new(dynamicSQLFactory)
	eFactory.SetInterfaces(eFactory)
	eFactory.Initialize()
	return eFactory
}

func (eFactory *dynamicSQLFactory) Create(eClass EClass) EObject {
	if eFactory.GetEPackage() != eClass.GetEPackage() || eClass.IsAbstract() {
		panic("The class '" + eClass.GetName() + "' is not a valid classifier")
	}
	if IsMapEntry(eClass) {
		eEntry := NewDynamicEMapEntryImpl()
		eEntry.SetEClass(eClass)
		return eEntry
	} else {
		eObject := newDynamicSQLEObjectImpl()
		eObject.SetEClass(eClass)
		return eObject
	}
}

func TestSQLStore_Constructor(t *testing.T) {
	s, err := NewSQLStore("testdata/library.store.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	require.Nil(t, s.Close())
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

func TestSQLStore_Get_Int(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// store
	s, err := NewSQLStore("testdata/library.store.sqlite", NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	v := s.Get(mockObject, eFeature, -1)
	assert.Equal(t, 4, v)
}

func TestSQLStore_Get_Enum(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("category")
	require.NotNil(t, eFeature)

	// store
	mockPackageRegitry := NewMockEPackageRegistry(t)
	s, err := NewSQLStore("testdata/library.store.sqlite", NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// Mystery == 0
	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(4)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	mockPackageRegitry.EXPECT().GetPackage("http:///org/eclipse/emf/examples/library/library.ecore/1.0.0").Return(ePackage).Once()
	v := s.Get(mockObject, eFeature, -1)
	assert.Equal(t, 0, v)

	// Biography == 2
	mockObject.EXPECT().GetSqlID().Return(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	v = s.Get(mockObject, eFeature, -1)
	assert.Equal(t, 2, v)
}

func TestSQLStore_Get_String_Null(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("name")
	require.NotNil(t, eFeature)

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(1)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	s, err := NewSQLStore("testdata/library.owner.sqlite", NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	v := s.Get(mockObject, eFeature, -1)
	assert.Equal(t, "", v)
}

func TestSQLStore_Get_Object_Nil(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("ownerPdg")
	require.NotNil(t, eFeature)

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	s, err := NewSQLStore("testdata/library.store.sqlite", NewURI(""), nil, nil, nil)
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

	eLibraryClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eLibraryClass)

	eLibraryOwnerFeature := eLibraryClass.GetEStructuralFeatureFromName("ownerPdg")
	require.NotNil(t, eLibraryOwnerFeature)

	ePersonClass, _ := ePackage.GetEClassifier("Person").(EClass)
	require.NotNil(t, ePersonClass)

	ePersonAdressAttribute, _ := ePersonClass.GetEStructuralFeatureFromName("address").(EAttribute)
	require.NotNil(t, ePersonAdressAttribute)

	ePersonFirstNameAttribute, _ := ePersonClass.GetEStructuralFeatureFromName("firstName").(EAttribute)
	require.NotNil(t, ePersonFirstNameAttribute)

	mockObject := NewMockSQLObject(t)
	mockPackageRegitry := NewMockEPackageRegistry(t)
	s, err := NewSQLStore("testdata/library.owner.sqlite", NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject.EXPECT().GetSqlID().Return(int64(1)).Once()
	mockObject.EXPECT().EClass().Return(eLibraryClass).Once()
	mockPackageRegitry.EXPECT().GetPackage("http:///org/eclipse/emf/examples/library/library.ecore/1.0.0").Return(ePackage).Once()
	v, _ := s.Get(mockObject, eLibraryOwnerFeature, -1).(SQLObject)
	require.NotNil(t, v)
	assert.Equal(t, ePersonClass, v.EClass())
	assert.Equal(t, int64(2), v.GetSqlID())
}

func TestSQLStore_Get_Reference(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("author")
	require.NotNil(t, eFeature)

	mockPackageRegitry := NewMockEPackageRegistry(t)
	mockObject := NewMockSQLObject(t)

	s, err := NewSQLStore("testdata/library.complex.sqlite", NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject.EXPECT().GetSqlID().Return(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	v, _ := s.Get(mockObject, eFeature, -1).(EObject)
	require.NotNil(t, v)
	assert.True(t, v.EIsProxy())
	assert.Equal(t, "#//@library/@writers.0", v.(EObjectInternal).EProxyURI().String())
}

func TestSQLStore_Set_Int(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	oldValue := s.Set(mockObject, eFeature, -1, 5)
	assert.Equal(t, 4, oldValue)

	// check store
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()
	var copies int
	row := db.QueryRow("SELECT copies FROM Lendable WHERE lendableID=3")
	err = row.Scan(&copies)
	assert.NoError(t, err)
	assert.Equal(t, 5, copies)
}

func TestSQLStore_Set_Reference_Nil(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("author")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.complex.sqlite")
	err := copyFile("testdata/library.complex.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// set
	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	oldValue, _ := s.Set(mockObject, eFeature, -1, nil).(EObject)
	require.NotNil(t, oldValue)
	assert.True(t, oldValue.EIsProxy())
	assert.Equal(t, "#//@library/@writers.0", oldValue.(EObjectInternal).EProxyURI().String())

	// check store
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()
	var author sql.NullString
	row := db.QueryRow("SELECT author FROM book WHERE bookID=3")
	err = row.Scan(&author)
	assert.NoError(t, err)
	assert.False(t, author.Valid)
}

func TestSQLStore_Set_List_Primitive(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/library.datalist.sqlite", dbPath)
	require.Nil(t, err)

	// create store
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	oldValue := s.Set(mockObject, eFeature, 1, "c4")
	assert.Equal(t, "c2", oldValue)

	// check store
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()
	var contents string
	row := db.QueryRow("SELECT contents FROM book_contents WHERE bookID=5 AND idx=1.0")
	err = row.Scan(&contents)
	assert.NoError(t, err)
	assert.Equal(t, "c4", contents)
}

func TestSQLStore_Get_List_String(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// create store
	s, err := NewSQLStore("testdata/library.datalist.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.Get(mockObject, eFeature, 1)
}

func TestSQLStore_IsSet_SingleValue_NotSet(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("ownerPdg")
	require.NotNil(t, eFeature)

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	s, err := NewSQLStore("testdata/library.store.sqlite", NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	isSet := s.IsSet(mockObject, eFeature)
	assert.False(t, isSet)
}

func TestSQLStore_IsSet_SingleValue(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("ownerPdg")
	require.NotNil(t, eFeature)

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(1)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	s, err := NewSQLStore("testdata/library.owner.sqlite", NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	isSet := s.IsSet(mockObject, eFeature)
	assert.True(t, isSet)
}

func TestSQLStore_IsSet_ManyValue_NotSet(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eFeature)

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(1)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	s, err := NewSQLStore("testdata/library.owner.sqlite", NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	isSet := s.IsSet(mockObject, eFeature)
	assert.False(t, isSet)
}

func TestSQLStore_IsSet_ManyValue_Set(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("books")
	require.NotNil(t, eFeature)

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(2)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()

	s, err := NewSQLStore("testdata/library.complex.sqlite", NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	isSet := s.IsSet(mockObject, eFeature)
	assert.True(t, isSet)
}

func TestSQLStore_UnSet_Single(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Lendable").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("copies")
	require.NotNil(t, eFeature)

	// database
	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	err := copyFile("testdata/library.store.sqlite", dbPath)
	require.Nil(t, err)

	// store
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(3)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.UnSet(mockObject, eFeature)

	// check store
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()
	var copies int
	row := db.QueryRow("SELECT copies FROM Lendable WHERE lendableID=3")
	err = row.Scan(&copies)
	assert.NoError(t, err)
	assert.Equal(t, 0, copies)
}

func TestSQLStore_UnSet_Many(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	dbPath := filepath.Join(t.TempDir(), "library.datalist.sqlite")
	err := copyFile("testdata/library.datalist.sqlite", dbPath)
	require.Nil(t, err)
	// create store
	s, err := NewSQLStore(dbPath, NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	s.UnSet(mockObject, eFeature)

	// check store
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()
	var count int
	row := db.QueryRow("SELECT COUNT(contents) FROM book_contents WHERE bookID=5")
	err = row.Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestSQLStore_IsEmpty_False(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// create store
	s, err := NewSQLStore("testdata/library.datalist.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.False(t, s.IsEmpty(mockObject, eFeature))
}

func TestSQLStore_IsEmpty_True(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("borrowers")
	require.NotNil(t, eFeature)

	// create store
	s, err := NewSQLStore("testdata/library.complex.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.True(t, s.IsEmpty(mockObject, eFeature))
}

func TestSQLStore_IsEmpty_NonExisting(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("borrowers")
	require.NotNil(t, eFeature)

	// create store
	s, err := NewSQLStore("testdata/library.datalist.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.True(t, s.IsEmpty(mockObject, eFeature))
}

func TestSQLStore_Size(t *testing.T) {
	ePackage := loadPackage("library.datalist.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Book").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("contents")
	require.NotNil(t, eFeature)

	// create store
	s, err := NewSQLStore("testdata/library.datalist.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.Equal(t, 3, s.Size(mockObject, eFeature))
}

func TestSQLStore_Size_Empty(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("borrowers")
	require.NotNil(t, eFeature)

	// create store
	s, err := NewSQLStore("testdata/library.complex.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.Equal(t, 0, s.Size(mockObject, eFeature))
}

func TestSQLStore_Size_NonExisting(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)

	eClass, _ := ePackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eClass)

	eFeature := eClass.GetEStructuralFeatureFromName("borrowers")
	require.NotNil(t, eFeature)

	// create store
	s, err := NewSQLStore("testdata/library.datalist.sqlite", NewURI(""), nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer s.Close()

	mockObject := NewMockSQLObject(t)
	mockObject.EXPECT().GetSqlID().Return(int64(5)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	assert.Equal(t, 0, s.Size(mockObject, eFeature))
}
