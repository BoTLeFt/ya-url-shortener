package shortener

import (
	cryptoRand "crypto/rand"
	"errors"

	"github.com/BoTLeFt/ya-url-shortener/internal/repository"
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
	for i := 0; i < 5; i++ { // 5 итераций, если не получили уникальный id, то возвращаем ошибку
		id, err := generateID()
		if err != nil {
			return "", err
		}
		if err := s.repo.Save(id, original); err != nil {
			if errors.Is(err, repository.ErrIDAlreadyExists) {
				// Коллизия id, пробуем сгенерировать другой
				continue
			}
			return "", err
		}
		return id, nil
	}
	return "", errors.New("failed to generate unique id for URL: " + original)
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
