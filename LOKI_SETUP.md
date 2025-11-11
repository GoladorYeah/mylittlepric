# Настройка логирования с Grafana Loki

## Обзор

Проект поддерживает два способа отправки логов в Loki:

1. **Promtail** - для сбора логов из Docker контейнеров (postgres, redis, и т.д.)
2. **Прямая отправка из Backend** - backend отправляет логи напрямую в Loki через HTTP API

## Способ 1: Promtail для Docker контейнеров

### Описание
Promtail собирает логи из всех Docker контейнеров и отправляет их в Loki.

### Архитектура:
```
Docker Контейнеры (postgres, redis, и т.д.)
         ↓ (логи)
    Promtail (собирает логи)
         ↓ (отправка)
      Loki (хранит логи)
         ↓ (запросы)
     Grafana (визуализация)
```

### Запуск:

```bash
# Остановить текущие сервисы мониторинга (если запущены)
docker compose -f docker-compose.monitoring.yml down

# Запустить сервисы
docker compose -f docker-compose.monitoring.yml up -d

# Проверить логи Promtail
docker logs mylittleprice_promtail

# Проверить статус
docker compose -f docker-compose.monitoring.yml ps
```

### Просмотр логов в Grafana:

1. Откройте Grafana: http://localhost:3001
2. Логин: `admin` / Пароль: `admin` (или из `GRAFANA_ADMIN_PASSWORD`)
3. Перейдите в `Explore` (иконка компаса)
4. Выберите источник данных `Loki`
5. Используйте запросы LogQL:

```logql
# Все логи Docker контейнеров
{job="docker_all"}

# Логи конкретного сервиса
{service="postgres"}
{service="redis"}

# Логи по имени контейнера
{container="mylittleprice-postgres"}

# Поиск по тексту
{job="docker_all"} |= "error"
```

---

## Способ 2: Прямая отправка из Backend (НОВОЕ)

### Описание
Backend приложение отправляет логи напрямую в Loki через HTTP API без использования Promtail. Логи пишутся одновременно в `stdout` и в Loki.

### Архитектура:

```
Backend Application
       ↓
logger (slog)
       ↓
io.MultiWriter
       ├→ stdout (терминал)
       └→ LokiWriter
              ├→ Buffer (буферизация)
              └→ HTTP POST → Loki API
                                ↓
                             Grafana
```

### Особенности:

- ✅ **Буферизация**: Логи батчуются (100 записей) перед отправкой
- ✅ **Периодическая отправка**: Буфер сбрасывается каждые 5 секунд
- ✅ **Graceful shutdown**: При остановке все логи отправляются в Loki
- ✅ **Отказоустойчивость**: Ошибки Loki не ломают приложение
- ✅ **Контекстная информация**: request_id, user_id, session_id автоматически добавляются

### Настройка:

#### 1. Добавьте переменные в `.env`:

```bash
# Loki Configuration
LOKI_ENABLED=true
LOKI_URL=http://localhost:3100/loki/api/v1/push
LOKI_SERVICE_NAME=mylittleprice-backend
```

**Параметры:**
- `LOKI_ENABLED` (default: `false`) - включить/выключить отправку в Loki
- `LOKI_URL` (default: `http://localhost:3100/loki/api/v1/push`) - URL Loki Push API
- `LOKI_SERVICE_NAME` (default: `mylittleprice-backend`) - название сервиса для лейблов

#### 2. Запустите Loki (если еще не запущен):

```bash
docker compose -f docker-compose.monitoring.yml up -d loki grafana
```

#### 3. Запустите backend:

```bash
cd backend
go run cmd/api/main.go
```

При запуске вы увидите:
```
level=INFO msg="Starting MyLittlePrice Backend" ... loki_enabled=true loki_url=http://localhost:3100/loki/api/v1/push
```

### Просмотр логов Backend в Grafana:

#### Запросы LogQL:

**Все логи backend:**
```logql
{service="mylittleprice-backend"}
```

**Логи по уровню:**
```logql
{service="mylittleprice-backend"} |= "level=error"
{service="mylittleprice-backend"} |= "level=info"
```

**Логи с request_id:**
```logql
{service="mylittleprice-backend"} | json | request_id != ""
```

**Ошибки за последний час:**
```logql
{service="mylittleprice-backend"} |= "level=error" [1h]
```

**Поиск по тексту:**
```logql
{service="mylittleprice-backend"} |= "database"
{service="mylittleprice-backend"} |= "redis"
{service="mylittleprice-backend"} |= "Failed"
```

**Статистика по запросам:**
```logql
rate({service="mylittleprice-backend"} [5m])
```

### Структура логов Backend:

**Лейблы Loki:**
- `service` - mylittleprice-backend
- `job` - mylittleprice-backend
- `level` - info/debug/warn/error

**JSON поля:**
- `time` - timestamp
- `level` - уровень логирования
- `msg` - сообщение
- `request_id` - ID запроса (автоматически)
- `user_id` - ID пользователя (автоматически)
- `session_id` - ID сессии (автоматически)
- Дополнительные контекстные поля

### Отключение Loki:

```bash
LOKI_ENABLED=false
```

Логи продолжат писаться в `stdout`.

---

## Troubleshooting

### Логи не появляются в Grafana

1. **Проверьте, что Loki запущен:**
   ```bash
   curl http://localhost:3100/ready
   # Должен вернуть "ready"
   ```

2. **Проверьте логи backend:**
   ```bash
   # Должно быть "loki_enabled=true"
   cd backend && go run cmd/api/main.go
   ```

3. **Проверьте переменные окружения:**
   ```bash
   cat backend/.env | grep LOKI
   ```

### Ошибки подключения к Loki

**Backend не может подключиться к Loki:**

1. Убедитесь, что URL правильный:
   - Локальный запуск: `http://localhost:3100/loki/api/v1/push`
   - В Docker: `http://loki:3100/loki/api/v1/push`

2. Проверьте доступность:
   ```bash
   docker compose -f docker-compose.monitoring.yml ps
   curl http://localhost:3100/ready
   ```

### Promtail не собирает логи

1. **Проверьте статус Promtail:**
   ```bash
   docker ps | grep promtail
   docker logs mylittleprice_promtail
   ```

2. **Проверьте метрики:**
   ```bash
   curl http://localhost:9080/metrics
   ```

3. **Проверьте права доступа к Docker socket:**
   ```bash
   docker exec mylittleprice_promtail ls -la /var/run/docker.sock
   ```

---

## Дополнительные возможности

### Добавление кастомных лейблов Backend

Отредактируйте `backend/internal/utils/logger.go`:

```go
labels := map[string]string{
    "service":     serviceName,
    "job":         serviceName,
    "level":       level,
    "environment": os.Getenv("ENV"),     // окружение
    "version":     "1.0.0",              // версия
    "hostname":    os.Hostname(),        // хост
}
```

### Создание дашбордов в Grafana

1. **Dashboards** → **New Dashboard**
2. Добавьте панели с запросами:

**Панель 1: Ошибки в минуту**
```logql
sum(rate({service="mylittleprice-backend"} |= "level=error" [1m]))
```

**Панель 2: Логи по уровням**
```logql
sum by (level) (rate({service="mylittleprice-backend"} [5m]))
```

**Панель 3: Последние ошибки**
```logql
{service="mylittleprice-backend"} |= "level=error"
```

---

## Полезные ссылки

- [Loki Documentation](https://grafana.com/docs/loki/latest/)
- [LogQL Query Language](https://grafana.com/docs/loki/latest/logql/)
- [Grafana Explore](https://grafana.com/docs/grafana/latest/explore/)
- [Promtail Configuration](https://grafana.com/docs/loki/latest/clients/promtail/configuration/)

---

## Итого

**Для Docker контейнеров (postgres, redis):**
- Используется Promtail
- Автоматический сбор логов
- Метки: `{job="docker_all"}`, `{service="postgres"}`

**Для Backend приложения:**
- Прямая отправка через HTTP API
- Контекстная информация (request_id, user_id)
- Метки: `{service="mylittleprice-backend"}`
- Включается через `LOKI_ENABLED=true`

Оба способа работают одновременно и дополняют друг друга!
