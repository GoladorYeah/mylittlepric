import 'package:chat_app/core/config/app_config.dart';

class ApiEndpoints {
  // Base URLs
  static String get baseUrl => AppConfig.apiBaseUrl;
  static String get wsUrl => AppConfig.wsBaseUrl;

  // Full URLs
  static String get chat => '$baseUrl${AppConfig.chatEndpoint}';
  static String get ws => '$wsUrl${AppConfig.wsEndpoint}';
  static String get productDetails => '$baseUrl${AppConfig.productDetailsEndpoint}';
  static String get session => '$baseUrl${AppConfig.sessionEndpoint}';
  static String get login => '$baseUrl${AppConfig.loginEndpoint}';
  static String get logout => '$baseUrl${AppConfig.logoutEndpoint}';
  static String get userPreferences => '$baseUrl${AppConfig.userPreferencesEndpoint}';
  static String get searchHistory => '$baseUrl${AppConfig.searchHistoryEndpoint}';

  // Dynamic endpoints
  static String sessionById(String sessionId) => '$session/$sessionId';
  static String sessionMessages(String sessionId) => '$session/$sessionId/messages';
}
