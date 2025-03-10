package locker

import "sync"

type Locker struct {
	mutex  sync.Mutex
	isLock bool
}

func NewLocker() *Locker {
	return &Locker{
		isLock: false,
	}
}

func (l *Locker) IsLock() bool {
	return l.isLock
}

func (l *Locker) Lock() {
	l.mutex.Lock()
	l.isLock = true
}
func (l *Locker) Unlock() {
	l.isLock = false
	l.mutex.Unlock()
}
