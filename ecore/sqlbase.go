package ecore

import "database/sql"

type sqlBase struct {
	db              *sql.DB
	schema          *sqlSchema
	uri             *URI
	idAttributeName string
	idManager       EObjectIDManager
}

func (s *sqlBase) encodeProperties() error {
	// properties
	propertiesQuery := `
	PRAGMA synchronous = NORMAL;
	PRAGMA journal_mode = WAL;
	`
	_, err := s.db.Exec(propertiesQuery)
	return err
}

func (s *sqlBase) encodeSchema() error {
	// tables
	for _, table := range []*sqlTable{
		s.schema.packagesTable,
		s.schema.classesTable,
		s.schema.objectsTable,
		s.schema.contentsTable,
		s.schema.enumsTable,
	} {
		if _, err := s.db.Exec(table.createQuery()); err != nil {
			return err
		}
	}
	return nil
}
