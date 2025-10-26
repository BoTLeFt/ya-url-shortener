package memory

import "sync"

type Memory struct {
	mu sync.RWMutex
	// id -> originalURL map
	data map[string]string
}

func New() *Memory {
	return &Memory{data: make(map[string]string)}
}

func (m *Memory) Save(id, original string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[id] = original
	return nil
}

func (m *Memory) Get(id string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[id]
	return v, ok
}
