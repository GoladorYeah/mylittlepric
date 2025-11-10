# MyLittlePrice Chat App

Flutter приложение для AI-powered поиска товаров с чат-интерфейсом.

## 📱 Платформы

- ✅ iOS
- ✅ Android
- ✅ Web
- ✅ macOS
- ✅ Windows
- ✅ Linux

## 🏗️ Архитектура

Проект следует **Feature-First Architecture** с четким разделением ответственности:

```
lib/
├── core/               # Основная функциональность
│   ├── config/        # Конфигурация, константы, роутинг
│   ├── network/       # HTTP client, WebSocket
│   └── storage/       # Локальное хранилище
├── features/          # Функциональные модули
│   ├── auth/          # Авторизация
│   ├── chat/          # Чат интерфейс
│   ├── history/       # История поисков
│   └── settings/      # Настройки
├── shared/            # Общие компоненты
│   ├── models/        # Модели данных
│   ├── widgets/       # Переиспользуемые виджеты
│   ├── providers/     # State management
│   └── utils/         # Утилиты
└── theme/             # Дизайн система
```

## 🔧 Технологический стек

### State Management
- **Riverpod** - Современный, type-safe state management

### Networking
- **Dio** - HTTP client с интерсепторами
- **WebSocket** - Real-time связь с сервером

### Storage
- **Hive** - Локальная NoSQL база данных
- **SharedPreferences** - Простое key-value хранилище

### Code Generation
- **Freezed** - Immutable data classes
- **json_serializable** - JSON сериализация

### UI/UX
- **Material 3** - Современный Material Design
- **go_router** - Декларативная навигация
- **cached_network_image** - Кеширование изображений
- **shimmer** - Loading эффекты

## 🚀 Начало работы

### Требования

- Flutter SDK >= 3.9.2
- Dart SDK >= 3.9.2

### Установка

1. **Установите зависимости:**
```bash
flutter pub get
```

2. **Сгенерируйте код (модели, freezed):**
```bash
dart run build_runner build --delete-conflicting-outputs
```

3. **Запустите приложение:**
```bash
# Development mode
flutter run

# Выбрать конкретную платформу
flutter run -d chrome        # Web
flutter run -d macos          # macOS
flutter run -d windows        # Windows
flutter run -d linux          # Linux
flutter run -d <device_id>    # iOS/Android
```

### Конфигурация

Создайте файл `.env` в корне проекта (опционально):

```env
API_BASE_URL=http://localhost:8080
WS_BASE_URL=ws://localhost:8080
ENABLE_LOGGING=true
```

Или используйте compile-time переменные:

```bash
flutter run --dart-define=API_BASE_URL=http://your-api.com
```

## 📦 Основные команды

```bash
# Установка зависимостей
flutter pub get

# Генерация кода
dart run build_runner build --delete-conflicting-outputs

# Continuous code generation (watch mode)
dart run build_runner watch

# Запуск тестов
flutter test

# Анализ кода
flutter analyze

# Форматирование
dart format lib/

# Очистка build кеша
flutter clean
```

## 🔌 Backend API

Приложение работает с Go backend на порте `8080`:

- **REST API**: `http://localhost:8080/api/*`
- **WebSocket**: `ws://localhost:8080/ws`

См. [backend README](../backend/README.md) для деталей API.

## 🎨 Дизайн система

### Цвета
- **Primary**: Indigo (#6366F1)
- **Secondary**: Purple (#8B5CF6)
- **Accent**: Green (#10B981)

Полная цветовая палитра в `lib/theme/app_colors.dart`

### Типографика
Используются предустановленные стили из Material 3 с кастомизацией в `lib/theme/app_text_styles.dart`

### Spacing & Sizing
Константы определены в `lib/core/config/constants.dart`:
- Spacing: XS(4), S(8), M(16), L(24), XL(32)
- Border Radius: S(4), M(8), L(12), XL(16)
- Icons: S(16), M(24), L(32), XL(48)

## 🧪 Тестирование

```bash
# Все тесты
flutter test

# Конкретный файл
flutter test test/features/chat/chat_test.dart

# С coverage
flutter test --coverage
```

## 📱 Сборка для продакшена

### Android
```bash
flutter build apk --release
flutter build appbundle --release
```

### iOS
```bash
flutter build ios --release
flutter build ipa --release
```

### Web
```bash
flutter build web --release
```

### Desktop
```bash
flutter build macos --release
flutter build windows --release
flutter build linux --release
```

## 🔄 Следующие шаги

После базовой настройки проекта:

1. ✅ Структура проекта создана
2. ✅ Модели данных портированы
3. ✅ Конфигурация настроена
4. ✅ Тема и дизайн-система готовы
5. ⏳ **Реализация сетевых сервисов** (HTTP, WebSocket)
6. ⏳ **Создание провайдеров** (ChatProvider, AuthProvider)
7. ⏳ **Разработка UI компонентов** (ChatScreen, MessageWidget, ProductCard)
8. ⏳ **Интеграция с backend API**

## 📚 Документация

- [Flutter Documentation](https://docs.flutter.dev/)
- [Riverpod Documentation](https://riverpod.dev/)
- [Freezed Documentation](https://pub.dev/packages/freezed)
- [Go Router Documentation](https://pub.dev/packages/go_router)

## 🤝 Связь с Backend

Приложение взаимодействует с Go backend через:

1. **REST API** для одиночных запросов
2. **WebSocket** для real-time чата
3. **Локальное кеширование** через Hive
4. **Синхронизация** между устройствами

## 📝 Лицензия

Proprietary - MyLittlePrice
