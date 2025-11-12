# Monitoring Quick Start Guide

## 🚀 5-Minute Setup

### 1. Start Monitoring Stack

```bash
docker-compose -f docker-compose.monitoring.yml up -d
```

This starts:
- ✅ Prometheus (metrics collection)
- ✅ Grafana (dashboards & visualization)
- ✅ Loki (log aggregation)
- ✅ Promtail (log collection)
- ✅ Alertmanager (alert management)

### 2. Start Backend with Metrics

```bash
cd backend
go run cmd/api/main.go
```

Backend will expose metrics at http://localhost:8080/metrics

### 3. Access Monitoring Interfaces

| Interface | URL | Login |
|-----------|-----|-------|
| **Grafana** | http://localhost:3001 | admin / admin |
| **Prometheus** | http://localhost:9090 | - |
| **Alertmanager** | http://localhost:9093 | - |

### 4. View Metrics in Grafana

1. Open http://localhost:3001
2. Login with `admin` / `admin`
3. Navigate to **Dashboards** → **Browse**
4. Open **WebSocket Monitoring** dashboard

You should see:
- Active WebSocket connections
- Message throughput
- Rate limiting activity
- Processing latency

## 📊 What's Being Monitored?

### WebSocket Metrics
- Active connections
- Connection rate & failures
- Message throughput (sent/received)
- Processing latency (p50, p95, p99)
- Rate limiting violations
- Pub/Sub broadcast activity

### HTTP API Metrics
- Request rate
- Error rate (5xx)
- Response latency
- In-flight requests

### Session Metrics
- Active sessions
- Cache hit/miss rate
- Message persistence failures
- Cleanup job status

### Rate Limiting
- HTTP rate limit violations
- WebSocket rate limit violations
- Redis errors

## 🚨 Critical Alerts

Alerts are pre-configured for:

- **BackendDown** - Backend API unreachable (>1min)
- **HighHTTPErrorRate** - >20% HTTP 5xx errors (>2min)
- **AllWebSocketConnectionsFailing** - >80% WS failures (>2min)
- **MessagePersistenceFailures** - Messages not saving (>5min)

## 🔔 Configure Notifications (Optional)

To receive alerts via Slack/Email:

1. Edit `alertmanager/alertmanager.yml`
2. Uncomment and configure receiver section:

```yaml
receivers:
  - name: 'critical-alerts'
    slack_configs:
      - api_url: 'YOUR_SLACK_WEBHOOK_URL'
        channel: '#alerts-critical'
```

3. Restart Alertmanager:
```bash
docker-compose -f docker-compose.monitoring.yml restart alertmanager
```

## 📈 Useful Prometheus Queries

**Active WebSocket connections:**
```promql
websocket_connections_active
```

**HTTP request rate:**
```promql
sum(rate(http_requests_total[5m]))
```

**Error rate percentage:**
```promql
sum(rate(http_requests_total{status=~"5.."}[5m]))
/
sum(rate(http_requests_total[5m]))
* 100
```

**Session cache hit rate:**
```promql
rate(session_cache_hit_total[5m])
/
(rate(session_cache_hit_total[5m]) + rate(session_cache_miss_total[5m]))
* 100
```

## 🔍 Troubleshooting

**No metrics in Grafana?**
```bash
# Check Prometheus targets
open http://localhost:9090/targets
# All should be "UP"

# Check backend metrics endpoint
curl http://localhost:8080/metrics
```

**Alerts not firing?**
```bash
# Check Prometheus rules
open http://localhost:9090/rules

# Check Alertmanager
open http://localhost:9093/#/alerts
```

## 📚 Full Documentation

For complete documentation see:
- **[MONITORING.md](MONITORING.md)** - Full monitoring guide
- **[BACKEND_IMPROVEMENTS.md](BACKEND_IMPROVEMENTS.md)** - All backend improvements

## 🎯 Next Steps

1. Create custom dashboards for your specific needs
2. Configure notification channels (Slack, Email, PagerDuty)
3. Tune alert thresholds based on your traffic
4. Set up log queries in Loki
5. Add business metrics

---

**Status:** ✅ Production Ready
**Version:** v2.2.0
**Last Updated:** November 12, 2024
