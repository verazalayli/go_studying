package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/verazalayli/go_studying/redis/pkg/model"
	"github.com/verazalayli/go_studying/redis/pkg/repository"
	"strings"
	"time"
)

/*
	Сервис — бизнес-логика. Здесь мы:
	- валидируем входные данные,
	- решаем, какой TTL использовать,
	- оборачиваем ошибки доменными.

	Сервис зависит от интерфейса репозитория, а не от Redis-клиента.
	Это позволяет легко подменить реализацию (например, на Postgres)
	или написать юнит‑тесты с моками.
*/

// Repository — минимальный контракт, который нужен сервису.
type Repository interface {
	Save(ctx context.Context, u model.User, ttl time.Duration) error
	GetByID(ctx context.Context, id string) (model.User, error)
	Delete(ctx context.Context, id string) error
}

// Service — публичный интерфейс сервиса.
type Service interface {
	CreateOrUpdateUser(ctx context.Context, u model.User, ttl *time.Duration) error
	GetUser(ctx context.Context, id string) (model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

// NewService — конструктор сервиса.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// validate — простейшая валидация сущности.
func (s *service) validate(u model.User) error {
	if strings.TrimSpace(u.ID) == "" {
		return errors.New("id is required")
	}
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}
	if u.Age < 0 {
		return errors.New("age must be >= 0")
	}
	return nil
}

// CreateOrUpdateUser — создаёт или обновляет пользователя.
// TTL можно передать (ttl!=nil), либо оставить nil — тогда используем TTL по умолчанию репозитория.
func (s *service) CreateOrUpdateUser(ctx context.Context, u model.User, ttl *time.Duration) error {
	if err := s.validate(u); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	var t time.Duration
	if ttl != nil {
		t = *ttl
	}
	if err := s.repo.Save(ctx, u, t); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

// GetUser — получить пользователя по id.
func (s *service) GetUser(ctx context.Context, id string) (model.User, error) {
	if strings.TrimSpace(id) == "" {
		return model.User{}, errors.New("id is required")
	}
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// Пробрасываем как есть — хендлер превратит в 404.
			return model.User{}, err
		}
		return model.User{}, fmt.Errorf("get by id: %w", err)
	}
	return u, nil
}

// DeleteUser — удалить пользователя по id.
func (s *service) DeleteUser(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id is required")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
