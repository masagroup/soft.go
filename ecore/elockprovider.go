package ecore

// ELockProvider is an object whose writes can be prevented while it is encoded. The lock is not
// reentrant: the goroutine holding it can read the object, not write it.
type ELockProvider interface {
	Lock()
	Unlock()
}

// eLockedStoreProvider is a store object whose store can be set by the goroutine holding its
// lock (ELockProvider).
type eLockedStoreProvider interface {
	setEStoreLocked(store EStore)
}
