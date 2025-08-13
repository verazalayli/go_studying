package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/verazalayli/go_studying/redis/pkg/model"
	"github.com/verazalayli/go_studying/redis/pkg/repository"
	"github.com/verazalayli/go_studying/redis/pkg/service"
	"io"
	"net/http"
	"strings"
	"time"
)

/*
	HTTP handler — слой входного транспорта.
	Он ничего не знает про Redis: только про сервисный интерфейс.

	Маршруты:
	POST   /users        — создать/обновить пользователя (тело JSON)
	GET    /users/{id}   — получить пользователя
	DELETE /users/{id}   — удалить пользователя
	GET    /health       — простая проверка живости
*/

// Handler хранит зависимости для HTTP.
type Handler struct {
	svc service.Service
}

// New — конструктор Handler.
func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// Routes — регистрирует маршруты в стандартном http.ServeMux.
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", h.createOrUpdateUser)
	mux.HandleFunc("GET /users/", h.getUserByID) // ожидаем /users/{id}
	mux.HandleFunc("DELETE /users/", h.deleteUserByID)
	mux.HandleFunc("GET /health", h.health)
	return mux
}

// health — простой healthcheck, полезно для readiness/liveness probe.
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// createOrUpdateUser — читает JSON из тела запроса и сохраняет в Redis через сервис.
// Пример тела:
//
//	{
//	  "id": "42",
//	  "name": "Alice",
//	  "email": "alice@example.com",
//	  "age": 33,
//	  "ttl_seconds": 3600  // необязательно: срок жизни записи в секундах
//	}
func (h *Handler) createOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Ограничиваем размер тела, чтобы защититься от слишком больших запросов.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "cannot read body: "+err.Error())
		return
	}

	// Вводной DTO: в него будем парсить JSON из запроса.
	var in struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Age        int    `json:"age"`
		TTLSeconds *int   `json:"ttl_seconds"` // необязательное поле
	}
	if err := json.Unmarshal(body, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Составляем доменную модель.
	u := model.User{
		ID:    strings.TrimSpace(in.ID),
		Name:  strings.TrimSpace(in.Name),
		Email: strings.TrimSpace(in.Email),
		Age:   in.Age,
	}

	// Опциональный TTL.
	var ttlPtr *time.Duration
	if in.TTLSeconds != nil {
		t := time.Duration(*in.TTLSeconds) * time.Second
		ttlPtr = &t
	}

	// Контекст с таймаутом, чтобы не зависнуть в сетевых операциях.
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.svc.CreateOrUpdateUser(ctx, u, ttlPtr); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"result": "saved"})
}

// getUserByID — достает пользователя по id из URL, например: GET /users/42
func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	id = strings.TrimSpace(id)
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required in path, e.g. /users/42")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	u, err := h.svc.GetUser(ctx, id)
	if err != nil {
		// Превращаем доменные ошибки в статусы HTTP.
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, u)
}

// deleteUserByID — удаляет пользователя по id: DELETE /users/42
func (h *Handler) deleteUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	id = strings.TrimSpace(id)
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required in path, e.g. /users/42")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.svc.DeleteUser(ctx, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
}

// --- Утилиты ответа JSON ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
