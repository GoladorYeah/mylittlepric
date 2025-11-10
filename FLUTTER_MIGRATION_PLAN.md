# План миграции на Flutter

## ✅ Этап 1: Базовая структура (ЗАВЕРШЕН)

### Что сделано:

1. **Структура проекта** - 32 директории
   - ✅ `lib/core/` - конфигурация, сеть, хранилище
   - ✅ `lib/features/` - auth, chat, history, settings
   - ✅ `lib/shared/` - модели, виджеты, провайдеры, утилиты
   - ✅ `lib/theme/` - дизайн-система

2. **Модели данных** - 8 файлов (всего 18 Dart файлов)
   - ✅ `Product` - товар с полями из SerpAPI
   - ✅ `ChatMessage` - сообщение чата
   - ✅ `User` - пользователь
   - ✅ `SessionResponse` - ответ сессии
   - ✅ `ChatResponse` - ответ чата
   - ✅ `SavedSearch` - сохраненный поиск
   - ✅ `ProductDetails` - детали товара

3. **Конфигурация**
   - ✅ `AppConfig` - константы конфигурации
   - ✅ `ApiEndpoints` - URL endpoints
   - ✅ `Constants` - константы приложения
   - ✅ `Router` - навигация с go_router

4. **Дизайн-система**
   - ✅ `AppColors` - цветовая палитра (light/dark)
   - ✅ `AppTextStyles` - типографика
   - ✅ `AppTheme` - Material 3 темы

5. **Инфраструктура**
   - ✅ `StorageService` - Hive + SharedPreferences
   - ✅ `AppLogger` - логирование
   - ✅ `main.dart` - entry point с Riverpod

6. **Зависимости** (pubspec.yaml)
   - ✅ Riverpod - state management
   - ✅ Dio - HTTP client
   - ✅ WebSocket - real-time
   - ✅ Hive - локальная БД
   - ✅ Freezed - code generation
   - ✅ go_router - навигация

## ✅ Этап 2: Сетевой слой (ЗАВЕРШЕН)

### Что сделано:

1. **HTTP Client**
   - ✅ `lib/core/network/dio_client.dart` - настройка Dio
   - ✅ `lib/core/network/api_exception.dart` - типизированные исключения
   - ✅ `lib/core/network/interceptors/auth_interceptor.dart` - авто-добавление токенов
   - ✅ `lib/core/network/interceptors/logging_interceptor.dart` - логирование запросов
   - ✅ `lib/core/network/interceptors/retry_interceptor.dart` - автоматические повторы
   - ✅ Обработка ошибок (timeout, network, server, unauthorized)
   - ✅ Exponential backoff для retry

2. **WebSocket Client**
   - ✅ `lib/core/network/websocket_client.dart` - полнофункциональный клиент
   - ✅ Автоматическое переподключение с exponential backoff
   - ✅ Heartbeat (ping/pong) каждые 30 секунд
   - ✅ Обработка событий через Streams
   - ✅ Управление состоянием подключения
   - ✅ Graceful disconnect

3. **API Сервисы**
   - ✅ `ChatApiService` - отправка сообщений, quick replies, детали продуктов
   - ✅ `SessionApiService` - CRUD сессий, история поисков
   - ✅ `AuthApiService` - login/logout, refresh токенов, preferences
   - ✅ `ProductApiService` - детали товаров, поиск, tracking
   - ✅ Все сервисы с Riverpod providers

## ✅ Этап 3: State Management (Riverpod Providers) (ЗАВЕРШЕН)

### Что сделано:

1. **Auth Providers**
   - ✅ `lib/features/auth/providers/auth_state.dart` - состояние аутентификации
   - ✅ `lib/features/auth/providers/auth_provider.dart` - управление авторизацией
   - ✅ Автоматическое восстановление сессии из storage
   - ✅ Обновление токенов (access/refresh)
   - ✅ Login/Logout с интеграцией OAuth
   - ✅ Обновление preferences пользователя
   - ✅ Helper providers: `isAuthenticatedProvider`, `currentUserProvider`, `authLoadingProvider`

2. **Chat Providers**
   - ✅ `lib/features/chat/providers/chat_state.dart` - состояние чата
   - ✅ `lib/features/chat/providers/chat_provider.dart` - управление чатом
   - ✅ WebSocket интеграция с автоматическим переподключением
   - ✅ REST API fallback при отсутствии WebSocket
   - ✅ Управление сообщениями с сохранением в Hive
   - ✅ Quick replies поддержка
   - ✅ Typing indicator
   - ✅ Helper providers: `chatMessagesProvider`, `chatQuickRepliesProvider`, `chatIsTypingProvider`, `chatIsConnectedProvider`, `chatIsSendingProvider`

3. **Settings Providers**
   - ✅ `lib/features/settings/providers/settings_state.dart` - состояние настроек
   - ✅ `lib/features/settings/providers/settings_provider.dart` - управление настройками
   - ✅ Theme mode (light/dark/system)
   - ✅ Страна, язык, валюта с списками доступных значений
   - ✅ Notifications и sound переключатели
   - ✅ Синхронизация с backend через AuthProvider
   - ✅ Helper providers: `themeModeProvider`, `countryProvider`, `languageProvider`, `currencyProvider`

4. **Session Providers**
   - ✅ `lib/features/history/providers/session_state.dart` - состояние истории
   - ✅ `lib/features/history/providers/session_provider.dart` - управление историей
   - ✅ Загрузка истории поисков с пагинацией
   - ✅ Сохранение и удаление поисков
   - ✅ Поиск по истории и фильтрация по категориям
   - ✅ Helper providers: `searchHistoryProvider`, `searchHistoryLoadingProvider`, `uniqueCategoriesProvider`

## ⏳ Этап 4: UI Компоненты

### Задачи:

1. **Chat Feature**
   - [ ] `ChatScreen` - основной экран чата
   - [ ] `MessageWidget` - виджет сообщения
   - [ ] `ChatInput` - поле ввода
   - [ ] `QuickReplies` - быстрые ответы
   - [ ] `ProductCard` - карточка товара
   - [ ] `TypingIndicator` - индикатор печати

2. **History Feature**
   - [ ] `HistoryScreen` - список истории
   - [ ] `SessionCard` - карточка сессии
   - [ ] Фильтры и поиск

3. **Settings Feature**
   - [ ] `SettingsScreen` - настройки
   - [ ] Выбор страны/языка/валюты
   - [ ] Переключение темы

4. **Auth Feature**
   - [ ] `LoginScreen` - экран входа
   - [ ] OAuth интеграция

5. **Shared Widgets**
   - [ ] `LoadingShimmer` - loading эффект
   - [ ] `ErrorWidget` - ошибки
   - [ ] `EmptyState` - пустое состояние

## ⏳ Этап 5: Интеграция

### Задачи:

1. **WebSocket Integration**
   - [ ] Подключение к `/ws`
   - [ ] Обработка сообщений
   - [ ] Синхронизация состояния

2. **REST API Integration**
   - [ ] CRUD операции
   - [ ] Кеширование
   - [ ] Offline mode

3. **Storage Integration**
   - [ ] Сохранение сессий в Hive
   - [ ] Кеш продуктов
   - [ ] Preferences

## ⏳ Этап 6: Тестирование и Оптимизация

### Задачи:

1. **Unit Tests**
   - [ ] Providers
   - [ ] Services
   - [ ] Models

2. **Widget Tests**
   - [ ] ChatScreen
   - [ ] MessageWidget
   - [ ] ProductCard

3. **Integration Tests**
   - [ ] E2E сценарии
   - [ ] WebSocket флоу

4. **Оптимизация**
   - [ ] Performance profiling
   - [ ] Memory leaks
   - [ ] Build size optimization

## 📊 Прогресс

```
Этап 1: ████████████████████ 100% (ЗАВЕРШЕН)
Этап 2: ████████████████████ 100% (ЗАВЕРШЕН)
Этап 3: ████████████████████ 100% (ЗАВЕРШЕН)
Этап 4: ░░░░░░░░░░░░░░░░░░░░   0%
Этап 5: ░░░░░░░░░░░░░░░░░░░░   0%
Этап 6: ░░░░░░░░░░░░░░░░░░░░   0%

Общий прогресс: 50%
```

## 🎯 Следующий шаг

**Начать Этап 4: UI Компоненты**

### Необходимо перед началом:

```bash
cd chat_app
flutter pub run build_runner build --delete-conflicting-outputs
```

Эта команда сгенерирует Freezed код для всех State классов (`*.freezed.dart` файлы).

### План:

1. Создать UI компоненты для Chat Feature
2. Реализовать History Feature UI
3. Создать Settings Screen
4. Реализовать Auth Screen с OAuth

## 🔄 Взаимодействие Flutter ↔ Next.js

### Flutter App (Основной функционал)
- Чат интерфейс
- История поисков
- Настройки пользователя
- Работа с продуктами
- Offline режим

### Next.js Site (Маркетинг)
- Landing page
- Политики (Privacy, Terms, Cookie, Advertising)
- OAuth redirect handler
- SEO оптимизация

### Общий Backend (Go)
```
┌─────────────────┐
│   Go Backend    │
│  (Port 8080)    │
│  REST + WS      │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼──────┐
│Flutter│ │Next.js  │
│ App   │ │Marketing│
└───────┘ └─────────┘
```

## 📝 Примечания

- **Freezed генерация**: Запустите `flutter pub run build_runner build` после создания моделей
- **API URL**: По умолчанию `http://localhost:8080`, настраивается через env
- **WebSocket**: Автоматическое переподключение при потере связи
- **Offline**: Hive кеширует сообщения и сессии
- **Multi-platform**: Поддержка iOS, Android, Web, Desktop

## 🚀 Быстрый старт

```bash
cd chat_app

# 1. Установить зависимости
flutter pub get

# 2. Сгенерировать код
flutter pub run build_runner build --delete-conflicting-outputs

# 3. Запустить (выберите платформу)
flutter run -d chrome      # Web
flutter run -d macos        # macOS
flutter run                 # Default device
```

## 📚 Ресурсы

- [CLAUDE.md](CLAUDE.md) - Контекст проекта
- [chat_app/README.md](chat_app/README.md) - Flutter документация
- [ARCHITECTURE.md](ARCHITECTURE.md) - Общая архитектура
