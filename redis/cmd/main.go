package main

import (
	"context"
	"github.com/verazalayli/go_studying/redis/pkg/handler"
	"github.com/verazalayli/go_studying/redis/pkg/repository"
	"github.com/verazalayli/go_studying/redis/pkg/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
	Это точка входа приложения.

	Здесь мы:
	1) Подключаемся к Redis.
	2) Собираем зависимости слоями в стиле "чистой архитектуры":
	   handler -> service -> repository -> redis.Client
	3) Поднимаем HTTP-сервер с простыми REST-эндпоинтами.
*/

// createRedisClient создаёт клиент Redis и проверяет соединение.
func createRedisClient() *redis.Client {
	// Обычно эти значения берём из переменных окружения .env
	addr := env("REDIS_ADDR", "127.0.0.1:6379")
	password := env("REDIS_PASSWORD", "") // если пароль не задан — пустая строка
	db := 0                               // номер БД в Redis (по умолчанию 0)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Проверяем, что Redis доступен.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
	log.Printf("connected to redis at %s", addr)
	return client
}

// env — утилита для чтения переменных окружения с значением по умолчанию.
func env(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func main() {
	// 1) Подключаемся к Redis
	rdb := createRedisClient()
	defer func() { _ = rdb.Close() }()

	// 2) Сборка зависимостей снизу вверх:
	// repository -> service -> handler
	userRepo := repository.NewUserRepository(rdb,
		repository.WithKeyPrefix("users:"), // ключи будут вида users:<id>
		repository.WithDefaultTTL(0),       // TTL=0 означает "без срока" (можно поменять)
	)

	userService := service.NewService(userRepo) // сервис зависит от интерфейса репозитория
	h := handler.New(userService)               // handler зависит от интерфейса сервиса

	// 3) HTTP сервер
	server := &http.Server{
		Addr:         ":8080",
		Handler:      h.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// 4) Запускаем сервер и аккуратно завершаем по Ctrl+C
	go func() {
		log.Printf("HTTP server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	// Ожидаем сигнал завершения (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// Плавное завершение сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}
