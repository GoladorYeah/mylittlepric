# MyLittlePrice - Monitoring & Alerting Guide

**Дата:** 12 ноября 2024
**Статус:** ✅ Production Ready

---

## 📋 Оглавление

- [Обзор](#обзор)
- [Архитектура мониторинга](#архитектура-мониторинга)
- [Быстрый старт](#быстрый-старт)
- [Метрики](#метрики)
- [Дашборды Grafana](#дашборды-grafana)
- [Алерты](#алерты)
- [Troubleshooting](#troubleshooting)

---

## 🎯 Обзор

MyLittlePrice использует полный стек мониторинга для production-ready наблюдаемости:

- **Prometheus** - сбор метрик
- **Grafana** - визуализация и дашборды
- **Loki** - сбор и анализ логов
- **Promtail** - агент для сбора логов из Docker контейнеров
- **Alertmanager** - управление алертами и уведомлениями

### Что мониторится

✅ **WebSocket соединения** - активные connections, rate, ошибки
✅ **Сообщения** - throughput, latency, failures
✅ **Rate Limiting** - violations, Redis errors
✅ **Сессии** - cache hit/miss, cleanup, sync
✅ **HTTP API** - request rate, latency, error rate
✅ **Pub/Sub** - broadcast activity

---

## 🏗️ Архитектура мониторинга

```
┌─────────────────────────────────────────────────┐
│                Backend API (Go)                 │
│  ┌──────────────┐      ┌──────────────┐        │
│  │   /metrics   │      │  Structured  │        │
│  │  (Prometheus)│      │    Logs      │        │
│  └──────┬───────┘      └──────┬───────┘        │
└─────────┼──────────────────────┼────────────────┘
          │                      │
          │                      │
    ┌─────▼──────┐         ┌────▼──────┐
    │ Prometheus │         │  Promtail │
    │  :9090     │         │           │
    └─────┬──────┘         └────┬──────┘
          │                     │
          │                ┌────▼──────┐
          │                │   Loki    │
          │                │  :3100    │
          │                └────┬──────┘
          │                     │
    ┌─────▼─────────────────────▼──────┐
    │         Grafana :3001             │
    │  - Dashboards                     │
    │  - Alerts                         │
    │  - Visualizations                 │
    └───────────────────────────────────┘
          │
    ┌─────▼──────────┐
    │  Alertmanager  │
    │     :9093      │
    └────────────────┘
```

---

## 🚀 Быстрый старт

### 1. Запуск мониторинга

```bash
# Старт всего стека мониторинга
docker-compose -f docker-compose.monitoring.yml up -d

# Проверка статуса
docker-compose -f docker-compose.monitoring.yml ps
```

### 2. Доступ к интерфейсам

| Сервис | URL | Credentials |
|--------|-----|-------------|
| **Grafana** | http://localhost:3001 | admin / admin |
| **Prometheus** | http://localhost:9090 | - |
| **Alertmanager** | http://localhost:9093 | - |
| **Loki** | http://localhost:3100 | - |

### 3. Первая настройка Grafana

1. Открыть http://localhost:3001
2. Войти (admin/admin)
3. Datasources уже преднастроены:
   - Prometheus (default)
   - Loki
4. Дашборды доступны в Dashboards → Browse

---

## 📊 Метрики

### WebSocket Metrics

#### Connections
```promql
# Активные соединения
websocket_connections_active

# Rate соединений
rate(websocket_connections_total[5m])

# Процент ошибок соединений
rate(websocket_connections_failed_total[5m])
/
rate(websocket_connections_total[5m])
```

#### Messages
```promql
# Throughput сообщений (sent)
sum(rate(websocket_messages_sent_total[5m])) by (type)

# Throughput сообщений (received)
sum(rate(websocket_messages_received_total[5m])) by (type)

# Latency (p95)
histogram_quantile(0.95,
  sum(rate(websocket_message_duration_seconds_bucket[5m])) by (le)
)

# Failures
rate(websocket_messages_sent_failed_total[5m])
```

#### Rate Limiting
```promql
# WebSocket rate limit violations
rate(websocket_rate_limit_exceeded_total[5m])

# По уровням (connection vs user)
sum(rate(websocket_rate_limit_exceeded_total[5m])) by (level)
```

#### Pub/Sub
```promql
# Broadcasts sent
rate(websocket_broadcasts_sent_total[5m])

# Broadcasts received
rate(websocket_broadcasts_received_total[5m])
```

### HTTP API Metrics

#### Request Rate
```promql
# Total requests
sum(rate(http_requests_total[5m]))

# By endpoint
sum(rate(http_requests_total[5m])) by (handler)

# By status code
sum(rate(http_requests_total[5m])) by (status)
```

#### Error Rate
```promql
# 5xx error rate
sum(rate(http_requests_total{status=~"5.."}[5m]))
/
sum(rate(http_requests_total[5m]))

# Per endpoint
sum(rate(http_requests_total{status=~"5.."}[5m])) by (handler)
```

#### Latency
```promql
# p50 latency
histogram_quantile(0.50,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (handler, le)
)

# p95 latency
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (handler, le)
)

# p99 latency
histogram_quantile(0.99,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (handler, le)
)
```

### Rate Limiting Metrics

```promql
# HTTP rate limit violations
sum(rate(rate_limit_exceeded_total[5m])) by (endpoint)

# Redis errors in rate limiter
rate_limiter_redis_errors_total
```

### Session Metrics

```promql
# Active sessions
active_sessions_total

# Session creation rate
rate(sessions_created_total[5m])

# Cache hit rate
rate(session_cache_hit_total[5m])
/
(rate(session_cache_hit_total[5m]) + rate(session_cache_miss_total[5m]))

# Message persistence failures
rate(message_persistence_failed_total[5m])
```

---

## 📈 Дашборды Grafana

### 1. WebSocket Monitoring

**Файл:** `grafana/dashboards/websocket-monitoring.json`
**UID:** `websocket-monitoring`

**Панели:**
- Active WebSocket Connections (Gauge)
- Connection Rate (Time Series)
- Message Rate by Type (Time Series)
- Rate Limiting Activity (Time Series)
- Message Processing Time (Time Series)
- Pub/Sub Activity (Time Series)

**Использование:**
```
Dashboards → Browse → WebSocket Monitoring
```

### 2. Создание кастомных дашбордов

#### Пример: HTTP Error Rate Dashboard

1. Grafana → Dashboards → New → New Dashboard
2. Add visualization
3. Query:
```promql
sum(rate(http_requests_total{status=~"5.."}[5m])) by (handler)
```
4. Panel title: "HTTP 5xx Error Rate by Endpoint"
5. Save dashboard

#### Пример: Session Health Dashboard

Query для cache hit rate:
```promql
rate(session_cache_hit_total[5m])
/
(rate(session_cache_hit_total[5m]) + rate(session_cache_miss_total[5m]))
```

---

## 🚨 Алерты

### Конфигурация алертов

Alerting rules находятся в:
- `prometheus/alerts/backend_alerts.yml`
- `prometheus/alerts/websocket_alerts.yml`
- `prometheus/alerts/ratelimit_alerts.yml`
- `prometheus/alerts/session_alerts.yml`

### Критические алерты

#### 1. Backend Down
```yaml
alert: BackendDown
expr: up{job="mylittleprice-backend"} == 0
for: 1m
severity: critical
```
**Действия:**
- Проверить статус backend контейнера
- Проверить логи: `docker logs mylittleprice_backend`
- Проверить health endpoint: `curl http://localhost:8080/health`

#### 2. High HTTP Error Rate
```yaml
alert: CriticalHTTPErrorRate
expr: |
  sum(rate(http_requests_total{status=~"5.."}[5m])) by (handler)
  /
  sum(rate(http_requests_total[5m])) by (handler)
  > 0.20
for: 2m
severity: critical
```
**Действия:**
- Проверить логи backend
- Проверить Grafana Loki для error logs
- Проверить database connectivity
- Проверить Redis connectivity

#### 3. WebSocket Connection Failures
```yaml
alert: AllWebSocketConnectionsFailing
expr: |
  (
    sum(rate(websocket_connections_failed_total[2m]))
    /
    sum(rate(websocket_connections_total[2m]))
  ) > 0.80
for: 2m
severity: critical
```
**Действия:**
- Проверить WebSocket endpoint
- Проверить rate limiting configuration
- Проверить Redis connectivity
- Проверить network issues

#### 4. Message Persistence Failures
```yaml
alert: MessagePersistenceFailures
expr: sum(rate(message_persistence_failed_total[5m])) > 1
for: 5m
severity: critical
```
**Действия:**
- Проверить PostgreSQL connectivity
- Проверить disk space
- Проверить database permissions
- Проверить логи backend

### Warning Alerts

#### 1. High WebSocket Rate Limiting
```yaml
alert: FrequentWebSocketRateLimiting
expr: sum(rate(websocket_rate_limit_exceeded_total[5m])) > 5
for: 5m
severity: warning
```

#### 2. High Session Cache Miss Rate
```yaml
alert: HighSessionCacheMissRate
expr: |
  (
    sum(rate(session_cache_miss_total[5m]))
    /
    sum(rate(session_lookups_total[5m]))
  ) > 0.30
for: 10m
severity: warning
```

### Настройка уведомлений

Редактируйте `alertmanager/alertmanager.yml`:

#### Slack Integration
```yaml
receivers:
  - name: 'critical-alerts'
    slack_configs:
      - api_url: 'YOUR_SLACK_WEBHOOK_URL'
        channel: '#alerts-critical'
        title: '🚨 Critical Alert: {{ .GroupLabels.alertname }}'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
```

#### Email Integration
```yaml
receivers:
  - name: 'critical-alerts'
    email_configs:
      - to: 'alerts@yourdomain.com'
        from: 'prometheus@yourdomain.com'
        smarthost: 'smtp.gmail.com:587'
        auth_username: 'your-email@gmail.com'
        auth_password: 'your-app-password'
```

После изменения конфигурации:
```bash
# Reload Alertmanager config
docker-compose -f docker-compose.monitoring.yml restart alertmanager
```

---

## 🔍 Troubleshooting

### Проблема: Метрики не собираются

**Симптомы:** Grafana показывает "No data"

**Решение:**
```bash
# 1. Проверить, что backend экспортирует метрики
curl http://localhost:8080/metrics

# 2. Проверить статус Prometheus targets
# Открыть http://localhost:9090/targets
# Все targets должны быть "UP"

# 3. Проверить логи Prometheus
docker logs mylittleprice_prometheus

# 4. Проверить prometheus.yml конфигурацию
cat prometheus/prometheus.yml
```

### Проблема: Alertmanager не отправляет уведомления

**Решение:**
```bash
# 1. Проверить статус Alertmanager
curl http://localhost:9093/-/healthy

# 2. Проверить конфигурацию
docker exec mylittleprice_alertmanager amtool check-config /etc/alertmanager/alertmanager.yml

# 3. Проверить логи
docker logs mylittleprice_alertmanager

# 4. Проверить firing alerts
# Открыть http://localhost:9093/#/alerts
```

### Проблема: Grafana не может подключиться к Prometheus

**Решение:**
```bash
# 1. Проверить, что Prometheus доступен из Grafana контейнера
docker exec mylittleprice_grafana wget -O- http://prometheus:9090/api/v1/status/config

# 2. Проверить datasource в Grafana
# Configuration → Data Sources → Prometheus → Test

# 3. Проверить network
docker network inspect monitoring
```

### Проблема: Высокая load на Prometheus

**Решение:**
```yaml
# Увеличить scrape interval в prometheus.yml
global:
  scrape_interval: 30s  # было 15s
  evaluation_interval: 30s  # было 15s

# Уменьшить retention
command:
  - '--storage.tsdb.retention.time=15d'  # было 30d
```

---

## 📚 Дополнительные ресурсы

### Prometheus Queries Cheat Sheet

**Rate vs Increase:**
```promql
# Rate - average per second over interval
rate(http_requests_total[5m])

# Increase - total increase over interval
increase(http_requests_total[5m])
```

**Aggregation:**
```promql
# Sum by label
sum(rate(http_requests_total[5m])) by (handler)

# Average
avg(rate(http_requests_total[5m]))

# Count
count(websocket_connections_active)
```

**Percentiles:**
```promql
# p95 latency
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
)
```

### Полезные ссылки

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
- [PromQL Tutorial](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Alertmanager Configuration](https://prometheus.io/docs/alerting/latest/configuration/)

---

## 🎯 Best Practices

### 1. Naming Conventions

Следуйте Prometheus naming conventions:
- Metric name: `<namespace>_<name>_<unit>_total`
- Labels в snake_case
- Counters заканчиваются на `_total`
- Histograms заканчиваются на `_bucket`, `_sum`, `_count`

### 2. Label Cardinality

❌ **Плохо** (высокая cardinality):
```go
httpRequests.WithLabelValues(userID, sessionID).Inc()
```

✅ **Хорошо** (низкая cardinality):
```go
httpRequests.WithLabelValues(endpoint, method).Inc()
```

### 3. Alert Fatigue

- Не создавайте алерты для каждой метрики
- Группируйте похожие алерты
- Используйте разные severity levels
- Настройте inhibition rules

### 4. Dashboard Organization

- Создавайте dashboards по компонентам (WebSocket, API, Sessions)
- Используйте переменные для фильтрации
- Добавляйте описания к панелям
- Используйте единые цветовые схемы

---

## 📊 Production Checklist

Перед запуском в production убедитесь:

- [ ] Prometheus scraping работает (проверить /targets)
- [ ] Grafana datasources настроены
- [ ] Dashboard WebSocket Monitoring работает
- [ ] Alertmanager получает alerts от Prometheus
- [ ] Уведомления настроены (Slack/Email/PagerDuty)
- [ ] Alert rules протестированы
- [ ] Retention policy настроена
- [ ] Backup strategy для Prometheus data
- [ ] Мониторинг самого мониторинга (meta-monitoring)
- [ ] Документация для oncall team

---

**Автор:** Claude AI Assistant
**Дата создания:** 12 ноября 2024
**Версия:** 1.0
