package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/verazalayli/go_studying/redis/pkg/model"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
	Repository — это "порт" к внешнему миру хранилищ.
	Здесь мы инкапсулируем детали работы с Redis.

	Почему храним JSON?
	— Это просто и прозрачно: мы сериализуем/десериализуем нашу сущность "как есть".
	— Redis хранит строки очень эффективно, а JSON хорошо читается человеком.
	— В будущем легко заменить на RedisJSON или другой подход, не меняя сервис/хендлер.

	Ключи будем строить так: <prefix><id>, например "users:42".
*/

// ErrNotFound — когда в Redis нет записи по ключу.
var ErrNotFound = errors.New("not found")

// UserRepository описывает операции над пользователем в Redis.
type UserRepository interface {
	Save(ctx context.Context, u model.User, ttl time.Duration) error
	GetByID(ctx context.Context, id string) (model.User, error)
	Delete(ctx context.Context, id string) error
}

// userRepository — конкретная реализация через go-redis.
type userRepository struct {
	rdb        *redis.Client // клиент Redis
	keyPrefix  string        // префикс для ключей (например "users:")
	defaultTTL time.Duration // "время жизни" записи по умолчанию
}

// Опции репозитория (функциональные опции — удобный паттерн настройки).
type Option func(*userRepository)

// WithKeyPrefix позволяет задать префикс ключа (по умолчанию "users:").
func WithKeyPrefix(prefix string) Option {
	return func(r *userRepository) { r.keyPrefix = prefix }
}

// WithDefaultTTL задаёт TTL по умолчанию для всех Save, где ttl==0.
func WithDefaultTTL(ttl time.Duration) Option {
	return func(r *userRepository) { r.defaultTTL = ttl }
}

// NewUserRepository — конструктор репозитория.
func NewUserRepository(rdb *redis.Client, opts ...Option) UserRepository {
	r := &userRepository{
		rdb:        rdb,
		keyPrefix:  "users:",
		defaultTTL: 0, // 0 = без TTL по умолчанию
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// key — собирает ключ из префикса и id.
func (r *userRepository) key(id string) string {
	return fmt.Sprintf("%s%s", r.keyPrefix, id)
}

// Save — сохраняет пользователя в Redis в виде JSON.
// Если ttl==0, используем defaultTTL. Если и он 0 — запись вечная.
func (r *userRepository) Save(ctx context.Context, u model.User, ttl time.Duration) error {
	// 1) Сериализуем в JSON
	data, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}

	// 2) Определяем финальный TTL
	finalTTL := ttl
	if finalTTL == 0 {
		finalTTL = r.defaultTTL
	}

	// 3) Пишем в Redis обычной строкой (SET key value EX seconds)
	//    Если finalTTL == 0, go-redis поставит TTL=0 => "вечно".
	if err := r.rdb.Set(ctx, r.key(u.ID), data, finalTTL).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}

// GetByID — достаём пользователя по id.
// Если ключа нет — возвращаем ErrNotFound.
func (r *userRepository) GetByID(ctx context.Context, id string) (model.User, error) {
	val, err := r.rdb.Get(ctx, r.key(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, fmt.Errorf("redis get: %w", err)
	}

	var u model.User
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		return model.User{}, fmt.Errorf("unmarshal user: %w", err)
	}
	return u, nil
}

// Delete — удаляет запись по ключу. Если записи нет — считаем успехом.
func (r *userRepository) Delete(ctx context.Context, id string) error {
	if err := r.rdb.Del(ctx, r.key(id)).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}
