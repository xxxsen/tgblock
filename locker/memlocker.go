package locker

import "sync"

type MemLocker struct {
	lockMap map[string]bool
	lck     sync.Mutex
}

func NewMemLocker() Locker {
	return &MemLocker{
		lockMap: make(map[string]bool),
	}
}

func (l *MemLocker) Lock(key string) bool {
	l.lck.Lock()
	defer l.lck.Unlock()
	if _, ok := l.lockMap[key]; ok {
		return false
	}
	l.lockMap[key] = true
	return true
}

func (l *MemLocker) Unlock(key string) {
	l.lck.Lock()
	defer l.lck.Unlock()
	delete(l.lockMap, key)
}
