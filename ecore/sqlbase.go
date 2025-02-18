package ecore

import (
	"context"
	"strings"

	"github.com/chebyrash/promise"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type executeQueryFn func(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error

type sqlBase struct {
	codecVersion     int64
	schema           *sqlSchema
	uri              *URI
	objectIDName     string
	objectIDManager  EObjectIDManager
	isObjectID       bool
	isContainerID    bool
	sqliteManager    *taskManager
	logger           *zap.Logger
	antsPool         *ants.Pool
	promisePool      promise.Pool
	connPool         *sqlitex.Pool
	connPoolProvider func() (*sqlitex.Pool, error)
	connPoolClose    func(conn *sqlitex.Pool) error
}

func (d *sqlBase) executeSqlite(fn executeQueryFn, query string, opts *sqlitex.ExecOptions) error {
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
			conn, err := d.connPool.Take(context.Background())
			if err != nil {
				return
			}
			if err = fn(conn, query, opts); err != nil {
				logger.Error("execute query", zap.Error(err))
			} else {
				logger.Debug("execute query")
			}
			d.connPool.Put(conn)
			return
		}).Await(context.Background())
	return err
}

func (d *sqlBase) executeQuery(query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.Execute, query, opts)
}

func (d *sqlBase) executeQueryTransient(query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.ExecuteTransient, query, opts)
}

func (d *sqlBase) executeQueryScript(query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.ExecuteScript, query, opts)
}

func (d *sqlBase) decodeProperties() (map[string]string, error) {
	// result
	properties := map[string]string{}

	// check if properties table exists
	tableExists := false
	if err := d.executeQuery(`SELECT name FROM sqlite_master WHERE type='table' AND name='.properties';`, &sqlitex.ExecOptions{
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
	if err := d.executeQuery(`SELECT key,value FROM ".properties" `, &sqlitex.ExecOptions{
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
