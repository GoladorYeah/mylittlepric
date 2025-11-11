# 📊 Мониторинг MyLittlePrice Backend

Этот проект использует **Prometheus** для сбора метрик и **Grafana** для визуализации.

## 🚀 Быстрый старт

### 1. Запуск мониторинга

```bash
# Запустить Prometheus + Grafana
docker-compose -f docker-compose.monitoring.yml up -d

# Проверить статус
docker-compose -f docker-compose.monitoring.yml ps
```

### 2. Доступ к сервисам

- **Grafana**: http://localhost:3000
  - Логин: `admin`
  - Пароль: `admin` (можно изменить через переменную окружения `GRAFANA_ADMIN_PASSWORD`)

- **Prometheus**: http://localhost:9090

- **Backend Metrics**: http://localhost:8080/metrics

### 3. Health Check эндпоинты

Backend предоставляет 3 health check эндпоинта:

- **Liveness probe**: `GET /health/live`
  - Простая проверка что приложение работает
  - Используется для Kubernetes liveness probe

- **Readiness probe**: `GET /health/ready`
  - Проверка зависимостей (PostgreSQL, Redis)
  - Используется для Kubernetes readiness probe

- **Detailed health**: `GET /health`
  - Детальная информация о здоровье системы
  - Включает статистику Redis connection pool

## 📈 Доступные метрики

### HTTP метрики
- `http_requests_total` - Общее количество HTTP запросов
- `http_request_duration_seconds` - Latency HTTP запросов

### WebSocket метрики
- `websocket_connections_active` - Активные WebSocket соединения
- `websocket_messages_total` - Количество отправленных/полученных сообщений

### Database метрики
- `db_queries_total` - Количество запросов к БД (PostgreSQL, Redis)
- `db_query_duration_seconds` - Latency запросов к БД
- `redis_connection_pool_active` - Активные подключения в Redis pool
- `redis_connection_pool_idle` - Idle подключения в Redis pool

### AI сервисы метрики
- `ai_requests_total` - Количество запросов к AI API (Gemini, SERP)
- `ai_request_duration_seconds` - Latency AI запросов
- `ai_tokens_used_total` - Количество использованных токенов

### Session метрики
- `active_sessions_total` - Количество активных сессий
- `messages_processed_total` - Количество обработанных сообщений

### Cleanup job метрики
- `cleanup_job_runs_total` - Количество запусков cleanup job
- `cleanup_job_records_deleted_total` - Количество удаленных записей
- `cleanup_job_duration_seconds` - Время выполнения cleanup job

### Error метрики
- `errors_total` - Общее количество ошибок (по типу и источнику)

## 🎨 Grafana Dashboards

После первого запуска Grafana автоматически подключится к Prometheus.

### Создание дашбордов

1. Откройте Grafana: http://localhost:3000
2. Войдите (admin/admin)
3. Перейдите в **Dashboards** → **New** → **New Dashboard**
4. Добавьте панели с нужными метриками

### Рекомендуемые дашборды

**1. Overview Dashboard**
- Total requests/second
- Active WebSocket connections
- Active chat sessions
- Error rate (%)
- P50/P95/P99 latency

**2. HTTP Performance Dashboard**
- Request rate by endpoint
- Latency heatmap
- Error rate by endpoint

**3. Database Dashboard**
- PostgreSQL/Redis query rate
- Query latency
- Redis connection pool usage

**4. AI Services Dashboard**
- Gemini API requests/minute
- Token usage (input/output)
- AI errors by service

**5. Session & Messages Dashboard**
- Active sessions over time
- Messages processed/minute

## 🔧 Конфигурация

### Prometheus

Конфигурация находится в `prometheus/prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'mylittleprice_backend'
    static_configs:
      - targets: ['host.docker.internal:8080']
    metrics_path: '/metrics'
    scrape_interval: 10s
```

**Примечание**: `host.docker.internal` используется для доступа к host machine из Docker контейнера.
- На **Linux** может потребоваться изменить на IP адрес host machine
- На **macOS/Windows** работает автоматически

### Grafana

Datasource автоматически настраивается через `grafana/provisioning/datasources/prometheus.yml`.

Чтобы изменить пароль администратора:

```bash
GRAFANA_ADMIN_PASSWORD=your_password docker-compose -f docker-compose.monitoring.yml up -d
```

## 📊 Примеры PromQL запросов

### Средняя latency по эндпоинтам
```promql
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])
```

### Error rate (%)
```promql
(sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))) * 100
```

### P95 latency
```promql
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, endpoint))
```

### Top 5 slowest endpoints
```promql
topk(5, sum(rate(http_request_duration_seconds_sum[5m])) by (endpoint))
```

### Active WebSocket connections
```promql
websocket_connections_active
```

### AI token usage per hour
```promql
sum(rate(ai_tokens_used_total[1h])) by (service, model, type)
```

## 🛑 Остановка мониторинга

```bash
# Остановить контейнеры
docker-compose -f docker-compose.monitoring.yml down

# Остановить контейнеры и удалить данные
docker-compose -f docker-compose.monitoring.yml down -v
```

## 🐛 Troubleshooting

### Prometheus не может подключиться к backend

**Проблема**: `Get "http://host.docker.internal:8080/metrics": dial tcp: lookup host.docker.internal`

**Решение для Linux**:
1. Узнайте IP адрес вашего host:
   ```bash
   ip addr show docker0 | grep -Po 'inet \K[\d.]+'
   ```
2. Замените `host.docker.internal` в `prometheus/prometheus.yml` на полученный IP
3. Перезапустите контейнеры:
   ```bash
   docker-compose -f docker-compose.monitoring.yml restart prometheus
   ```

**Альтернатива**: Используйте network mode `host` для Prometheus.

### Grafana не показывает данные

1. Проверьте что Prometheus собирает метрики:
   - Откройте http://localhost:9090/targets
   - Убедитесь что `mylittleprice_backend` в состоянии **UP**

2. Проверьте datasource в Grafana:
   - **Configuration** → **Data Sources** → **Prometheus**
   - Нажмите **Test** - должно быть "Data source is working"

3. Проверьте что backend работает:
   ```bash
   curl http://localhost:8080/health
   curl http://localhost:8080/metrics
   ```

## 📚 Дополнительные ресурсы

- [Prometheus Documentation](https://prometheus.io/docs/introduction/overview/)
- [Grafana Documentation](https://grafana.com/docs/grafana/latest/)
- [PromQL Basics](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana Dashboard Best Practices](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/best-practices/)

---

**Дата создания**: 11 ноября 2025
**Версия**: 1.0
