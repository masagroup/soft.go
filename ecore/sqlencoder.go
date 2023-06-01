package ecore

import (
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type sqlEncoderFeatureData struct {
}

type sqlEncoderObjectData struct {
	id        int64
	classData *sqlEncoderClassData
}

type sqlEncoderClassData struct {
	id    int64
	table string
}

type sqlEncoderPackageData struct {
	id        int64
	classData []*sqlEncoderClassData
}

type SQLEncoder struct {
	resource          EResource
	writer            io.Writer
	driver            string
	db                *sql.DB
	packageDataMap    map[EPackage]*sqlEncoderPackageData
	classDataMap      map[EClass]*sqlEncoderClassData
	insertContentStmt *sql.Stmt
	insertPackageStmt *sql.Stmt
	insertClassStmt   *sql.Stmt
	insertObjectStms  map[EClass]*sql.Stmt
	idAttributeName   string
}

func NewSQLEncoder(resource EResource, w io.Writer, options map[string]any) *SQLEncoder {
	e := &SQLEncoder{
		resource:         resource,
		writer:           w,
		driver:           "sqlite",
		packageDataMap:   map[EPackage]*sqlEncoderPackageData{},
		classDataMap:     map[EClass]*sqlEncoderClassData{},
		insertObjectStms: map[EClass]*sql.Stmt{},
	}
	if options != nil {
		if driver, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			e.driver = driver.(string)
		}
		e.idAttributeName, _ = options[JSON_OPTION_ID_ATTRIBUTE_NAME].(string)
	}
	return e
}

func (e *SQLEncoder) createDB() (*sql.DB, error) {
	fileName := filepath.Base(e.resource.GetURI().Path())
	dbPath, err := sqlTmpDB(fileName)
	if err != nil {
		return nil, err
	}

	// open db
	db, err := sql.Open(e.driver, dbPath)
	if err != nil {
		return nil, err
	}

	// version info
	version := fmt.Sprintf(`PRAGMA user_version = %v`, sqlCodecVersion)
	_, err = db.Exec(version)
	if err != nil {
		return nil, err
	}

	// properties infos
	properties := `
	PRAGMA synchronous = NORMAL;
	PRAGMA journal_mode = WAL;
	`
	_, err = db.Exec(properties)
	if err != nil {
		return nil, err
	}

	// common tables
	tables := `
	CREATE TABLE packages ( 
		packageID INTEGER PRIMARY KEY AUTOINCREMENT,
		uri TEXT
	);
	CREATE TABLE classes (
		classID INTEGER PRIMARY KEY AUTOINCREMENT,
		packageID INTEGER,
		name TEXT,
		FOREIGN KEY(packageID) REFERENCES packages(packageID)
	);
	CREATE TABLE contents (
		classID INTEGER,
		objectID INTEGER,
		FOREIGN KEY(classID) REFERENCES classes(classID)
	);
	`
	_, err = db.Exec(tables)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (e *SQLEncoder) EncodeResource() {
	var err error
	e.db, err = e.createDB()
	if err != nil {
		e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), e.resource.GetURI().String(), 0, 0))
		return
	}
	defer func() {
		_ = e.db.Close()
	}()

	if contents := e.resource.GetContents(); !contents.Empty() {
		object := contents.Get(0).(EObject)
		if err := e.encodeContent(object); err != nil {
			e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), e.resource.GetURI().String(), 0, 0))
			return
		}
	}
}

func (e *SQLEncoder) EncodeObject(object EObject) error {
	return nil
}

func (e *SQLEncoder) encodeContent(eObject EObject) error {
	objectData, err := e.encodeObject(eObject)
	if err != nil {
		return err
	}

	if e.insertContentStmt == nil {
		insertContentStmt, err := e.db.Prepare(`INSERT INTO contents (classID,objectID) VALUES (?,?)`)
		if err != nil {
			return err
		}
		e.insertContentStmt = insertContentStmt
	}

	if _, err := e.insertContentStmt.Exec(objectData.classData.id, objectData.id); err != nil {
		return err
	}

	return nil
}

func (e *SQLEncoder) encodeObject(eObject EObject) (*sqlEncoderObjectData, error) {
	// encode object class
	eClass := eObject.EClass()
	classData, err := e.encodeClass(eClass)
	if err != nil {
		return nil, err
	}

	idManager := e.resource.GetObjectIDManager()
	if len(e.idAttributeName) == 0 {
		idManager = nil
	}

	// create table
	insertObjectStmt, isObjectInsertStmt := e.insertObjectStms[eClass]
	if !isObjectInsertStmt {
		// create table
		{
			var query strings.Builder
			query.WriteString("CREATE TABLE ")
			query.WriteString(classData.table)
			query.WriteString(" (objectID INTEGER PRIMARY KEY AUTOINCREMENT")
			if idManager != nil {
				query.WriteString(",")
				query.WriteString(e.idAttributeName)
			}
			query.WriteString(")")

			if _, err := e.db.Exec(query.String()); err != nil {
				return nil, err
			}
		}
		// create stmt
		{
			var query strings.Builder
			query.WriteString("INSERT INTO ")
			query.WriteString(classData.table)
			query.WriteString(" (objectID")

			if idManager != nil {
				query.WriteString(",")
				query.WriteString(e.idAttributeName)
				query.WriteString(") VALUES (NULL,?)")
			} else {
				query.WriteString(") VALUES (NULL)")
			}

			insertObjectStmt, err = e.db.Prepare(query.String())
			if err != nil {
				return nil, err
			}

		}
	}

	// insert object in table
	args := []any{}
	if idManager != nil {
		args = append(args, idManager.GetID(eObject))
	}

	sqlResult, err := insertObjectStmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	// retrieve object id
	objectID, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// features
	for itFeature := eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
		eFeature := itFeature.Next().(EStructuralFeature)
		e.encodeFeatureValue(eObject, objectID, eFeature)
	}

	return &sqlEncoderObjectData{
		id:        objectID,
		classData: classData,
	}, nil
}

func (e *SQLEncoder) encodeObjectReference(eObject EObject) error {
	// encode object class
	_, err := e.encodeClass(eObject.EClass())
	if err != nil {
		return err
	}

	return nil
}

func (e *SQLEncoder) encodeFeatureValue(eObject EObject, objectID int64, eFeature EStructuralFeature) {
}

func (e *SQLEncoder) encodeClass(eClass EClass) (*sqlEncoderClassData, error) {
	eClassData := e.classDataMap[eClass]
	if eClassData == nil {
		// encode package
		ePackage := eClass.GetEPackage()
		packageData, err := e.encodePackage(ePackage)
		if err != nil {
			return nil, err
		}
		// create statement if needed
		if e.insertClassStmt == nil {
			stmt, err := e.db.Prepare(`INSERT INTO classes (packageID,name) VALUES (?,?)`)
			if err != nil {
				return nil, err
			}
			e.insertClassStmt = stmt
		}
		// insert new class
		sqlResult, err := e.insertClassStmt.Exec(packageData.id, eClass.GetName())
		if err != nil {
			return nil, err
		}
		// retrieve index
		id, err := sqlResult.LastInsertId()
		if err != nil {
			return nil, err
		}
		// create data
		eClassData = &sqlEncoderClassData{
			id:    id,
			table: ePackage.GetNsPrefix() + "_" + strings.ToLower(eClass.GetName()),
		}
		e.classDataMap[eClass] = eClassData
	}
	return eClassData, nil
}

func (e *SQLEncoder) encodePackage(ePackage EPackage) (*sqlEncoderPackageData, error) {
	ePackageData := e.packageDataMap[ePackage]
	if ePackageData == nil {
		// create statement if needed
		if e.insertPackageStmt == nil {
			stmt, err := e.db.Prepare(`INSERT INTO packages (uri) VALUES (?)`)
			if err != nil {
				return nil, err
			}
			e.insertPackageStmt = stmt
		}
		// insert new package
		sqlResult, err := e.insertPackageStmt.Exec(ePackage.GetNsURI())
		if err != nil {
			return nil, err
		}
		// retrieve package index
		id, err := sqlResult.LastInsertId()
		if err != nil {
			return nil, err
		}
		// create data
		ePackageData = &sqlEncoderPackageData{
			id:        id,
			classData: make([]*sqlEncoderClassData, ePackage.GetEClassifiers().Size()),
		}
		e.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData, nil
}
