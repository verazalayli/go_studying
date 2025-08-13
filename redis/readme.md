# Redis Clean Architecture (Go) — Учебный проект

Этот репозиторий — минимальный и **максимально обучающий** пример сервиса на Go в стиле “чистой архитектуры”, где данные хранятся в **Redis** в виде **JSON**.
Задача проекта: показать, **как разложить код по слоям**, как **писать/читать** записи в Redis, **как проверить руками**, что в Redis всё лежит, и **как думать про бизнес‑логику**.

---

## Что мы строим

* HTTP‑сервис с 3 эндпоинтами:

    * `POST /users` — создать/обновить пользователя (с опциональным TTL)
    * `GET /users/{id}` — получить пользователя по `id`
    * `DELETE /users/{id}` — удалить пользователя по `id`
* Данные храним в Redis как **строки JSON** с ключом вида `users:<id>`.
* Архитектура по слоям:

    * **model (entity)** → доменная сущность `User`
    * **repository** → конкретная работа с Redis (знание ключей, TTL, JSON)
    * **service** → бизнес‑правила: валидация, выбор TTL, обработка ошибок
    * **handler (transport)** → HTTP: принимает/отдаёт JSON, даёт коды статусов
    * **cmd** → сборка зависимостей и запуск сервера

---

## Структура проекта

```
redis/
├─ cmd/
│  └─ app/
│     └─ main.go                  # точка входа, сборка слоёв, запуск HTTP
├─ pkg/
│  ├─ handler/
│  │  └─ http.go                  # HTTP-эндпоинты (POST/GET/DELETE)
│  ├─ model/
│  │  └─ user.go                  # доменная модель User
│  ├─ repository/
│  │  └─ user_redis.go            # Redis-логика: Set/Get/Del JSON-строк
│  └─ service/
│     └─ user_service.go          # бизнес-логика и валидация
└─ go.mod
```

---

## Требования

* Go 1.22+
* Docker (для быстрого поднятия Redis)

---

## Запуск за 2 минуты

### 1) Поднять Redis (Docker)

```bash
docker run --rm -p 6379:6379 --name redis redis:7-alpine
```

> Если нужен пароль:
> `docker run --rm -p 6379:6379 --name redis redis:7-alpine redis-server --requirepass mypass`

### 2) Запустить сервис

```bash
go mod tidy
# Если Redis с паролем - укажи окружение:
# export REDIS_ADDR=127.0.0.1:6379
# export REDIS_PASSWORD=mypass
go run ./cmd/app
```

По умолчанию сервер слушает `:8080`.
При старте он делает `PING` к Redis — если Redis недоступен, упадёт с понятной ошибкой.

---

## Переменные окружения

* `REDIS_ADDR` — адрес Redis, по умолчанию `127.0.0.1:6379`
* `REDIS_PASSWORD` — пароль (если задан в Redis), по умолчанию пусто
* (в коде выставлен `DB=0` — можно поменять в `cmd/app/main.go`)

---

## API: как отправить/достать данные

### Создать/обновить пользователя (без TTL)

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"id":"42","name":"Alice","email":"alice@example.com","age":33}'
```

**Ответ:**

```json
{"result":"saved"}
```

### Создать/обновить пользователя (с TTL = 1 час)

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"id":"100","name":"Bob","email":"bob@example.com","age":28,"ttl_seconds":3600}'
```

**Ответ:**

```json
{"result":"saved"}
```

### Получить пользователя

```bash
curl http://localhost:8080/users/42
```

**Ответ:**

```json
{
  "id": "42",
  "name": "Alice",
  "email": "alice@example.com",
  "age": 33
}
```

Если пользователя нет: `404 {"error":"user not found"}`

### Удалить пользователя

```bash
curl -X DELETE http://localhost:8080/users/42
```

**Ответ:**

```json
{"result":"deleted"}
```

---

## Как проверить, что записи реально лежат в Redis

### Вариант A: через консоль (`redis-cli`)

1. Подключиться:

```bash
# Без пароля:
docker exec -it redis redis-cli
# Или, если локально установлен redis-cli:
# redis-cli -h 127.0.0.1 -p 6379
# С паролем:
# redis-cli -h 127.0.0.1 -p 6379 -a mypass
```

2. Посмотреть ключи:

```redis
127.0.0.1:6379> KEYS users:*
1) "users:42"
2) "users:100"
```

> Примечание: в проде вместо `KEYS` используют `SCAN` (не блокирует):

```redis
127.0.0.1:6379> SCAN 0 MATCH users:* COUNT 100
```

3. Увидеть значение (JSON-строка):

```redis
127.0.0.1:6379> GET users:42
"{\"id\":\"42\",\"name\":\"Alice\",\"email\":\"alice@example.com\",\"age\":33}"
```

4. Проверить TTL:

```redis
127.0.0.1:6379> TTL users:100
(integer) 3476   # оставшиеся секунды (или -1 если без TTL)
```

5. Удалить ключ вручную:

```redis
127.0.0.1:6379> DEL users:42
(integer) 1
```

### Вариант B: через UI (RedisInsight)

RedisInsight — бесплатный графический UI от Redis (ручками смотреть ключи/TTL/значения удобно).

Шаги:

1. Скачай RedisInsight со страницы Redis (доступен под macOS/Windows/Linux).
2. Запусти и добавь подключение: Host `127.0.0.1`, Port `6379` (и пароль, если есть).
3. Открой вкладку **Browser** → увидишь ключи `users:*`.
4. Клик по ключу — увидишь JSON. Там же можно посмотреть TTL, удалить ключ и т.д.

---

## Почему JSON и как это работает

Мы сознательно выбрали самый понятный путь: **сериализуем `User` в JSON и ложим строкой** командой `SET key value [EX seconds]`.
Плюсы:

* Прозрачность: в `GET` видно тот же JSON, что и в API.
* Минимум инфраструктурного кода.
* Легко мигрировать на RedisJSON/Protobuf/MsgPack позже без изменения API/сервиса.

Минусы:

* Это просто строка: нет поиска по полям на стороне Redis (если нужно — можно добавить отдельные индексы, например, хранить `SADD users:emails <email>:<id>` и т.п.).

---

## Что делает код с точки зрения бизнес‑логики

1. **Handler (HTTP):**

    * Принимает JSON от клиента.
    * Ограничивает размер тела (1 MiB) и ставит таймауты на контекст, чтобы не зависнуть в сети.
    * Превращает HTTP‑детали в вызовы **Service**.
    * Возвращает человеку понятные **HTTP‑коды** и JSON‑ответы/ошибки.

2. **Service (бизнес‑правила):**

    * Валидирует сущность: `id`, `name`, `email` не пустые; `age >= 0`.
    * Решает, какой TTL использовать:

        * Если пришёл `ttl_seconds` в запросе — используем его.
        * Если нет — используем `defaultTTL` репозитория (в примере это 0 → “вечно”).
    * Преобразует ошибки хранилища в доменные/понятные для handler (например, `ErrNotFound` → 404).

3. **Repository (инфраструктура Redis):**

    * Знает **как строятся ключи** (`users:<id>`).
    * Делает `SET`/`GET`/`DEL`.
    * Сериализует/десериализует JSON.
    * Превращает `redis.Nil` из клиента в нашу доменную ошибку `ErrNotFound`.

Итого: **все детали Redis** изолированы в одном месте (репозиторий). Остальные слои не знают, какая именно БД используется — легко подменить.

---

## Как устроены ключи и TTL

* Ключ: `users:<id>` (например, `users:42`)
* Значение: строка JSON (`{"id":"42","name":"Alice",...}`)
* TTL:

    * Если в POST передать `"ttl_seconds": N` → запись живёт `N` секунд (Redis сам удалит).
    * Если не передать — живёт вечно (до удаления).
    * Проверить TTL можно через `TTL users:<id>`.

---

## Частые вопросы / Траблшутинг

**Q: Получаю `redis ping failed` при старте.**
A: Проверь, что контейнер Redis запущен и слушает на `127.0.0.1:6379`. Если Redis с паролем — выставь `REDIS_PASSWORD`.

**Q: `POST /users` возвращает 400 “invalid JSON”.**
A: Проверь заголовок `Content-Type: application/json` и формат JSON.

**Q: `GET /users/{id}` возвращает 404, но я уверен, что создавал.**
A: Проверь ключи командой `KEYS users:*` в `redis-cli`. Возможно, другой `id` или запись протухла по TTL.

**Q: Хочу хранить не JSON, а поля по отдельности.**
A: Сделай альтернативный репозиторий, который пишет `HSET users:<id> name ... email ... age ...` (Hash). Сервис/handler менять не придётся.

**Q: Хочу искать по email.**
A: Добавь в репозитории индексацию: например, `SET email:<email> <id>`; при `GET by email` сначала `GET email:<email>`, затем `GET users:<id>`.

---

## Краткий обзор кода (куда смотреть)

* **`pkg/model/user.go`** — структура `User`, это просто данные.
* **`pkg/repository/redisrepo/user_repo.go`** — ключи, JSON, `Set/Get/Del`, TTL.
* **`pkg/service/user/service.go`** — валидация, доменные ошибки, выбор TTL.
* **`pkg/handler/http/handler.go`** — HTTP‑маршруты, парсинг/рендеринг JSON, коды ответов.
* **`cmd/app/main.go`** — создание Redis‑клиента, сборка зависимостей, запуск сервера.

---

## Полезные команды `redis-cli` (шпаргалка)

```redis
# Показать все ключи по паттерну (осторожно в проде)
KEYS users:*

# Неблокирующий обход ключей (рекомендуется)
SCAN 0 MATCH users:* COUNT 100

# Получить значение
GET users:42

# Установить значение (как если бы писал руками)
SET users:42 "{\"id\":\"42\",\"name\":\"Alice\",\"email\":\"alice@example.com\",\"age\":33}"

# TTL в секундах (-1 = вечный ключ, -2 = ключа нет)
TTL users:42

# Удалить
DEL users:42
```

---

## Почему “чистая архитектура” здесь полезна

* **Изоляция инфраструктуры**: Redis спрятан в одном месте — легко заменить БД.
* **Тестируемость**: `service` зависит от **интерфейса** репозитория → легко мокать.
* **Читабельность**: в handler — HTTP, в service — правила, в repository — БД. Лёгкая навигация и разделение ответственности.

---

## Локальная проверка end-to-end

1. Создай пользователя:

```bash
curl -s -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"id":"u1","name":"Neo","email":"neo@matrix.io","age":29,"ttl_seconds":30}' | jq
```

2. Убедись в Redis:

```bash
docker exec -it redis redis-cli
127.0.0.1:6379> GET users:u1
"{\"id\":\"u1\",\"name\":\"Neo\",\"email\":\"neo@matrix.io\",\"age\":29}"
127.0.0.1:6379> TTL users:u1
(integer) 25
```

3. Прочитай по API:

```bash
curl -s http://localhost:8080/users/u1 | jq
```

4. Удали:

```bash
curl -s -X DELETE http://localhost:8080/users/u1 | jq
```

5. Проверь, что пропал:

```bash
curl -s http://localhost:8080/users/u1 | jq
# {"error":"user not found"}
```

---

Готово! Если хочешь, добавлю сюда `docker-compose.yml`, Makefile (линт/тест/ран), моки для `service`, или альтернативную реализацию репозитория (Hash/RedisJSON/Postgres) — скажи, что именно нужно.
