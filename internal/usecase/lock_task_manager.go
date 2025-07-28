package usecase

import "sync"

type LockTaskManager struct {
	mu    sync.Mutex
	locks map[uint64]*sync.Mutex
}

func NewLockTaskManager() *LockTaskManager {
	return &LockTaskManager{
		locks: make(map[uint64]*sync.Mutex),
	}
}

func (m *LockTaskManager) GetLock(id uint64) *sync.Mutex {
	m.mu.Lock()
	defer m.mu.Unlock()

	if l, ok := m.locks[id]; ok {
		return l
	}

	lock := &sync.Mutex{}
	m.locks[id] = lock
	return lock
}
