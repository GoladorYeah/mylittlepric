# 📋 План рефакторинга и оптимизации Backend

**Дата создания**: 11 ноября 2025
**Версия**: 1.0
**Проект**: MyLittlePrice Backend

---

## 📊 Общая информация

### Текущий стек технологий
- **Backend Framework**: Go + Fiber v2.52.9
- **ORM**: Ent v0.14.5
- **Databases**: PostgreSQL + Redis v9.16.0
- **AI/Search**: Google Gemini API v1.34.0, SerpAPI
- **Authentication**: JWT (golang-jwt/jwt v5.3.0)

### Архитектура
```
backend/
├── cmd/api/          # Точка входа
├── ent/              # Ent ORM (generated)
├── internal/
│   ├── app/          # Роутинг
│   ├── config/       # Конфигурация
│   ├── container/    # DI контейнер
│   ├── domain/       # Domain models
│   ├── handlers/     # HTTP handlers + WebSocket
│   ├── middleware/   # Auth middleware
│   ├── models/       # Бизнес-модели
│   ├── services/     # Бизнес-логика (21 сервис)
│   └── utils/        # Утилиты
└── migrations/       # DB migrations
```

### Статистика
- **Сервисов**: 21
- **Handlers**: 8
- **Endpoints**: ~20 REST + WebSocket
- **Ent Entities**: 5 (User, ChatSession, Message, SearchHistory, UserPreference)

---

## 🚨 Выявленные проблемы

### Критичность проблем

| № | Проблема | Файл | Серьезность | Тип |
|---|----------|------|-------------|-----|
| 1 | N+1 GetSession() вызовы | `processor.go` | 🔴 **КРИТИЧНО** | Performance |
| 2 | Неиспользуемый sqlx | `container.go` | 🔴 **КРИТИЧНО** | Dead Code |
| 3 | CleanupExpiredAnonymousHistory не вызывается | `search_history_service.go` | 🔴 **КРИТИЧНО** | Missing Logic |
| 4 | No Redis fallback для users | `auth_service.go` | 🟡 Средне | Reliability |
| 5 | SessionService too large (SRP) | `session.go` | 🟡 Средне | Architecture |
| 6 | Дублирование getUserBy* методов | `auth_service.go` | 🟡 Средне | Code Quality |
| 7 | Игнорирование критичных ошибок | `processor.go` | 🟡 Средне | Error Handling |
| 8 | Отсутствие индексов в БД | `ent/schema/` | 🟡 Средне | Performance |
| 9 | Debug логирование fmt.Printf | `processor.go` | 🟢 Низко | Code Quality |
| 10 | JSONB вместо таблиц | `ent/schema/` | 🟢 Низко | Design |
| 11 | Устаревший Gemini SDK | `go.mod` | 🟡 Средне | Dependencies |
| 12 | Неоптимальная конфигурация Redis | `container.go` | 🟢 Низко | Performance |

---

## 🎯 Этапы рефакторинга

---

## Этап 1: Критические исправления (1-2 дня)

**Цель**: Устранить критические проблемы, влияющие на производительность и стабильность

### 1.1 Удалить неиспользуемую зависимость sqlx

**Приоритет**: 🔴 Критичный
**Сложность**: Низкая (1-2 часа)
**Файлы**:
- `backend/internal/container/container.go`
- `backend/go.mod`

**Проблема**:
После миграции на Ent ORM, sqlx все еще импортируется и инициализируется, создавая лишний connection pool к PostgreSQL.

**Задачи**:
- [x] Удалить импорт `github.com/jmoiron/sqlx` из `container.go:11`
- [x] Удалить поле `DB *sqlx.DB` из `Container` struct (строка 25)
- [x] Удалить код инициализации sqlx (строки 73-82)
- [x] Обновить метод `Close()` - удалить закрытие sqlx DB (строки 224-228)
- [x] Запустить `go mod tidy` для очистки зависимостей
- [x] Протестировать запуск приложения

**Статус**: ✅ **ЗАВЕРШЕНО** (11 ноября 2025)

**Ожидаемый результат**:
- Освобождение лишнего connection pool (5 idle + 25 max connections)
- Уменьшение потребления памяти
- Очистка dead code

---

### 1.2 Добавить Cron job для очистки истории поиска

**Приоритет**: 🔴 Критичный
**Сложность**: Средняя (2-3 часа)
**Файлы**:
- `backend/internal/jobs/cleanup.go` (новый)
- `backend/cmd/api/main.go`

**Статус**: ✅ **ЗАВЕРШЕНО** (11 ноября 2025)

**Проблема**:
Функция `CleanupExpiredAnonymousHistory()` определена в `search_history_service.go:287`, но никогда не вызывается. Это приводит к бесконечному накоплению истории анонимных пользователей в БД.

**Задачи**:
- [x] Создать новый пакет `backend/internal/jobs/`
- [x] Создать файл `cleanup.go` с CleanupJob
- [x] Реализовать ticker для запуска каждые 24 часа
- [x] Вызывать `SearchHistoryService.CleanupExpiredAnonymousHistory()`
- [x] Добавить логирование результатов (сколько записей удалено)
- [x] Запустить goroutine в `main.go`
- [x] Добавить graceful shutdown для cleanup job
- [ ] Добавить метрики (опционально)

**Пример кода**:
```go
// backend/internal/jobs/cleanup.go
package jobs

import (
    "context"
    "log"
    "time"
    "mylittleprice/internal/services"
)

type CleanupJob struct {
    searchHistoryService *services.SearchHistoryService
    interval             time.Duration
    ctx                  context.Context
    cancel               context.CancelFunc
}

func NewCleanupJob(shs *services.SearchHistoryService) *CleanupJob {
    ctx, cancel := context.WithCancel(context.Background())
    return &CleanupJob{
        searchHistoryService: shs,
        interval:             24 * time.Hour,
        ctx:                  ctx,
        cancel:               cancel,
    }
}

func (j *CleanupJob) Start() {
    ticker := time.NewTicker(j.interval)
    go func() {
        for {
            select {
            case <-ticker.C:
                j.runCleanup()
            case <-j.ctx.Done():
                ticker.Stop()
                return
            }
        }
    }()
    log.Println("🧹 Cleanup job started (runs every 24h)")
}

func (j *CleanupJob) runCleanup() {
    count, err := j.searchHistoryService.CleanupExpiredAnonymousHistory(j.ctx)
    if err != nil {
        log.Printf("❌ Cleanup job failed: %v", err)
    } else {
        log.Printf("✅ Cleanup job completed: %d records deleted", count)
    }
}

func (j *CleanupJob) Stop() {
    j.cancel()
    log.Println("🛑 Cleanup job stopped")
}
```

**Интеграция в main.go**:
```go
// После инициализации container
cleanupJob := jobs.NewCleanupJob(container.SearchHistoryService)
cleanupJob.Start()

// В graceful shutdown
defer cleanupJob.Stop()
```

**Ожидаемый результат**:
- Автоматическая очистка старых данных каждые 24 часа
- Предотвращение неограниченного роста таблицы `search_history`
- Логи с информацией о количестве удаленных записей

**Реализованные улучшения**:
- ✅ Создан пакет `internal/jobs/` с CleanupJob
- ✅ Cleanup запускается немедленно при старте приложения, затем каждые 24 часа
- ✅ Добавлено подробное логирование: количество удаленных записей, время выполнения
- ✅ Graceful shutdown: cleanup job корректно останавливается при завершении приложения
- ✅ Context-based cancellation для безопасной остановки ticker

---

### 1.3 Миграция на новый Google Gemini SDK

**Приоритет**: 🔴 Критичный
**Сложность**: Средняя (4-6 часов)
**Дедлайн**: 30 ноября 2025
**Файлы**:
- `backend/go.mod`
- `backend/internal/services/gemini.go`
- `backend/internal/services/embedding.go`
- `backend/internal/container/container.go`

**Проблема**:
Текущий SDK `google.golang.org/genai v1.34.0` устарел. Google объявил о новом унифицированном SDK и прекращает поддержку старого **30 ноября 2025**.

**Старый SDK**: `github.com/google/generative-ai-go`
**Новый SDK**: `github.com/googleapis/go-genai` (пакет `google.golang.org/genai`)

**Задачи**:
- [ ] Изучить документацию нового SDK: https://github.com/googleapis/go-genai
- [ ] Изучить guide по миграции: https://ai.google.dev/gemini-api/docs/migrate
- [ ] Обновить `go.mod`: `go get google.golang.org/genai@latest`
- [ ] Рефакторить `GeminiService`:
  - [ ] Обновить инициализацию клиента
  - [ ] Обновить методы генерации текста
  - [ ] Обновить методы работы с grounding
  - [ ] Обновить обработку ошибок (новый формат)
- [ ] Рефакторить `EmbeddingService`:
  - [ ] Обновить методы генерации embeddings
  - [ ] Проверить совместимость размерности векторов
- [ ] Обновить инициализацию в `container.go`
- [ ] Протестировать все AI features:
  - [ ] Генерация ответов
  - [ ] Category classification
  - [ ] Product extraction
  - [ ] Embeddings для кэша
  - [ ] Grounding с Google Search
- [ ] Удалить старый SDK из зависимостей

**Breaking changes (ожидаемые)**:
- Новые импорты
- Изменения в структуре `ClientConfig`
- Новые методы API для content generation
- Изменения в обработке streaming responses

**Ожидаемый результат**:
- Поддержка новых моделей (Gemini 2.0, Veo, Imagen)
- Актуальные bug fixes и improvements
- Соответствие требованиям Google

---

## Этап 2: Улучшение производительности (3-5 дней)

**Цель**: Устранить N+1 проблемы и оптимизировать запросы к БД

### 2.1 Устранить N+1 проблему с GetSession() в ProcessChat

**Приоритет**: 🔴 Критичный
**Сложность**: Высокая (6-8 часов)
**Файлы**:
- `backend/internal/handlers/processor.go`
- `backend/internal/services/session.go`

**Проблема**:
Функция `ProcessChat()` вызывает `GetSession()` 5+ раз для обработки одного сообщения:
- Строка 85: `getOrCreateSession()`
- Строка 156: после добавления сообщения
- Строка 397: после добавления истории
- Строка 441: после цикла
- Косвенные вызовы через `AddToCycleHistory()`, `IncrementCycleIteration()`, `StartNewCycle()`

Каждый вызов = запрос к Redis + потенциальный fallback к PostgreSQL + десериализация JSON.

**Задачи**:

**Фаза 1: Рефакторинг SessionService**
- [ ] Изменить методы SessionService для работы с in-memory объектом:
  ```go
  // Было:
  func (s *SessionService) IncrementCycleIteration(sessionID string) error

  // Стало:
  func (s *SessionService) IncrementCycleIteration(session *models.Session) error
  ```
- [ ] Обновить методы:
  - [ ] `AddToCycleHistory(session *models.Session, ...)`
  - [ ] `IncrementCycleIteration(session *models.Session)`
  - [ ] `StartNewCycle(session *models.Session, ...)`
  - [ ] `AddMessage(session *models.Session, ...)`
  - [ ] `StartNewSearch(session *models.Session, ...)`

**Фаза 2: Рефакторинг ProcessChat**
- [ ] Получить session **один раз** в начале (строка 85)
- [ ] Передавать указатель на session во все методы
- [ ] Удалить все промежуточные вызовы GetSession()
- [ ] Сохранить session **один раз** в конце через `UpdateSession()`

**Фаза 3: Добавить явное сохранение**
- [ ] Создать метод `SaveSession(session *models.Session) error`
- [ ] Вызывать в конце ProcessChat для persist изменений
- [ ] Добавить оптимистическую блокировку (version field) для concurrency

**Фаза 4: Тестирование**
- [ ] Unit tests для обновленных методов SessionService
- [ ] Integration test для ProcessChat
- [ ] Измерить количество Redis/PostgreSQL запросов (до/после)
- [ ] Load testing для проверки concurrency

**Ожидаемый результат**:
- **5x сокращение** запросов к Redis/PostgreSQL на сообщение
- Уменьшение latency обработки сообщений на 30-50%
- Более предсказуемое поведение при высокой нагрузке

---

### 2.2 Добавить индексы в Ent schema

**Приоритет**: 🟡 Средний
**Сложность**: Средняя (3-4 часа)
**Файлы**:
- `backend/ent/schema/user.go`
- `backend/ent/schema/searchhistory.go`
- `backend/ent/schema/chatsession.go`
- `backend/migrations/`

**Проблема**:
Отсутствуют индексы на часто используемых полях, что замедляет queries:
- `user.email` - поиск при логине
- `search_history.user_id`, `session_id` - фильтрация истории
- `search_history.created_at` - сортировка
- `chat_session.expires_at` - cleanup запросы

**Задачи**:

**User schema**:
```go
// backend/ent/schema/user.go
func (User) Indexes() []ent.Index {
    return []ent.Index{
        // Email уже unique, но добавим для поиска
        index.Fields("email"),
        // Для поиска по provider
        index.Fields("provider", "google_id"),
    }
}
```

**SearchHistory schema**:
```go
// backend/ent/schema/searchhistory.go
func (SearchHistory) Indexes() []ent.Index {
    return []ent.Index{
        // Для GetUserSearchHistory
        index.Fields("user_id", "created_at"),
        // Для анонимных пользователей
        index.Fields("session_id", "created_at"),
        // Для cleanup job
        index.Fields("expires_at").
            Annotations(entsql.IndexWhere("user_id IS NULL")), // partial index
    }
}
```

**ChatSession schema**:
```go
// backend/ent/schema/chatsession.go
func (ChatSession) Indexes() []ent.Index {
    return []ent.Index{
        // Для GetActiveSessionForUser
        index.Fields("user_id", "expires_at"),
        // Для cleanup
        index.Fields("expires_at"),
    }
}
```

**Migration steps**:
- [ ] Добавить метод `Indexes()` в каждую schema
- [ ] Сгенерировать миграцию: `go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema`
- [ ] Создать SQL миграцию для существующей БД
- [ ] Протестировать на dev окружении
- [ ] Измерить производительность запросов (EXPLAIN ANALYZE)
- [ ] Применить на production

**Ожидаемый результат**:
- **10-100x ускорение** поисковых запросов
- Оптимизация cleanup операций
- Улучшение производительности при росте данных

---

## Этап 3: Улучшение надежности (3-5 дней)

**Цель**: Добавить fallback механизмы и правильную обработку ошибок

### 3.1 Добавить Redis fallback для getUserByID

**Приоритет**: 🟡 Средний
**Сложность**: Средняя (3-4 часа)
**Файлы**:
- `backend/internal/services/auth_service.go`

**Проблема**:
Методы `getUserByID()`, `getUserByEmail()`, `getUserByProviderID()` полагаются только на Redis. Если Redis недоступен или данные отсутствуют, система падает. Это критично для:
- `RefreshAccessToken()` - нельзя обновить токен
- `GoogleLogin()` - OAuth не работает
- `Login()` - обычный вход не работает

**Текущий код** (строки 452-485):
```go
func (s *AuthService) getUserByID(userID uuid.UUID) (*models.User, error) {
    userData, err := s.redis.HGetAll(s.ctx, userKey).Result()
    if err != nil {
        return nil, err  // ← Redis упал = ошибка
    }
    if len(userData) == 0 {
        return nil, redis.Nil  // ← Нет в Redis = ошибка
    }
    // ...
}
```

**Задачи**:

**Фаза 1: Добавить fallback к Ent**
```go
func (s *AuthService) getUserByID(userID uuid.UUID) (*models.User, error) {
    // 1. Попробовать Redis
    userKey := fmt.Sprintf("user:id:%s", userID.String())
    userData, err := s.redis.HGetAll(s.ctx, userKey).Result()

    if err == nil && len(userData) > 0 {
        // Redis hit - parse и return
        return s.parseUserFromRedis(userData)
    }

    // 2. Fallback к PostgreSQL через Ent
    log.Printf("⚠️ Redis miss for user %s, falling back to PostgreSQL", userID)
    entUser, err := s.entClient.User.Get(s.ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    // 3. Sync обратно в Redis для следующих запросов
    user := s.entUserToModel(entUser)
    if err := s.syncUserToRedis(user); err != nil {
        log.Printf("⚠️ Failed to sync user to Redis: %v", err)
        // Не возвращаем ошибку - user получен из БД
    }

    return user, nil
}
```

**Фаза 2: Создать helper методы**
- [ ] `parseUserFromRedis(data map[string]string) (*models.User, error)`
- [ ] `entUserToModel(entUser *ent.User) *models.User`
- [ ] `syncUserToRedis(user *models.User) error`

**Фаза 3: Обновить другие методы**
- [ ] Обновить `getUserByEmail()` для fallback
- [ ] Обновить `getUserByProviderID()` для fallback
- [ ] Добавить аналогичную логику в другие get методы

**Фаза 4: Тестирование**
- [ ] Unit test: Redis доступен, пользователь есть
- [ ] Unit test: Redis недоступен, fallback к PostgreSQL
- [ ] Unit test: Пользователь не найден нигде
- [ ] Integration test: Полный сценарий login/refresh
- [ ] Chaos testing: Отключить Redis и проверить работу

**Ожидаемый результат**:
- Система работает даже при отказе Redis
- Автоматическая синхронизация данных
- Улучшение resilience

---

### 3.2 Правильная обработка ошибок в processor

**Приоритет**: 🟡 Средний
**Сложность**: Средняя (4-5 часов)
**Файлы**:
- `backend/internal/handlers/processor.go`

**Проблема**:
Множество критичных ошибок игнорируются через `fmt.Printf()`:

```go
// Строка 101: Ошибка игнорируется полностью
session, _ = p.container.SessionService.GetSession(req.SessionID)

// Строка 397: Логируется но продолжает работу
if err != nil {
    fmt.Printf("⚠️ Failed to re-fetch session: %v\n", err)
    // session может быть устаревшей!
}

// Строка 422: Критичная ошибка игнорируется
if err := p.container.SessionService.IncrementCycleIteration(req.SessionID); err != nil {
    fmt.Printf("⚠️ Failed to increment cycle: %v\n", err)
    // Цикл не обновился!
}
```

**Задачи**:

**Фаза 1: Классификация ошибок**
- [ ] Определить критичные ошибки (должны прервать обработку)
- [ ] Определить recoverable ошибки (можно retry)
- [ ] Определить non-critical ошибки (только логировать)

**Фаза 2: Добавить retry logic**
```go
// backend/internal/utils/retry.go
func RetryWithBackoff(fn func() error, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }

        if !isRetriable(err) {
            return err
        }

        backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
        time.Sleep(backoff)
    }
    return fmt.Errorf("max retries exceeded: %w", err)
}
```

**Фаза 3: Обновить processor.go**
- [ ] Заменить все игнорируемые ошибки на proper handling
- [ ] Добавить retry для Redis/PostgreSQL операций
- [ ] Возвращать ошибки для критичных failures
- [ ] Добавить context timeout для long operations

**Пример рефакторинга**:
```go
// Было:
session, _ = p.container.SessionService.GetSession(req.SessionID)

// Стало:
session, err = p.container.SessionService.GetSession(req.SessionID)
if err != nil {
    return &ChatProcessorResponse{
        Error: &ErrorInfo{
            Code:    "session_fetch_failed",
            Message: "Failed to get session",
            Details: err.Error(),
        },
    }
}
```

**Фаза 4: Добавить circuit breaker для внешних сервисов**
```go
// Для Gemini API, SERP API
type CircuitBreaker struct {
    failureThreshold int
    resetTimeout     time.Duration
    state            string // "closed", "open", "half-open"
}
```

**Ожидаемый результат**:
- Надежная обработка временных сбоев
- Предсказуемое поведение при ошибках
- Лучшая observability проблем

---

## Этап 4: Рефакторинг архитектуры (5-7 дней)

**Цель**: Улучшить разделение ответственности и maintainability

### 4.1 Разделить SessionService на 3 сервиса

**Приоритет**: 🟡 Средний
**Сложность**: Высокая (8-12 часов)
**Файлы**:
- `backend/internal/services/session.go` → разделить
- `backend/internal/services/message.go` (новый)
- `backend/internal/services/cycle.go` (новый)
- `backend/internal/container/container.go`
- `backend/internal/handlers/*.go`

**Проблема**:
`SessionService` имеет ~25 методов и нарушает Single Responsibility Principle. Один сервис управляет:
1. Сессиями (Redis + PostgreSQL)
2. Сообщениями (Redis)
3. Циклами промптов (Redis)
4. Состоянием поиска (JSONB)
5. Контекстом разговора (JSONB)

**Новая архитектура**:

```
SessionService (10 методов)
├── CreateSession()
├── GetSession()
├── UpdateSession()
├── DeleteSession()
├── GetActiveSessionForUser()
├── LinkSessionToUser()
├── StartNewSearch()
├── SetCategory()
├── IsSearchCompleted()
└── GetSessionInfo()

MessageService (5 методов)
├── AddMessage()
├── GetMessages()
├── GetRecentMessages()
├── GetConversationHistory()
└── IncrementMessageCount()

CycleService (4 методов)
├── IncrementCycleIteration()
├── StartNewCycle()
├── AddToCycleHistory()
└── GetConversationContext()
```

**Задачи**:

**Фаза 1: Создать MessageService**
- [ ] Создать `backend/internal/services/message.go`
- [ ] Переместить методы работы с сообщениями
- [ ] Обновить зависимости (Redis, Config)
- [ ] Добавить unit tests

**Фаза 2: Создать CycleService**
- [ ] Создать `backend/internal/services/cycle.go`
- [ ] Переместить методы работы с циклами
- [ ] Добавить зависимость от SessionService (для получения session)
- [ ] Добавить unit tests

**Фаза 3: Упростить SessionService**
- [ ] Удалить перемещенные методы
- [ ] Оставить только core функциональность
- [ ] Обновить существующие tests

**Фаза 4: Обновить Container**
```go
// backend/internal/container/container.go
type Container struct {
    // ...
    SessionService *services.SessionService
    MessageService *services.MessageService  // Новый
    CycleService   *services.CycleService    // Новый
    // ...
}

func (c *Container) initServices() error {
    // Сначала SessionService (базовый)
    c.SessionService = services.NewSessionService(...)

    // Потом зависимые сервисы
    c.MessageService = services.NewMessageService(c.Redis, c.Config)
    c.CycleService = services.NewCycleService(c.SessionService, c.Redis)
    // ...
}
```

**Фаза 5: Обновить все handlers**
- [ ] processor.go - использовать MessageService, CycleService
- [ ] chat.go - обновить вызовы
- [ ] websocket.go - обновить вызовы
- [ ] Обновить integration tests

**Ожидаемый результат**:
- Лучшее разделение ответственности (SRP)
- Проще тестировать каждый сервис отдельно
- Легче добавлять новую функциональность
- Меньший размер каждого файла (лучше читаемость)

---

### 4.2 Объединить и упростить getUserBy* методы

**Приоритет**: 🟢 Низкий
**Сложность**: Средняя (2-3 часа)
**Файлы**:
- `backend/internal/services/auth_service.go`

**Проблема**:
Дублирование кода в 3 методах:
- `getUserByID()` (строка 452)
- `getUserByEmail()` (строка 437) → вызывает getUserByID
- `getUserByProviderID()` (строка 487) → вызывает getUserByID

**Задачи**:

**Создать унифицированный метод**:
```go
type UserLookup struct {
    ByID         *uuid.UUID
    ByEmail      *string
    ByProviderID *struct {
        Provider string
        ID       string
    }
}

func (s *AuthService) getUser(lookup UserLookup) (*models.User, error) {
    var userID uuid.UUID
    var err error

    // 1. Determine userID from lookup criteria
    switch {
    case lookup.ByID != nil:
        userID = *lookup.ByID
    case lookup.ByEmail != nil:
        userID, err = s.lookupUserIDByEmail(*lookup.ByEmail)
    case lookup.ByProviderID != nil:
        userID, err = s.lookupUserIDByProvider(lookup.ByProviderID.Provider, lookup.ByProviderID.ID)
    default:
        return nil, fmt.Errorf("no lookup criteria provided")
    }

    if err != nil {
        return nil, err
    }

    // 2. Get user with Redis fallback (из Этапа 3.1)
    return s.getUserWithFallback(userID)
}
```

**Обновить публичные методы**:
```go
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
    return s.getUser(UserLookup{ByID: &userID})
}

func (s *AuthService) getUserByEmail(email string) (*models.User, error) {
    return s.getUser(UserLookup{ByEmail: &email})
}

func (s *AuthService) getUserByProviderID(provider, providerID string) (*models.User, error) {
    return s.getUser(UserLookup{
        ByProviderID: &struct{Provider, ID string}{provider, providerID},
    })
}
```

**Ожидаемый результат**:
- Единая точка для получения пользователей
- Консистентное поведение (fallback, caching)
- Меньше дублирования кода

---

### 4.3 Заменить JSONB на отдельные таблицы

**Приоритет**: 🟢 Низкий
**Сложность**: Очень высокая (16-24 часа)
**Файлы**:
- `backend/ent/schema/searchstate.go` (новый)
- `backend/ent/schema/cyclestate.go` (новый)
- `backend/ent/schema/chatsession.go`
- `backend/internal/services/session.go`
- `backend/migrations/`

**Проблема**:
JSONB поля `search_state` и `cycle_state` в `ChatSession`:
- Теряют типизацию
- Требуют постоянного преобразования `mapToStruct()`/`structToMap()`
- Сложно делать queries по вложенным полям
- Нет foreign key constraints

**Новая структура БД**:

```
chat_sessions (1) ──> (1) search_states
chat_sessions (1) ──> (1) cycle_states
```

**Задачи**:

**Фаза 1: Создать новые Ent schemas**

`backend/ent/schema/searchstate.go`:
```go
type SearchState struct {
    ent.Schema
}

func (SearchState) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).Default(uuid.New),
        field.String("category").Optional(),
        field.String("search_phrase").Optional(),
        field.Bool("is_completed").Default(false),
        field.JSON("products_found", []map[string]interface{}{}).Default([]map[string]interface{}{}),
        field.Time("search_started_at").Optional(),
        field.Time("search_completed_at").Optional(),
        field.Time("created_at").Default(time.Now),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}

func (SearchState) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("session", ChatSession.Type).
            Ref("search_state").
            Unique().
            Required(),
    }
}
```

`backend/ent/schema/cyclestate.go`:
```go
type CycleState struct {
    ent.Schema
}

func (CycleState) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).Default(uuid.New),
        field.Int("current_iteration").Default(0),
        field.Int("max_iterations").Default(3),
        field.JSON("cycle_history", []map[string]interface{}{}).Default([]map[string]interface{}{}),
        field.JSON("context", map[string]interface{}{}).Default(map[string]interface{}{}),
        field.Time("created_at").Default(time.Now),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}

func (CycleState) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("session", ChatSession.Type).
            Ref("cycle_state").
            Unique().
            Required(),
    }
}
```

**Фаза 2: Обновить ChatSession schema**
```go
// Удалить JSONB поля
// field.JSON("search_state", ...)
// field.JSON("cycle_state", ...)

// Добавить edges
func (ChatSession) Edges() []ent.Edge {
    return []ent.Edge{
        // Existing edges...
        edge.To("search_state", SearchState.Type).Unique(),
        edge.To("cycle_state", CycleState.Type).Unique(),
    }
}
```

**Фаза 3: Создать data migration**
```go
// backend/migrations/migrate_jsonb_to_tables.go
func MigrateJSONBToTables(ctx context.Context, client *ent.Client) error {
    sessions, _ := client.ChatSession.Query().All(ctx)

    for _, session := range sessions {
        // Parse JSONB
        var searchState models.SearchState
        mapToStruct(session.SearchState, &searchState)

        // Create SearchState entity
        client.SearchState.Create().
            SetSession(session).
            SetCategory(searchState.Category).
            // ... другие поля
            Save(ctx)

        // То же для CycleState
    }
}
```

**Фаза 4: Обновить SessionService**
- [ ] Убрать `structToMap()`/`mapToStruct()` вызовы
- [ ] Использовать Ent queries для работы с состояниями
- [ ] Обновить методы для eager loading (WithSearchState, WithCycleState)

**Фаза 5: Миграция на production**
- [ ] Тестировать миграцию на копии production БД
- [ ] Создать rollback план
- [ ] Запустить миграцию в maintenance window
- [ ] Удалить старые JSONB колонки после проверки

**Ожидаемый результат**:
- Типобезопасность
- Возможность делать SQL queries по состояниям
- Foreign key constraints
- Проще добавлять новые поля

---

## Этап 5: Улучшение качества кода (2-3 дня)

**Цель**: Внедрить best practices и улучшить observability

### 5.1 Внедрить structured logging

**Приоритет**: 🟡 Средний
**Сложность**: Средняя (4-6 часов)
**Файлы**:
- `backend/internal/utils/logger.go` (новый)
- Все файлы с `fmt.Printf()` и `log.Printf()`

**Проблема**:
- 34+ вызовов `fmt.Printf()` с эмодзи в processor.go
- Неструктурированные логи затрудняют анализ
- Нет context-aware логирования (request_id, user_id)
- Невозможно фильтровать по уровням (debug/info/error)

**Выбор библиотеки**:
**Рекомендация**: `log/slog` (стандартная библиотека Go 1.21+)

Альтернативы:
- `github.com/rs/zerolog` (быстрее, но внешняя зависимость)
- `go.uber.org/zap` (очень быстрая, но сложнее setup)

**Задачи**:

**Фаза 1: Создать logger wrapper**
```go
// backend/internal/utils/logger.go
package utils

import (
    "context"
    "log/slog"
    "os"
)

type ContextKey string

const (
    RequestIDKey ContextKey = "request_id"
    UserIDKey    ContextKey = "user_id"
    SessionIDKey ContextKey = "session_id"
)

var logger *slog.Logger

func InitLogger(level string) {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }

    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
    })

    logger = slog.New(handler)
}

func GetLogger() *slog.Logger {
    return logger
}

// Context-aware logging
func LogInfo(ctx context.Context, msg string, args ...any) {
    logger.InfoContext(ctx, msg, extractContextAttrs(ctx, args)...)
}

func LogError(ctx context.Context, msg string, err error, args ...any) {
    attrs := append(extractContextAttrs(ctx, args), slog.Any("error", err))
    logger.ErrorContext(ctx, msg, attrs...)
}

func extractContextAttrs(ctx context.Context, args []any) []any {
    attrs := make([]any, 0, len(args)+3)

    if reqID := ctx.Value(RequestIDKey); reqID != nil {
        attrs = append(attrs, slog.String("request_id", reqID.(string)))
    }
    if userID := ctx.Value(UserIDKey); userID != nil {
        attrs = append(attrs, slog.String("user_id", userID.(string)))
    }
    if sessionID := ctx.Value(SessionIDKey); sessionID != nil {
        attrs = append(attrs, slog.String("session_id", sessionID.(string)))
    }

    return append(attrs, args...)
}
```

**Фаза 2: Инициализация в main.go**
```go
// backend/cmd/api/main.go
func main() {
    // Load config
    cfg := config.Load()

    // Initialize logger
    utils.InitLogger(cfg.LogLevel)
    logger := utils.GetLogger()

    logger.Info("Starting MyLittlePrice Backend",
        slog.String("version", "1.0.0"),
        slog.String("env", cfg.Environment),
    )
    // ...
}
```

**Фаза 3: Добавить request context middleware**
```go
// backend/internal/middleware/request_context.go
func RequestContextMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        requestID := uuid.New().String()

        ctx := c.UserContext()
        ctx = context.WithValue(ctx, utils.RequestIDKey, requestID)

        // Add user_id if authenticated
        if userID := c.Locals("userID"); userID != nil {
            ctx = context.WithValue(ctx, utils.UserIDKey, userID)
        }

        c.SetUserContext(ctx)
        c.Set("X-Request-ID", requestID)

        return c.Next()
    }
}
```

**Фаза 4: Заменить все fmt.Printf()**

Примеры замены:
```go
// Было:
fmt.Printf("🔄 New search for session %s\n", req.SessionID)

// Стало:
logger.Info("new search started",
    slog.String("session_id", req.SessionID),
)

// Было:
fmt.Printf("⚠️ Failed to increment cycle: %v\n", err)

// Стало:
logger.Error("failed to increment cycle",
    slog.Any("error", err),
    slog.String("session_id", req.SessionID),
)
```

**Фаза 5: Обновить все сервисы**
- [ ] processor.go (34 замены)
- [ ] session.go
- [ ] auth_service.go
- [ ] gemini.go
- [ ] Все остальные сервисы с логами

**Фаза 6: Конфигурация через env**
```env
# .env
LOG_LEVEL=info  # debug, info, warn, error
LOG_FORMAT=json # json, text
```

**Ожидаемый результат**:
```json
{
  "time": "2025-11-11T15:30:00Z",
  "level": "INFO",
  "msg": "new search started",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "session_id": "abc123",
  "category": "smartphones"
}
```

**Преимущества**:
- Structured logs для анализа (ELK, Grafana Loki)
- Context propagation (request_id через весь flow)
- Фильтрация по уровням
- Production-ready format

---

### 5.2 Оптимизация Redis конфигурации

**Приоритет**: 🟢 Низкий
**Сложность**: Низкая (1-2 часа)
**Файлы**:
- `backend/internal/container/container.go`
- `backend/go.mod`

**Проблема**:
Используется базовая конфигурация Redis без оптимизаций для high-throughput приложений.

**Задачи**:

**Фаза 1: Обновить go-redis**
```bash
go get github.com/redis/go-redis/v9@latest
go mod tidy
```

**Фаза 2: Оптимизировать connection**
```go
// backend/internal/container/container.go
func (c *Container) initRedis() error {
    c.Redis = redis.NewClient(&redis.Options{
        Addr:     c.Config.RedisURL,
        Password: c.Config.RedisPassword,
        DB:       c.Config.RedisDB,

        // Connection pool
        PoolSize:     50,               // Увеличить для high-throughput
        MinIdleConns: 10,               // Поддерживать минимум idle
        MaxIdleConns: 20,               // Максимум idle

        // Timeouts
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
        PoolTimeout:  4 * time.Second,

        // Buffers (go-redis v9.12+)
        ReadBufferSize:  1024 * 1024,   // 1MiB для high-throughput
        WriteBufferSize: 1024 * 1024,   // 1MiB

        // Retry
        MaxRetries:      3,
        MinRetryBackoff: 8 * time.Millisecond,
        MaxRetryBackoff: 512 * time.Millisecond,

        // Maintenance notifications
        MaintNotificationsConfig: &maintnotifications.Config{
            Mode: maintnotifications.ModeDisabled,
        },
    })

    // Health check with context
    ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
    defer cancel()

    if err := c.Redis.Ping(ctx).Err(); err != nil {
        return fmt.Errorf("Redis ping failed: %w", err)
    }

    log.Println("✅ Connected to Redis with optimized config")
    return nil
}
```

**Фаза 3: Добавить OpenTelemetry (опционально)**
```go
import (
    "github.com/redis/go-redis/extra/redisotel/v9"
)

func (c *Container) initRedis() error {
    // ... создание клиента

    // Enable instrumentation
    if err := redisotel.InstrumentTracing(c.Redis); err != nil {
        log.Printf("⚠️ Failed to enable Redis tracing: %v", err)
    }

    if err := redisotel.InstrumentMetrics(c.Redis); err != nil {
        log.Printf("⚠️ Failed to enable Redis metrics: %v", err)
    }

    return nil
}
```

**Фаза 4: Добавить конфигурацию через env**
```go
// backend/internal/config/config.go
type Config struct {
    // ...
    RedisPoolSize     int           `env:"REDIS_POOL_SIZE" envDefault:"50"`
    RedisMinIdle      int           `env:"REDIS_MIN_IDLE" envDefault:"10"`
    RedisReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"3s"`
    RedisWriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"3s"`
}
```

**Ожидаемый результат**:
- Лучшая производительность Redis операций
- Правильные timeouts для production
- Connection pooling для concurrency
- Observability через OpenTelemetry (опционально)

---

## 📅 Временные рамки и приоритизация

### Рекомендуемый порядок выполнения:

**Неделя 1: Критические исправления**
- День 1-2: Этап 1 (Удалить sqlx, Cleanup job, Gemini SDK)
- День 3-5: Этап 2.1 (N+1 проблема GetSession)

**Неделя 2: Производительность и надежность**
- День 1-2: Этап 2.2 (Добавить индексы)
- День 3-4: Этап 3.1 (Redis fallback)
- День 5: Этап 3.2 (Обработка ошибок)

**Неделя 3: Рефакторинг архитектуры**
- День 1-3: Этап 4.1 (Разделить SessionService)
- День 4: Этап 4.2 (Объединить getUserBy*)
- День 5: Этап 5.1 (Structured logging)

**Неделя 4: Долгосрочные улучшения (опционально)**
- День 1-5: Этап 4.3 (JSONB → таблицы) - если нужно

---

## 📈 Метрики успеха

После завершения рефакторинга ожидаются следующие улучшения:

### Performance
- ✅ **5x** сокращение запросов к БД на сообщение
- ✅ **10-100x** ускорение queries с индексами
- ✅ **30-50%** уменьшение latency обработки сообщений
- ✅ Освобождение лишнего connection pool

### Reliability
- ✅ Работа при отказе Redis (fallback к PostgreSQL)
- ✅ Правильная обработка ошибок с retry logic
- ✅ Автоматическая очистка устаревших данных
- ✅ Circuit breaker для внешних API

### Maintainability
- ✅ Лучшее разделение ответственности (SRP)
- ✅ Меньше дублирования кода
- ✅ Structured logging для debugging
- ✅ Меньший размер файлов (проще читать)

### Code Quality
- ✅ Удаление dead code (sqlx)
- ✅ Актуальные зависимости (Gemini SDK)
- ✅ Type safety (убрать JSONB)
- ✅ Лучшее покрытие тестами

---

## ⚠️ Риски и митигация

### Риск 1: Миграция Gemini SDK может сломать AI features
**Митигация**:
- Тщательное тестирование всех AI features
- Поэтапный rollout (canary deployment)
- Подготовить rollback план

### Риск 2: Рефакторинг SessionService может ввести баги
**Митигация**:
- Полное покрытие тестами до рефакторинга
- Feature flags для постепенного переключения
- Extensive integration testing

### Риск 3: JSONB → таблицы требует data migration
**Митигация**:
- Тестировать миграцию на копии production БД
- Maintenance window для миграции
- Rollback план с восстановлением JSONB

### Риск 4: Изменение обработки ошибок может изменить поведение API
**Митигация**:
- Документировать новые error codes
- Обратная совместимость для clients
- Graceful degradation

---

## 🔗 Дополнительные ресурсы

### Документация библиотек
- **Ent ORM**: https://entgo.io/docs/getting-started
- **Fiber v2**: https://docs.gofiber.io/
- **Go-Redis v9**: https://redis.io/docs/latest/develop/clients/go/
- **Google GenAI SDK**: https://github.com/googleapis/go-genai
- **log/slog**: https://pkg.go.dev/log/slog

### Best Practices
- **Go Error Handling**: https://go.dev/blog/error-handling-and-go
- **Database Indexing**: https://use-the-index-luke.com/
- **Redis Best Practices**: https://redis.io/docs/latest/develop/use/patterns/
- **Structured Logging**: https://www.honeycomb.io/blog/structured-logging-and-your-team

---

## 📝 Чеклист готовности к началу

Перед началом убедитесь:
- [ ] Есть полный бэкап production БД
- [ ] Настроен staging environment для тестирования
- [ ] Все разработчики ознакомлены с планом
- [ ] Подготовлены тестовые сценарии для каждого этапа
- [ ] Настроен monitoring для отслеживания метрик
- [ ] Согласован maintenance window для критичных изменений

---

**Составлено**: Claude AI
**Дата**: 11 ноября 2025
**Версия плана**: 1.0
