package ecore

import (
	"database/sql"
	"fmt"
)

type SQLStore struct {
	schema *sqlSchema
	db     *sql.DB
}

func NewSQLStore(driver string, dbPath string, options map[string]any) (*SQLStore, error) {
	// Open db
	db, err := sql.Open(driver, dbPath)
	if err != nil {
		return nil, err
	}

	// properties
	propertiesQuery := `
	PRAGMA synchronous = NORMAL;
	PRAGMA journal_mode = WAL;
	`
	_, err = db.Exec(propertiesQuery)
	if err != nil {
		return nil, err
	}

	// version
	if row := db.QueryRow("PRAGMA user_version;"); row == nil {
		// create version
		versionQuery := fmt.Sprintf(`PRAGMA user_version = %v`, sqlCodecVersion)
		_, err = db.Exec(versionQuery)
		if err != nil {
			return nil, err
		}
	} else {
		// retrieve version
		var v int
		if err := row.Scan(&v); err != nil {
			return nil, err
		}
		if v != sqlCodecVersion {
			return nil, fmt.Errorf("history version %v is not supported", v)
		}
	}

	// schema options
	schemaOptions := []sqlSchemaOption{}
	idAttributeName, _ := options[SQL_OPTION_ID_ATTRIBUTE_NAME].(string)
	if len(idAttributeName) > 0 {
		schemaOptions = append(schemaOptions, withIDAttributeName(idAttributeName))
	}

	return &SQLStore{
		schema: newSqlSchema(schemaOptions...),
		db:     db,
	}, nil
}
