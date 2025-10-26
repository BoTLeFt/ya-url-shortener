package shortener

import (
	cryptoRand "crypto/rand"
)

// Параметры для генерации ID
const (
	idLen    = 7
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Repository interface {
	Save(id, original string) error
	Get(id string) (string, bool)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// Shorten генерирует уникальный id и сохраняет соответствие между id и original URL
func (s *Service) Shorten(original string) (string, error) {
	for { // иттерация пока не получим уникальный id
		id, err := generateID()
		if err != nil {
			return "", err
		}
		if _, exists := s.repo.Get(id); !exists {
			if err := s.repo.Save(id, original); err != nil {
				return "", err
			}
			return id, nil
		}
	}
}

// Resolve возвращает оригинальный URL по id
func (s *Service) Resolve(id string) (string, bool) {
	if !IsValidID(id) {
		return "", false
	}
	return s.repo.Get(id)
}

// IsValidID проверяет, соответствует ли ID формату из конфигурации
func IsValidID(id string) bool {
	if len(id) != idLen {
		return false
	}
	for _, char := range id {
		valid := false
		for _, allowed := range alphabet {
			if char == allowed {
				valid = true
				break
			}
		}
		if !valid {
			return false
		}
	}
	return true
}

func generateID() (string, error) {
	b := make([]byte, idLen)
	if _, err := cryptoRand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b), nil
}
