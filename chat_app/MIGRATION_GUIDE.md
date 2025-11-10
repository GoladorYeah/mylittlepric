# Руководство по миграции зависимостей

Этот документ описывает важные изменения при обновлении зависимостей проекта.

## 🔄 Обзор изменений

Проект был обновлен до последних стабильных версий всех зависимостей по состоянию на ноябрь 2025.

## ⚠️ Breaking Changes

### 1. Riverpod 3.0 (flutter_riverpod 2.5.1 → 3.0.3)

**Что изменилось:**
- Добавлена поддержка code generation через `@riverpod` аннотации
- Новые возможности: automatic retry, offline persistence, mutations
- Улучшенная типобезопасность

**Как мигрировать:**

#### Старый способ (все еще работает):
```dart
final myProvider = Provider<String>((ref) {
  return 'Hello World';
});

final counterProvider = StateNotifierProvider<CounterNotifier, int>((ref) {
  return CounterNotifier();
});
```

#### Новый способ с code generation:
```dart
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'my_provider.g.dart';

@riverpod
String myProvider(MyProviderRef ref) {
  return 'Hello World';
}

@riverpod
class Counter extends _$Counter {
  @override
  int build() => 0;

  void increment() => state++;
}
```

**Действия:**
1. Добавьте `import 'package:riverpod_annotation/riverpod_annotation.dart';`
2. Добавьте `part 'filename.g.dart';`
3. Используйте `@riverpod` аннотацию
4. Запустите: `dart run build_runner build`

**Новые возможности:**

- **Automatic Retry**: провайдеры автоматически повторяют запросы при ошибках
```dart
@Riverpod(keepAlive: true)
Future<Data> myData(MyDataRef ref) async {
  // Автоматически повторится при network errors
  return await fetchData();
}
```

- **ref.mounted**: проверка, жив ли еще провайдер
```dart
@riverpod
Future<void> fetchUser(FetchUserRef ref) async {
  final data = await api.fetch();
  if (!ref.mounted) return; // Безопасная проверка
  ref.state = data;
}
```

**Ссылки:**
- [Riverpod 3.0 Announcement](https://riverpod.dev/docs/introduction/getting_started)
- [Code Generation Guide](https://riverpod.dev/docs/concepts/about_code_generation)

---

### 2. go_router (14.6.2 → 17.0.0)

**Что изменилось:**
- Улучшенная типобезопасность маршрутов
- Новый API для навигации
- Улучшенная поддержка deep linking

**Как мигрировать:**

#### Изменения в redirect:
```dart
// Старый способ
redirect: (context, state) {
  if (!isLoggedIn) return '/login';
  return null;
}

// Новый способ (тот же)
redirect: (context, state) {
  if (!isLoggedIn) return '/login';
  return null;
}
```

#### Type-safe routes (опционально):
```dart
@TypedGoRoute<HomeRoute>(path: '/')
class HomeRoute extends GoRouteData {
  @override
  Widget build(BuildContext context, GoRouterState state) {
    return const HomeScreen();
  }
}

// Использование
HomeRoute().go(context);
```

**Действия:**
1. Проверьте работу существующих роутов
2. Обновите redirect логику если нужно
3. Рассмотрите переход на type-safe routes

**Ссылки:**
- [go_router Changelog](https://pub.dev/packages/go_router/changelog)
- [Type-safe routes](https://pub.dev/packages/go_router#type-safe-routes)

---

### 3. Freezed (2.5.7 → 3.2.3)

**Что изменилось:**
- Улучшенная генерация кода
- Лучшая поддержка generics
- Оптимизация производительности

**Как мигрировать:**

Код остается прежним, но нужно перегенерировать файлы:

```dart
@freezed
class User with _$User {
  const factory User({
    required String id,
    required String name,
    String? email,
  }) = _User;

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
}
```

**Действия:**
1. Запустите: `dart run build_runner clean`
2. Запустите: `dart run build_runner build --delete-conflicting-outputs`

**Ссылки:**
- [Freezed Documentation](https://pub.dev/packages/freezed)

---

### 4. flutter_lints (5.0.0 → 6.0.0)

**Что изменилось:**
- Новые правила линтинга
- Более строгая проверка кода
- Улучшенные рекомендации

**Как мигрировать:**

Созданный `analysis_options.yaml` включает все рекомендуемые правила.

**Частые изменения:**

1. **Trailing commas (требуются)**:
```dart
// До
Widget build(BuildContext context) {
  return Container(
    child: Text('Hello')
  );
}

// После
Widget build(BuildContext context) {
  return Container(
    child: Text('Hello'), // trailing comma
  );
}
```

2. **Const constructors**:
```dart
// До
final widget = Container();

// После
const widget = Container();
```

3. **Prefer single quotes**:
```dart
// До
final text = "Hello";

// После
final text = 'Hello';
```

**Действия:**
1. Запустите: `flutter analyze`
2. Исправьте предупреждения
3. Или запустите: `dart fix --apply` для автоматических исправлений

---

## ✅ Минорные обновления (без breaking changes)

### Networking
- **dio**: 5.7.0 → 5.9.0
  - Улучшения производительности
  - Исправления багов

- **web_socket_channel**: 3.0.1 → 3.0.3
  - Стабильность соединений

### Storage
- **shared_preferences**: 2.3.2 → 2.5.3
  - Улучшенная поддержка платформ

### UI/UX
- **flutter_svg**: 2.0.10 → 2.2.2
  - Улучшенный рендеринг SVG
  - Поддержка новых features

### Utils
- **uuid**: 4.5.1 → 4.5.2
- **intl**: 0.19.0 → 0.20.2
  - Обновленные локализации
- **equatable**: 2.0.5 → 2.0.7
- **logger**: 2.4.0 → 2.6.2
  - Улучшенное форматирование

### Build Tools
- **build_runner**: 2.4.13 → 2.10.1
  - Быстрее генерация
- **json_serializable**: 6.8.0 → 6.11.1

---

## 📦 Пошаговая миграция

### Шаг 1: Очистка
```bash
flutter clean
flutter pub get
```

### Шаг 2: Генерация кода
```bash
dart run build_runner clean
dart run build_runner build --delete-conflicting-outputs
```

### Шаг 3: Анализ кода
```bash
flutter analyze
```

### Шаг 4: Автоматические исправления
```bash
dart fix --apply
```

### Шаг 5: Тестирование
```bash
flutter test
```

### Шаг 6: Запуск
```bash
flutter run
```

---

## 🔧 Решение проблем

### Ошибки компиляции после обновления

1. **Очистите build cache**:
```bash
flutter clean
rm -rf .dart_tool/
flutter pub get
```

2. **Перегенерируйте код**:
```bash
dart run build_runner build --delete-conflicting-outputs
```

3. **Проверьте imports**:
Убедитесь что все импорты используют `package:` вместо относительных путей.

### Конфликты зависимостей

Если `flutter pub get` выдает ошибки:

1. Удалите `pubspec.lock`
2. Запустите `flutter pub get`
3. Проверьте совместимость версий

### Проблемы с go_router

Если роутинг не работает:
1. Проверьте initialLocation
2. Убедитесь что все routes имеют уникальные paths
3. Проверьте redirect логику

---

## 📚 Дополнительные ресурсы

- [Flutter Breaking Changes](https://docs.flutter.dev/release/breaking-changes)
- [Riverpod Migration Guide](https://riverpod.dev/docs/introduction/getting_started)
- [go_router Migration](https://pub.dev/packages/go_router/changelog)
- [Freezed Documentation](https://pub.dev/packages/freezed)

---

## ✨ Что дальше?

После успешной миграции:

1. ✅ Обновите CI/CD пайплайны
2. ✅ Обновите документацию команды
3. ✅ Начните использовать новые возможности Riverpod 3.0
4. ✅ Рассмотрите переход на type-safe routes в go_router
5. ✅ Примените dart fix для улучшения кода
