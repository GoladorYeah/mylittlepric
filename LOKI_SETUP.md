# Настройка логирования с Grafana Loki

## Проблема
Логи в Grafana Loki были пустые, потому что не было агента для сбора логов.

## Решение
Добавлен **Promtail** - официальный агент для сбора логов и отправки их в Loki.

## Что было сделано:

1. **Создана конфигурация Promtail** (`promtail/promtail-config.yml`):
   - Настроен сбор логов из всех Docker контейнеров
   - Добавлены метки для идентификации контейнеров (имя, ID, сервис, проект)
   - Настроена отправка логов в Loki

2. **Добавлен сервис Promtail** в `docker-compose.monitoring.yml`:
   - Монтирован Docker socket для доступа к логам контейнеров
   - Подключен к сетям `monitoring` и `mylittleprice`
   - Настроена зависимость от Loki

## Запуск:

```bash
# Остановить текущие сервисы мониторинга
docker compose -f docker-compose.monitoring.yml down

# Запустить обновленные сервисы
docker compose -f docker-compose.monitoring.yml up -d

# Проверить логи Promtail
docker logs mylittleprice_promtail

# Проверить статус
docker compose -f docker-compose.monitoring.yml ps
```

## Проверка в Grafana:

1. Откройте Grafana: http://localhost:3001
2. Логин: `admin` / Пароль: из переменной окружения `GRAFANA_ADMIN_PASSWORD` (по умолчанию `admin`)
3. Перейдите в `Explore` (иконка компаса слева)
4. Выберите источник данных `Loki`
5. Используйте запросы LogQL:
   ```logql
   # Все логи
   {job="docker_all"}

   # Логи конкретного сервиса
   {service="postgres"}
   {service="redis"}

   # Логи по имени контейнера
   {container="mylittleprice-postgres"}

   # Поиск по тексту
   {job="docker_all"} |= "error"
   ```

## Архитектура:

```
Docker Контейнеры (postgres, redis, и т.д.)
         ↓ (логи)
    Promtail (собирает логи)
         ↓ (отправка)
      Loki (хранит логи)
         ↓ (запросы)
     Grafana (визуализация)
```

## Устранение неполадок:

Если логи все еще пустые:

1. Проверьте, что Promtail запущен:
   ```bash
   docker ps | grep promtail
   ```

2. Проверьте логи Promtail на ошибки:
   ```bash
   docker logs mylittleprice_promtail
   ```

3. Проверьте, что Loki доступен:
   ```bash
   curl http://localhost:3100/ready
   ```

4. Проверьте метрики Promtail:
   ```bash
   curl http://localhost:9080/metrics
   ```
