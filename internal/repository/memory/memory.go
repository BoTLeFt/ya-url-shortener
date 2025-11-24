package memory

import (
	"sync"

	"github.com/BoTLeFt/ya-url-shortener/internal/repository"
	"github.com/BoTLeFt/ya-url-shortener/internal/repository/file"
)

type Memory struct {
	mu sync.RWMutex
	// id -> originalURL map
	data map[string]string
	file file.Producer
}

func New(file file.Producer) *Memory {
	return &Memory{data: make(map[string]string), file: file}
}

func (m *Memory) Save(id, original string, saveToFile bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.data[id]; exists {
		return repository.ErrIDAlreadyExists
	}
	m.data[id] = original
	if saveToFile {
		event := &file.ShortenedURL{
			ShortURL:    id,
			OriginalURL: original,
		}
		m.file.WriteEvent(event)
	}
	return nil
}

func (m *Memory) Get(id string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[id]
	return v, ok
}
