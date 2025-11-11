# MyLittlePrice Backend Monitoring

Этот документ описывает настройку и использование мониторинга для MyLittlePrice Backend с помощью Prometheus и Grafana.

## 📋 Содержание

- [Архитектура мониторинга](#архитектура-мониторинга)
- [Быстрый старт](#быстрый-старт)
- [Метрики](#метрики)
- [Grafana дашборды](#grafana-дашборды)
- [Troubleshooting](#troubleshooting)

## 🏗️ Архитектура мониторинга

```
MyLittlePrice Backend (:8080)
        ↓ /metrics endpoint
    Prometheus (:9090)
        ↓ scrape metrics
    Grafana (:3001)
        ↓ visualize
    Dashboards
```

## 🚀 Быстрый старт

### 1. Запуск мониторинга

```bash
# Запустить Prometheus и Grafana
docker-compose -f docker-compose.monitoring.yml up -d

# Проверить статус
docker-compose -f docker-compose.monitoring.yml ps
```

### 2. Доступ к сервисам

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001
  - Логин: `admin`
  - Пароль: `admin` (можно изменить через `GRAFANA_ADMIN_PASSWORD` в .env)

### 3. Проверка метрик

Проверьте, что backend экспортирует метрики:

```bash
curl http://localhost:8080/metrics
```

### 4. Настройка Grafana

1. Откройте Grafana (http://localhost:3001)
2. Войдите (admin/admin)
3. Datasource Prometheus уже настроен автоматически
4. Импортируйте дашборды из `grafana/dashboards/`

## 📊 Метрики

Backend экспортирует следующие категории метрик:

### HTTP Метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `http_requests_total` | Counter | Общее количество HTTP запросов (method, endpoint, status) |
| `http_request_duration_seconds` | Histogram | Latency HTTP запросов (method, endpoint) |

**Примеры запросов:**

```promql
# Request rate по endpoint
rate(http_requests_total[5m])

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Error rate
rate(http_requests_total{status=~"5.."}[5m])
```

### WebSocket Метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `websocket_connections_active` | Gauge | Активные WebSocket соединения |
| `websocket_messages_total` | Counter | Количество WebSocket сообщений (type: sent/received) |

### Database Метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `db_queries_total` | Counter | Количество DB запросов (database, operation) |
| `db_query_duration_seconds` | Histogram | Latency DB запросов |
| `redis_connection_pool_active` | Gauge | Активные Redis connections |
| `redis_connection_pool_idle` | Gauge | Idle Redis connections |

### AI Сервисы Метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `ai_requests_total` | Counter | AI API запросы (service, model, status) |
| `ai_request_duration_seconds` | Histogram | Latency AI запросов |
| `ai_tokens_used_total` | Counter | Использованные AI токены (service, model, type: input/output) |

**Примеры:**

```promql
# Gemini API usage rate
rate(ai_requests_total{service="gemini"}[5m])

# Token usage per minute
rate(ai_tokens_used_total{service="gemini",type="input"}[1m]) * 60

# AI error rate
rate(ai_requests_total{status="error"}[5m])
```

### Session & Message Метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `active_sessions_total` | Gauge | Активные chat сессии |
| `messages_processed_total` | Counter | Обработанные сообщения (status) |
| `message_processing_duration_seconds` | Histogram | Время обработки сообщения |

### Cleanup Job Метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `cleanup_job_runs_total` | Counter | Количество запусков cleanup |
| `cleanup_job_records_deleted_total` | Counter | Удаленные записи |
| `cleanup_job_duration_seconds` | Histogram | Длительность cleanup |

### Error Метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `errors_total` | Counter | Ошибки по типам (type, source) |

## 📈 Grafana Дашборды

### Создание дашбордов

Рекомендуемые дашборды (создайте вручную или импортируйте):

#### 1. Overview Dashboard
- Total requests/second
- Active WebSocket connections  
- Active sessions
- Error rate (%)
- P95 latency

#### 2. HTTP Performance
- Request rate by endpoint
- Latency distribution
- Status code breakdown

#### 3. AI Services
- Gemini API requests/minute
- Token usage (input/output)
- API latency
- Error rate

#### 4. Database
- Query rate (PostgreSQL/Redis)
- Redis pool utilization
- Query latency

#### 5. Sessions & Messages
- Active sessions over time
- Messages processed/minute
- Processing latency

## 🔧 Troubleshooting

### Prometheus не видит метрики

1. Проверьте, что backend запущен:
   ```bash
   curl http://localhost:8080/health
   ```

2. Проверьте endpoint метрик:
   ```bash
   curl http://localhost:8080/metrics
   ```

3. Проверьте конфигурацию Prometheus:
   ```bash
   docker logs mylittleprice_prometheus
   ```

4. Проверьте targets в Prometheus UI:
   - Откройте http://localhost:9090/targets
   - Убедитесь, что `mylittleprice_backend` в состоянии UP

### Grafana не показывает данные

1. Проверьте datasource:
   - Configuration → Data Sources → Prometheus
   - Test connection

2. Проверьте, что Prometheus собирает метрики:
   - Откройте Prometheus UI (http://localhost:9090)
   - Попробуйте query: `up{job="mylittleprice_backend"}`

3. Проверьте временной диапазон в Grafana (правый верхний угол)

### Backend не экспортирует метрики

1. Проверьте, что код скомпилирован с метриками:
   ```bash
   cd backend
   go build -o bin/api ./cmd/api
   ```

2. Проверьте логи backend на ошибки:
   ```bash
   docker logs mylittleprice_backend
   ```

## 📝 Дополнительные ресурсы

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Prometheus Query Language (PromQL)](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Go Prometheus Client](https://github.com/prometheus/client_golang)

## 🔐 Production Considerations

При развертывании в production:

1. **Защита доступа**:
   - Настройте аутентификацию для Prometheus
   - Используйте сложный пароль для Grafana admin
   - Настройте HTTPS

2. **Retention**:
   - Настройте retention period для Prometheus (по умолчанию 15 дней)
   - Используйте remote storage для long-term хранения

3. **Alerting**:
   - Настройте alerting rules в Prometheus
   - Интегрируйте с Alertmanager для уведомлений

4. **Backup**:
   - Регулярно делайте backup Grafana дашбордов
   - Бэкапьте Prometheus данные

