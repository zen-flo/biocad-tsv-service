package queue

import "sync"

// Manager manages a set of queued items safely
type Manager struct {
	mu    sync.Mutex
	files map[string]struct{}
}

// New creates a new Queue Manager
func New() *Manager {
	return &Manager{
		files: make(map[string]struct{}),
	}
}

// Add adds a file to the queue if it's not already present
func (qm *Manager) Add(file string) bool {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	if _, ok := qm.files[file]; ok {
		return false
	}
	qm.files[file] = struct{}{}
	return true
}

// Remove removes a file from the queue
func (qm *Manager) Remove(file string) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	delete(qm.files, file)
}

// Exists checks if a file is already queued
func (qm *Manager) Exists(file string) bool {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	_, ok := qm.files[file]
	return ok
}
