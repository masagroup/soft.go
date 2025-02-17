package ecore

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type executeQueryFn func(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error

type sqlBase struct {
	codecVersion    int64
	schema          *sqlSchema
	uri             *URI
	objectIDName    string
	objectIDManager EObjectIDManager
	isObjectID      bool
	isContainerID   bool
	sqliteManager   *taskManager
	logger          *zap.Logger
}

func (d *sqlBase) executeSqlite(fn executeQueryFn, conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error {
	// retrieve task type
	var taskType TaskType
	switch operation, _, _ := strings.Cut(query, " "); operation {
	case "SELECT":
		taskType = TaskRead
	default:
		taskType = TaskWrite
	}

	// schedule sqlite task :
	// only one write to db is active
	// multiple read to db is allowed
	_, err := d.sqliteManager.ScheduleTask([]any{d}, taskType, query,
		func() (res any, err error) {
			args := []zap.Field{zap.String("query", query)}
			if opts != nil {
				args = append(args, zap.Any("args", opts.Args))
			}

			logger := d.logger.With(args...)
			err = fn(conn, query, opts)
			if err != nil {
				logger.Error("execute query", zap.Error(err))
			} else {
				logger.Debug("execute query")
			}
			return
		}).Await(context.Background())
	return err
}

func (d *sqlBase) executeQuery(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.Execute, conn, query, opts)
}

func (d *sqlBase) executeQueryTransient(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.ExecuteTransient, conn, query, opts)
}

func (d *sqlBase) executeQueryScript(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.ExecuteScript, conn, query, opts)
}

func (d *sqlBase) decodeProperties(conn *sqlite.Conn) (map[string]string, error) {
	// result
	properties := map[string]string{}

	// check if properties table exists
	tableExists := false
	if err := d.executeQuery(conn, `SELECT name FROM sqlite_master WHERE type='table' AND name='.properties';`, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			tableExists = true
			return nil
		},
	}); err != nil {
		return nil, err
	}
	if !tableExists {
		return properties, nil
	}

	// retrieve properties from table
	if err := d.executeQuery(conn, `SELECT key,value FROM ".properties" `, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			key := stmt.ColumnText(0)
			value := stmt.ColumnText(1)
			properties[key] = value
			return nil
		},
	}); err != nil {
		return nil, err
	}
	return properties, nil
}
