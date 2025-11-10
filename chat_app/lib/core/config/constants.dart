class AppConstants {
  // App Info
  static const String appName = 'MyLittlePrice';
  static const String appVersion = '0.1.0';

  // Routes
  static const String splashRoute = '/';
  static const String chatRoute = '/chat';
  static const String historyRoute = '/history';
  static const String settingsRoute = '/settings';
  static const String loginRoute = '/login';

  // Message Types
  static const String messageTypeDialogue = 'dialogue';
  static const String messageTypeSearch = 'search';

  // Response Types
  static const String responseTypeDialogue = 'dialogue';
  static const String responseTypeSearch = 'search';

  // Search Types
  static const String searchTypeExact = 'exact_model';
  static const String searchTypeParameters = 'by_parameters';
  static const String searchTypeCategory = 'by_category';

  // Message Roles
  static const String roleUser = 'user';
  static const String roleAssistant = 'assistant';

  // Session Status
  static const String sessionStatusIdle = 'idle';
  static const String sessionStatusSearching = 'searching';
  static const String sessionStatusCompleted = 'completed';

  // UI Strings
  static const String defaultErrorMessage = 'Something went wrong. Please try again.';
  static const String networkErrorMessage = 'No internet connection. Please check your network.';
  static const String timeoutErrorMessage = 'Request timed out. Please try again.';
  static const String unauthorizedErrorMessage = 'Please log in to continue.';

  // Animation Durations (milliseconds)
  static const int shortAnimationDuration = 200;
  static const int mediumAnimationDuration = 300;
  static const int longAnimationDuration = 500;

  // Spacing
  static const double spacingXs = 4.0;
  static const double spacingS = 8.0;
  static const double spacingM = 16.0;
  static const double spacingL = 24.0;
  static const double spacingXl = 32.0;

  // Border Radius
  static const double radiusS = 4.0;
  static const double radiusM = 8.0;
  static const double radiusL = 12.0;
  static const double radiusXl = 16.0;

  // Icon Sizes
  static const double iconSizeS = 16.0;
  static const double iconSizeM = 24.0;
  static const double iconSizeL = 32.0;
  static const double iconSizeXl = 48.0;

  // Font Sizes
  static const double fontSizeXs = 12.0;
  static const double fontSizeS = 14.0;
  static const double fontSizeM = 16.0;
  static const double fontSizeL = 20.0;
  static const double fontSizeXl = 24.0;
  static const double fontSizeXxl = 32.0;

  // Available Countries (country code -> country name)
  static const Map<String, String> availableCountries = {
    'US': 'United States',
    'GB': 'United Kingdom',
    'DE': 'Germany',
    'FR': 'France',
    'ES': 'Spain',
    'IT': 'Italy',
    'CA': 'Canada',
    'AU': 'Australia',
    'JP': 'Japan',
    'BR': 'Brazil',
    'MX': 'Mexico',
    'IN': 'India',
    'CN': 'China',
    'RU': 'Russia',
    'KR': 'South Korea',
  };
}
