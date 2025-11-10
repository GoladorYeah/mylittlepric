# ✅ State Management (Этап 3) - Готово!

## Что было реализовано

Созданы все необходимые Riverpod провайдеры для управления состоянием приложения:

### 1. 🔐 AuthProvider
**Файлы:**
- `lib/features/auth/providers/auth_state.dart`
- `lib/features/auth/providers/auth_provider.dart`

**Функционал:**
- ✅ Автоматическое восстановление сессии из storage при запуске
- ✅ Login/Logout с OAuth интеграцией
- ✅ Автоматическое обновление токенов (access/refresh)
- ✅ Обновление пользовательских preferences (country, language, currency)
- ✅ Удаление аккаунта

**Использование:**
```dart
// В виджете
final authState = ref.watch(authProvider);
final isAuthenticated = ref.watch(isAuthenticatedProvider);
final currentUser = ref.watch(currentUserProvider);

// Действия
await ref.read(authProvider.notifier).login(
  provider: 'google',
  token: 'oauth_token',
);

await ref.read(authProvider.notifier).logout();

await ref.read(authProvider.notifier).updatePreferences(
  country: 'US',
  language: 'en',
  currency: 'USD',
);
```

---

### 2. 💬 ChatProvider
**Файлы:**
- `lib/features/chat/providers/chat_state.dart`
- `lib/features/chat/providers/chat_provider.dart`

**Функционал:**
- ✅ WebSocket соединение с автоматическим переподключением
- ✅ REST API fallback если WebSocket недоступен
- ✅ Управление сообщениями с автоматическим сохранением в Hive
- ✅ Quick replies поддержка
- ✅ Typing indicator
- ✅ Восстановление сессии и сообщений из storage

**Использование:**
```dart
// В виджете
final chatState = ref.watch(chatProvider);
final messages = ref.watch(chatMessagesProvider);
final quickReplies = ref.watch(chatQuickRepliesProvider);
final isTyping = ref.watch(chatIsTypingProvider);
final isConnected = ref.watch(chatIsConnectedProvider);

// Отправка сообщений
await ref.read(chatProvider.notifier).sendMessageViaWebSocket('Hello!');
// или через REST API
await ref.read(chatProvider.notifier).sendMessageViaRest('Hello!');

// Quick reply
await ref.read(chatProvider.notifier).sendQuickReply('Show me laptops');

// Очистка чата
await ref.read(chatProvider.notifier).clearChat();
```

---

### 3. ⚙️ SettingsProvider
**Файлы:**
- `lib/features/settings/providers/settings_state.dart`
- `lib/features/settings/providers/settings_provider.dart`

**Функционал:**
- ✅ Theme mode управление (light/dark/system)
- ✅ Выбор страны (10 стран)
- ✅ Выбор языка (8 языков)
- ✅ Выбор валюты (6 валют)
- ✅ Notifications и sound переключатели
- ✅ Автоматическая синхронизация с backend через AuthProvider
- ✅ Сохранение в SharedPreferences

**Использование:**
```dart
// В виджете
final settings = ref.watch(settingsProvider);
final themeMode = ref.watch(themeModeProvider);
final country = ref.watch(countryProvider);

// Изменение темы
await ref.read(settingsProvider.notifier).setThemeMode(ThemeMode.dark);
await ref.read(settingsProvider.notifier).toggleTheme();

// Изменение настроек (автоматически синхронизируется с backend)
await ref.read(settingsProvider.notifier).setCountry('UA');
await ref.read(settingsProvider.notifier).setLanguage('uk');
await ref.read(settingsProvider.notifier).setCurrency('UAH');

// Или все сразу
await ref.read(settingsProvider.notifier).setPreferences(
  country: 'UA',
  language: 'uk',
  currency: 'UAH',
);

// Доступные списки
print(availableCountries); // List<Country>
print(availableLanguages); // List<Language>
print(availableCurrencies); // List<Currency>
```

---

### 4. 📜 SessionProvider (History)
**Файлы:**
- `lib/features/history/providers/session_state.dart`
- `lib/features/history/providers/session_provider.dart`

**Функционал:**
- ✅ Загрузка истории поисков с сервера
- ✅ Пагинация (load more)
- ✅ Сохранение новых поисков
- ✅ Удаление поисков
- ✅ Поиск по истории
- ✅ Фильтрация по категориям
- ✅ Получение уникальных категорий

**Использование:**
```dart
// В виджете
final sessionState = ref.watch(sessionProvider);
final searches = ref.watch(searchHistoryProvider);
final isLoading = ref.watch(searchHistoryLoadingProvider);
final categories = ref.watch(uniqueCategoriesProvider);

// Загрузка истории
await ref.read(sessionProvider.notifier).loadSearches();
await ref.read(sessionProvider.notifier).loadMoreSearches(); // пагинация
await ref.read(sessionProvider.notifier).refreshSearches(); // обновить

// Сохранение поиска
await ref.read(sessionProvider.notifier).saveSearch(
  query: 'gaming laptop',
  category: 'Electronics',
  productIds: ['id1', 'id2'],
);

// Удаление
await ref.read(sessionProvider.notifier).deleteSearch('search_id');

// Поиск в истории
final filtered = ref.read(sessionProvider.notifier).searchInHistory('laptop');

// Фильтрация по категории
final electronics = ref.read(sessionProvider.notifier).getSearchesByCategory('Electronics');
```

---

## 🚀 Следующие шаги

### 1. Генерация Freezed кода

**ВАЖНО:** Перед запуском приложения нужно сгенерировать Freezed файлы:

```bash
cd chat_app
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
```

Эта команда создаст все `*.freezed.dart` файлы для State классов.

### 2. Обновить main.dart

В `lib/main.dart` нужно использовать `themeModeProvider`:

```dart
@override
Widget build(BuildContext context, WidgetRef ref) {
  final themeMode = ref.watch(themeModeProvider);

  return MaterialApp.router(
    title: 'MyLittlePrice',
    debugShowCheckedModeBanner: false,
    theme: AppTheme.lightTheme,
    darkTheme: AppTheme.darkTheme,
    themeMode: themeMode, // вместо ThemeMode.system
    routerConfig: AppRouter.router,
  );
}
```

### 3. Исправить импорты в services

Нужно добавить Provider для StorageService в `lib/core/storage/storage_service.dart`:

```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';

// В конце файла добавить:
final storageServiceProvider = Provider<StorageService>((ref) {
  // Это статический сервис, так что просто возвращаем экземпляр
  return StorageService();
});
```

### 4. Обновить WebSocketClient

В `lib/core/network/websocket_client.dart` нужно исправить метод для работы с StorageService:

```dart
// Было:
final sessionId = await _storageService.getString(AppConfig.sessionIdKey);

// Должно быть:
final prefs = StorageService.prefs;
final sessionId = prefs.getString(AppConfig.sessionIdKey);
```

### 5. Добавить недостающие константы в AppConfig

В `lib/core/config/app_config.dart` добавить:

```dart
// Storage keys
static const String accessTokenKey = 'access_token';
static const String refreshTokenKey = 'refresh_token';
static const String sessionIdKey = 'session_id';
```

---

## 📝 Архитектура State Management

```
┌─────────────────────────────────────────────┐
│           ProviderScope (main.dart)         │
└───────────────────┬─────────────────────────┘
                    │
        ┌───────────┴───────────┐
        │                       │
┌───────▼────────┐      ┌──────▼───────┐
│  AuthProvider  │◄─────┤SettingsProvider│
│                │      │                │
│ - login()      │      │ - setCountry() │
│ - logout()     │      │ - setLanguage()│
│ - refresh()    │      │ - setTheme()   │
└────────┬───────┘      └────────────────┘
         │
         │ uses
         ▼
┌────────────────┐      ┌────────────────┐
│  ChatProvider  │      │SessionProvider │
│                │      │                │
│ - sendMessage()│      │ - loadSearches()│
│ - clearChat()  │      │ - saveSearch() │
└────────────────┘      └────────────────┘
```

## 🔗 Интеграция с Backend

Все провайдеры интегрированы с API сервисами:
- `AuthProvider` → `AuthApiService`
- `ChatProvider` → `ChatApiService` + `WebSocketClient`
- `SettingsProvider` → `AuthApiService` (для синхронизации preferences)
- `SessionProvider` → `SessionApiService`

## 📦 Зависимости

Убедитесь что в `pubspec.yaml` есть:
```yaml
dependencies:
  flutter_riverpod: ^2.4.9
  freezed_annotation: ^2.4.1
  uuid: ^4.2.2

dev_dependencies:
  freezed: ^2.4.5
  build_runner: ^2.4.7
```

---

## ✨ Готово к использованию!

После выполнения шагов выше, все провайдеры готовы к использованию в UI компонентах.

Следующий этап: **Этап 4 - UI Компоненты**
