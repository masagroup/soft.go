package ecore

type ELockProvider interface {
	Lock()
	Unlock()
}
