# Changelog

All notable changes to the MyLittlePrice Chat App will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-11-09

### Added
- Initial Flutter project structure with feature-first architecture
- Core configuration and setup (AppConfig, ApiEndpoints, Constants)
- Complete theme system with Material 3 (light/dark modes)
- Data models with Freezed (Product, ChatMessage, User, Session, etc.)
- Storage service (Hive + SharedPreferences)
- Routing with go_router
- Logger utility for centralized logging
- Comprehensive documentation (README.md, FLUTTER_MIGRATION_PLAN.md)

### Changed
- **BREAKING**: Updated to Riverpod 3.0 with code generation support
  - flutter_riverpod: 2.5.1 → 3.0.3
  - Added riverpod_annotation and riverpod_generator
  - Enables `@riverpod` annotation for provider generation

- **BREAKING**: Updated go_router to latest version
  - go_router: 14.6.2 → 17.0.0
  - New routing API with improved type safety

- **BREAKING**: Updated Freezed to version 3.x
  - freezed: 2.5.7 → 3.2.3
  - freezed_annotation: 2.4.4 → 3.1.0

- **BREAKING**: Updated flutter_lints
  - flutter_lints: 5.0.0 → 6.0.0
  - Added comprehensive analysis_options.yaml

- Updated networking dependencies:
  - dio: 5.7.0 → 5.9.0
  - web_socket_channel: 3.0.1 → 3.0.3

- Updated storage dependencies:
  - shared_preferences: 2.3.2 → 2.5.3

- Updated UI/UX dependencies:
  - flutter_svg: 2.0.10 → 2.2.2

- Updated utility dependencies:
  - uuid: 4.5.1 → 4.5.2
  - intl: 0.19.0 → 0.20.2
  - equatable: 2.0.5 → 2.0.7
  - logger: 2.4.0 → 2.6.2

- Updated build tools:
  - build_runner: 2.4.13 → 2.10.1
  - json_serializable: 6.8.0 → 6.11.1

### Removed
- flutter_markdown (discontinued package)
  - Will use flutter_markdown_plus when needed

### Migration Notes

#### Riverpod 3.0
To use the new code generation:
```dart
// Old way (still works)
final myProvider = Provider<String>((ref) => 'Hello');

// New way with code generation
@riverpod
String myProvider(MyProviderRef ref) => 'Hello';
```

Run `dart run build_runner build` after creating providers.

#### go_router 17.0
- Review router configuration for breaking changes
- Some redirect and route builder APIs have changed
- See: https://pub.dev/packages/go_router/changelog

#### Freezed 3.0
- Code generation unchanged
- Some internal API improvements
- Run `dart run build_runner build --delete-conflicting-outputs`

### Dependencies Status

**Current stable versions (as of 2025-11-09):**
- ✅ All dependencies updated to latest stable
- ✅ No deprecated packages
- ✅ Compatible with Flutter SDK ^3.9.2

## [Unreleased]

### Planned
- Network layer (HTTP client with Dio, WebSocket client)
- State management providers (ChatProvider, AuthProvider, SettingsProvider)
- UI components (ChatScreen, MessageWidget, ProductCard)
- Backend API integration
- Offline mode support
- Multi-platform builds (iOS, Android, Web, Desktop)
