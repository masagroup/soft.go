package ecore

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/chebyrash/promise"
	"github.com/panjf2000/ants/v2"
	"github.com/petermattis/goid"
	"go.uber.org/zap"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type executeQueryFn func(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error

var queryID atomic.Int64

type queryType uint8

const (
	queryRead  queryType = 1 << 0
	queryWrite queryType = 1 << 1
)

type query struct {
	id      int64
	type_   queryType
	cmd     string
	promise *promise.Promise[any]
}

func newQuery(cmd string) *query {
	var queryType queryType
	switch operation, _, _ := strings.Cut(cmd, " "); operation {
	case "SELECT":
		queryType = queryRead
	default:
		queryType = queryWrite
	}
	return &query{
		id:    queryID.Add(1),
		type_: queryType,
		cmd:   cmd,
	}
}

type sqlBase struct {
	codecVersion     int64
	schema           *sqlSchema
	uri              *URI
	objectIDName     string
	objectIDManager  EObjectIDManager
	isObjectID       bool
	isContainerID    bool
	sqliteMutex      sync.Mutex
	sqliteQueries    []*query
	logger           *zap.Logger
	antsPool         *ants.Pool
	promisePool      promise.Pool
	connPool         *sqlitex.Pool
	connPoolProvider func() (*sqlitex.Pool, error)
	connPoolClose    func(conn *sqlitex.Pool) error
}

// execute sqlite cmd
func (s *sqlBase) executeSqlite(fn executeQueryFn, cmd string, opts *sqlitex.ExecOptions) error {
	s.sqliteMutex.Lock()

	// create query
	q := newQuery(cmd)

	// log
	args := []zap.Field{zap.Int64("goid", goid.Get()), zap.Int64("id", q.id), zap.String("query", cmd)}
	if opts != nil {
		args = append(args, zap.Any("args", opts.Args))
	}
	s.logger.Named("sqlite").Debug("schedule", args...)

	// compute previous query
	// only one write to db is active
	// multiple read to db is allowed
	var previous *query
	switch q.type_ {
	case queryRead:
		for i := len(s.sqliteQueries) - 1; i >= 0; i-- {
			if operation := s.sqliteQueries[i]; operation.type_ == queryWrite {
				previous = operation
				break
			}
		}
	case queryWrite:
		if len(s.sqliteQueries) > 0 {
			previous = s.sqliteQueries[len(s.sqliteQueries)-1]
		}
	}

	// add query to queries
	s.sqliteQueries = append(s.sqliteQueries, q)

	// create query promise
	q.promise = promise.NewWithPool(func(resolve func(any), reject func(error)) {
		logger := s.logger.Named("sqlite").With(zap.Int64("goid", goid.Get()), zap.Int64("id", q.id))
		// wait for previous query to be finished
		if previous != nil {
			if e := logger.Check(zap.DebugLevel, "waiting previous query"); e != nil {
				e.Write(zap.Int64("previous", previous.id))
			}
			_, err := previous.promise.Await(context.Background())
			if err != nil {
				logger.Error("error in previous", zap.Error(err))
				reject(err)
				return
			}
		}

		// execute query
		conn, err := s.connPool.Take(context.Background())
		if err != nil {
			reject(err)
			return
		}
		defer s.connPool.Put(conn)

		if err := fn(conn, cmd, opts); err != nil {
			logger.Error("execute", zap.Error(err))
			reject(err)
			return
		} else {
			logger.Error("execute")
		}

		// clean query
		logger.Debug("cleaning")
		s.sqliteMutex.Lock()
		defer s.sqliteMutex.Unlock()

		// retrieve operation index
		index := slices.Index(s.sqliteQueries, q)
		if index == -1 {
			reject(errors.New("unable to find query index"))
			return
		}
		// remove query from queries
		copy(s.sqliteQueries[index:], s.sqliteQueries[index+1:])
		s.sqliteQueries[len(s.sqliteQueries)-1] = nil
		s.sqliteQueries = s.sqliteQueries[:len(s.sqliteQueries)-1]
		logger.Debug("cleaned")
		if len(s.sqliteQueries) == 0 {
			logger.Debug("no pending")
		}
		resolve(nil)
	}, s.promisePool)

	s.sqliteMutex.Unlock()

	// wait for query to be finished
	_, err := q.promise.Await(context.Background())
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
