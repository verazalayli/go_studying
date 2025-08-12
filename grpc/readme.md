# Clean Architecture + gRPC (Go) — README

Этот учебный проект показывает, как собрать минимальный gRPC‑сервис в стиле «чистой архитектуры» с двумя сервисами: **server** и **client**.
Внутри — разделение на слои: handler (gRPC), service (use cases), repository (in‑memory).

---

Да, давай разложу по шагам **что именно мы делаем в этом учебном проекте** и какую логику реализует каждая RPC.

---

## Общая идея

У нас есть сервис управления заметками (**NoteService**), работающий по gRPC.
Он умеет:

1. Принимать от клиента запрос на создание новой заметки.
2. Хранить созданные заметки (в примере — в памяти).
3. Возвращать по запросу одну заметку по её ID.
4. Возвращать список всех заметок.

То есть клиент и сервер общаются в формате **запрос → ответ**, и весь обмен идёт в protobuf по HTTP/2.

---

## Как это выглядит по шагам

### 1. **Создание заметки** (`CreateNote`)

**Что делает клиент**:

* Формирует сообщение `CreateNoteRequest` с полями:

  * `title` — заголовок.
  * `content` — текст заметки.
* Отправляет это сообщение на сервер через RPC `CreateNote`.

**Что делает сервер**:

* Получает `CreateNoteRequest`.
* В хендлере (на стороне сервера) вызывается метод `svc.Create(...)`, который:

  * Генерирует уникальный ID (UUID).
  * Ставит текущее время (`CreatedAt`).
  * Сохраняет заметку в хранилище (в примере — просто в `map` в памяти).
* Формирует ответ `CreateNoteResponse` с полной заметкой (`Note`), включая:

  * `id` (сгенерированный сервером).
  * `title` (из запроса).
  * `content` (из запроса).
  * `created_at` (в секундах Unix).
* Отправляет этот ответ обратно клиенту.

**Что получает клиент**:

* Готовую заметку с ID и временем создания.
* Может использовать этот ID, чтобы потом запросить заметку.

---

### 2. **Получение заметки по ID** (`GetNote`)

**Что делает клиент**:

* Отправляет `GetNoteRequest` с полем:

  * `id` — идентификатор заметки.
* Ждёт ответ.

**Что делает сервер**:

* Получает запрос, достаёт заметку из хранилища по ID.
* Если нашёл:

  * Отправляет `GetNoteResponse` с полем `note` (сама заметка).
* Если не нашёл:

  * Отправляет gRPC-ошибку `NotFound`.

**Что получает клиент**:

* Либо полную заметку, либо ошибку `NotFound`.

---

### 3. **Список всех заметок** (`ListNotes`)

**Что делает клиент**:

* Отправляет пустой `ListNotesRequest`.
* Ждёт ответ.

**Что делает сервер**:

* Берёт все заметки из хранилища.
* Формирует `ListNotesResponse`:

  * Поле `notes` — массив всех заметок.
* Отправляет клиенту.

**Что получает клиент**:

* Список заметок (каждая — `Note` с id, title, content, created\_at).

---

## Смысл обмена

По сути:

* **Клиент** — инициатор, он отправляет *запросы* с данными и ждёт *ответы*.
* **Сервер** — обработчик, он:

  * Принимает запрос,
  * Выполняет нужную бизнес-логику (создание/поиск/листинг),
  * Возвращает результат или ошибку.

В нашем примере:

* Когда мы создаём заметку — мы передаём серверу **неполную** информацию (только текст и заголовок).
* Сервер сам добавляет то, что клиент не знает: **ID** и **время создания**.
* Ответ — это **подтверждение создания** с полной информацией о заметке.

---

## Структура

```
.
├─ cmd/
│  ├─ server/                # сервер (composition root)
│  │  └─ main.go
│  └─ client/                # клиент (демо-вызовы)
│     └─ main.go
├─ pkg/
│  ├─ handler/
│  │  └─ grpc/
│  │     └─ note_handler.go  # входной адаптер: gRPC -> сервис
│  ├─ repository/
│  │  └─ memory/
│  │     └─ note_repo.go     # репозиторий в памяти (адаптер к сервисному интерфейсу)
│  └─ service/
│     └─ note_service.go     # логика
├─ pkg/pb/
│  └─ note.pb.go             # сгенерированный код protobuf/gRPC (не редактировать)
└─ proto/
   └─ note.proto             # gRPC контракт (protobuf)
```

Ключевые моменты:

* `pkg/service` — доменная модель `Note`, интерфейс `NoteRepository`, и реализация `NoteService` (use cases).
* `pkg/repository/memory` — конкретная реализация хранилища (in‑memory).
* `pkg/handler/grpc` — gRPC‑обработчик, маппит protobuf <-> доменные сущности и вызывает сервис.
* `cmd/server` — composition root: собирает зависимости (repo → service → handler), стартует gRPC.
* `cmd/client` — простой клиент, который вызывает методы сервиса.

---

## Требования

* Go 1.21+ (рекомендуется)
* protoc 25.x/26.x (желательно не ниже 21.x)
* Плагины:

  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

  Убедитесь, что `$(go env GOPATH)/bin` в `PATH`.

---

## Установка зависимостей

```bash
go mod tidy
```

---

## Генерация protobuf кода

Проверьте, что в `proto/note.proto` корректно указан `go_package` под ваш module path:

```proto
option go_package = "github.com/<your-username>/<your-repo>/grpc/pkg/pb;pb";
```

Затем сгенерируйте:

```bash
protoc -I proto \
  --go_out=pkg/pb --go_opt=paths=source_relative \
  --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
  proto/note.proto
```

> Если раньше уже генерировали, можно очистить:
> `rm -f pkg/pb/*.pb.go && go mod tidy && (re)generate`.

---

## Запуск

### Сервер

```bash
go run ./grpc/cmd/server
# Логи:
# 2025/08/11 18:44:58 gRPC server starting on :50051
```

Порт можно переопределить переменной окружения:

```bash
PORT=60000 go run ./grpc/cmd/server
```

### Клиент

В отдельном терминале:

```bash
go run ./grpc/cmd/client
# Вы увидите что-то вроде:
# Created: 8b250a24-... => First
# Created: 13b28b8f-... => Second
# GetNote: id:"8b250a24-..." title:"First" content:"Hello world" created_at:...
# ListNotes:
# - 8b250a24-... | First
# - 13b28b8f-... | Second
# GetNote(2): 13b28b8f-...: Second
```

---

## API (кратко)

Прото‑контракт: `proto/note.proto`

Методы:

* `CreateNote(CreateNoteRequest) -> CreateNoteResponse`
* `GetNote(GetNoteRequest) -> GetNoteResponse`
* `ListNotes(ListNotesRequest) -> ListNotesResponse`

Сообщения:

* `Note { id, title, content, created_at }`
* `CreateNoteRequest { title, content }`
* `GetNoteRequest { id }`

В хендлере есть маппинг ошибок сервиса на gRPC‑коды:

* `ErrBadRequest` → `InvalidArgument`
* `ErrNotFound`   → `NotFound`
* прочее          → `Internal`

---

## Как это связано «чистой архитектурой»

* **Domain/Application (pkg/service)**: бизнес‑модель (`Note`), порт хранилища (`NoteRepository`), и юзкейсы (`NoteService`).
  Никаких деталей транспорта или СУБД.
* **Adapters**:

    * **Repository**: `pkg/repository/memory` реализует `NoteRepository`.
    * **Transport**: `pkg/handler/grpc` реализует gRPC‑интерфейс и вызывает `NoteService`.
* **Composition root**: `cmd/server/main.go` связывает адаптеры, сервис и транспорт.
  Хотите поменять хранилище или добавить другой транспорт (HTTP) — меняете адаптер или композицию, не трогая бизнес‑логику.


---

## Частые ошибки и их решения

**Паника `slice bounds out of range` внутри `google.golang.org/protobuf/internal/filedesc`**
→ Несовместимые версии `protoc`/плагинов/рантайма.
Решение:

* Обновить `protoc` (до 25/26), `protoc-gen-go`, `protoc-gen-go-grpc`.
* Убедиться, что `option go_package` совпадает с вашим module path.
* Перегенерировать `.pb.go`.
* `go get google.golang.org/protobuf@latest`, `go get google.golang.org/grpc@latest`, `go mod tidy`.

**`undefined: grpc.NewServer` / `cannot find package`**
→ Неверные импорты.
Решение:

* Импорт транспорта: `import "google.golang.org/grpc"`.
* Путь к хендлеру/репозиторию соответствует вашей структуре/модулю.
* `go mod tidy`.

**Клиент не подключается к серверу**
→ Сервер не запущен или порт занят/закрыт.
Решение:

* Проверьте, что сервер слушает `:50051` (или ваш порт).
* Измените порт через `PORT` или `cmd/client` адрес.

---

## Пример «псевдо‑трассы» CreateNote

**Client:**

* Build CreateNoteRequest{title, content}
* Protobuf marshal → bytes
* HTTP/2 → Server

**Server:**
* HTTP/2 → bytes
* Protobuf unmarshal → *pb.CreateNoteRequest
* Handler: svc.Create(title, content)
* Service: validate → enrich (uuid, time) → repo.Save(note)
* Repo: store in map
* Service: return Note
* Handler: map to pb.Note → CreateNoteResponse
* Protobuf marshal → bytes
* HTTP/2 → Client

**Client:**
* bytes → Protobuf unmarshal → *pb.CreateNoteResponse
* print
