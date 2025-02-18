package ecore

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/chebyrash/promise"
	"github.com/petermattis/goid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/semaphore"
)

type TaskType uint8

const (
	TaskRead  = 1 << 0
	TaskWrite = 1 << 1
)

type TaskManager interface {
	ScheduleTask(objects []any, taskType TaskType, desc string, task func() (any, error)) *promise.Promise[any]

	WaitTasks(context context.Context, object any) error
}

type task struct {
	id_      int64
	type_    TaskType
	promise_ *promise.Promise[any]
}

type taskManager struct {
	mutex  sync.Mutex
	tasks  map[any][]*task
	logger *zap.Logger
	pool   promise.Pool
}

func newTaskManager(pool promise.Pool, logger *zap.Logger) *taskManager {
	return &taskManager{
		tasks:  map[any][]*task{},
		pool:   pool,
		logger: logger,
	}
}

func (s *taskManager) Close() error {
	return s.WaitTasks(context.Background(), nil)
}

func (s *taskManager) WaitTasks(context context.Context, object any) error {
	// compute tasks to wait for
	var allTasks []*task
	s.mutex.Lock()
	if object == nil {
		allTasks = make([]*task, 0)
		for _, tasks := range s.tasks {
			allTasks = append(allTasks, tasks...)
		}
	} else {
		allTasks = s.tasks[object]
	}
	s.mutex.Unlock()

	// wait for tasks to be finished
	if len(allTasks) > 0 {
		// debug
		if e := s.logger.Check(zap.DebugLevel, "waiting tasks"); e != nil {
			e.Write(zap.Int64s("ops", mapSlice(allTasks, func(index int, op *task) int64 { return op.id_ })))
		}
		// compute promises
		allPromises := mapSlice(allTasks, func(index int, op *task) *promise.Promise[any] { return op.promise_ })
		// wait for promises
		_, err := promise.AllWithPool(context, s.pool, allPromises...).Await(context)
		if err != nil {
			return err
		}
		s.logger.Debug("waiting tasks finished")
	}
	return nil
}

func (s *taskManager) ScheduleTask(objects []any, taskType TaskType, desc string, op func() (any, error)) *promise.Promise[any] {
	if objects == nil {
		return s.scheduleTaskExclusive(taskType, desc, op)
	} else {
		return s.scheduleTaskObject(objects, taskType, desc, op)
	}
}

func (s *taskManager) scheduleTaskExclusive(taskType TaskType, desc string, op func() (any, error)) *promise.Promise[any] {
	s.logger.Debug("schedule exclusive access", zap.String("desc", desc))
	return promise.NewWithPool(func(resolve func(any), reject func(error)) {
		s.mutex.Lock()
		objects := slices.Collect(maps.Keys(s.tasks))
		size := int64(len(objects))
		s.mutex.Unlock()

		run := make(chan struct{})
		locked := semaphore.NewWeighted(size)
		if err := locked.Acquire(context.Background(), size); err != nil {
			reject(err)
			return
		}

		for _, object := range objects {
			s.scheduleTaskObject([]any{object}, taskType, "exclusive", func() (any, error) {
				// the object is locked
				locked.Release(1)

				// wait for the op to be run
				<-run

				return nil, nil
			})
		}

		s.logger.Debug("waiting for exclusive access")

		// wait for all tables to be locked
		if err := locked.Acquire(context.Background(), size); err != nil {
			reject(err)
			return
		}

		// indicate all queries that operation is run
		defer close(run)

		s.logger.Debug("executing with exclusive access")
		if result, err := op(); err != nil {
			reject(err)
		} else {
			resolve(result)
		}
	}, s.pool)
}

func filterSlice[S ~[]E, E any](slice S, filter func(int, E) bool) []E {
	filteredSlice := make([]E, 0, len(slice))
	for i, v := range slice {
		if filter(i, v) {
			filteredSlice = append(filteredSlice, v)
		}
	}
	return slices.Clip(filteredSlice)
}

func mapSlice[S ~[]E, E, R any](slice S, mapper func(int, E) R) []R {
	mappedSlice := make([]R, len(slice))
	for i, v := range slice {
		mappedSlice[i] = mapper(i, v)
	}
	return mappedSlice
}

type anys []any

func (as anys) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for _, a := range as {
		arr.AppendString(fmt.Sprintf("%p", a))
	}
	return nil
}

var taskID atomic.Int64

func (s *taskManager) scheduleTaskObject(objects []any, operationType TaskType, desc string, operationFn func() (any, error)) *promise.Promise[any] {
	// create operation
	taskID := taskID.Add(1)
	op := &task{
		id_:   taskID,
		type_: operationType,
	}

	// only keep objects
	objects = filterSlice(objects, func(index int, a any) bool {
		switch a.(type) {
		case EObject, EList, EMap:
			return true
		default:
			return false
		}
	})

	s.logger.Debug("schedule", zap.Int64("goid", goid.Get()), zap.Int64("id", taskID), zap.Array("locks", anys(objects)), zap.String("desc", desc))

	s.mutex.Lock()
	// compute previous tasks
	previous := []*task{}
	for _, object := range objects {
		tasks := s.tasks[object]
		switch operationType {
		case TaskRead:
			for i := len(tasks) - 1; i >= 0; i-- {
				operation := tasks[i]
				if operation.type_ == TaskWrite {
					previous = append(previous, operation)
					break
				}
			}
		case TaskWrite:
			for i := len(tasks) - 1; i >= 0; i-- {
				operation := tasks[i]
				previous = append(previous, operation)
				if operation.type_ == TaskWrite {
					break
				}
			}
		}
	}

	op.promise_ = promise.NewWithPool(func(resolve func(any), reject func(error)) {
		logger := s.logger.With(zap.Int64("goid", goid.Get()), zap.Int64("id", taskID))
		// wait for all previous promises
		if len(previous) > 0 {
			if e := logger.Check(zap.DebugLevel, "waiting previous tasks"); e != nil {
				e.Write(zap.Int64s("previous", mapSlice(previous, func(index int, op *task) int64 { return op.id_ })))
			}
			// compute promises
			promises := mapSlice(previous, func(index int, op *task) *promise.Promise[any] { return op.promise_ })
			// wait for promises
			_, err := promise.All(context.Background(), promises...).Await(context.Background())
			if err != nil {
				logger.Error("error in previous task", zap.Error(err))
				reject(err)
				return
			}
		}
		logger.Debug("execute")
		result, err := operationFn()

		logger.Debug("cleaning")
		s.mutex.Lock()
		defer s.mutex.Unlock()
		for _, object := range objects {
			tasks := s.tasks[object]
			// retrieve operation index
			index := slices.Index(tasks, op)
			if index == -1 {
				logger.Error("unable to find task index")
				reject(errors.New("unable to find task index"))
				return
			}
			// remove operation from collection
			copy(tasks[index:], tasks[index+1:])
			tasks[len(tasks)-1] = nil
			tasks = tasks[:len(tasks)-1]
			if len(tasks) == 0 {
				// no more tasks - remove object from map
				delete(s.tasks, object)
			} else {
				// set remaining tasks
				s.tasks[object] = tasks
			}
		}
		logger.Debug("cleaned")
		if len(s.tasks) == 0 {
			logger.Debug("no pending")
		}

		// operation result
		if err != nil {
			reject(err)
		} else {
			resolve(result)
		}
	}, s.pool)
	// add operation
	for _, object := range objects {
		s.tasks[object] = append(s.tasks[object], op)
	}
	s.mutex.Unlock()
	return op.promise_
}
