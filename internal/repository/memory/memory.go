package memory

import (
	"sync"

	"github.com/BoTLeFt/ya-url-shortener/internal/repository"
)

type Memory struct {
	mu sync.Mutex
	// id -> originalURL map
	data map[string]string
}

func New() *Memory {
	return &Memory{data: make(map[string]string)}
}

func (m *Memory) Save(id, original string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.data[id]; exists {
		return repository.ErrIDAlreadyExists
	}
	m.data[id] = original
	return nil
}

func (m *Memory) Get(id string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[id]
	return v, ok
}
