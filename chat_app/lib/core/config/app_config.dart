class AppConfig {
  // API Configuration
  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );

  static const String wsBaseUrl = String.fromEnvironment(
    'WS_BASE_URL',
    defaultValue: 'ws://localhost:8080',
  );

  // API Endpoints
  static const String chatEndpoint = '/api/chat';
  static const String wsEndpoint = '/ws';
  static const String productDetailsEndpoint = '/api/product-details';
  static const String sessionEndpoint = '/api/session';
  static const String loginEndpoint = '/api/auth/login';
  static const String logoutEndpoint = '/api/auth/logout';
  static const String userPreferencesEndpoint = '/api/user/preferences';
  static const String searchHistoryEndpoint = '/api/search-history';

  // Session Configuration
  static const int maxMessagesPerSession = 8;
  static const int maxSearchesPerSession = 3;
  static const Duration sessionTimeout = Duration(hours: 24);

  // Cache Configuration
  static const Duration cacheTimeout = Duration(hours: 24);
  static const int maxCachedMessages = 100;

  // UI Configuration
  static const int messageAnimationDuration = 300; // ms
  static const int typingIndicatorDelay = 500; // ms
  static const int maxQuickReplies = 3;

  // Storage Keys
  static const String sessionIdKey = 'chat_session_id';
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String userKey = 'user';
  static const String countryKey = 'country';
  static const String languageKey = 'language';
  static const String currencyKey = 'currency';
  static const String themeKey = 'theme_mode';

  // Hive Box Names
  static const String messagesBox = 'messages';
  static const String sessionBox = 'sessions';
  static const String preferencesBox = 'preferences';
  static const String cacheBox = 'cache';

  // Default Values
  static const String defaultCountry = 'US';
  static const String defaultLanguage = 'en';
  static const String defaultCurrency = 'USD';

  // Network Configuration
  static const Duration connectionTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);
  static const int maxRetries = 3;

  // WebSocket Configuration
  static const Duration wsReconnectDelay = Duration(seconds: 2);
  static const int wsMaxReconnectAttempts = 5;
  static const Duration wsPingInterval = Duration(seconds: 30);

  // Logging
  static const bool enableLogging = bool.fromEnvironment(
    'ENABLE_LOGGING',
    defaultValue: true,
  );
}
