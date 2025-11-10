import 'package:flutter/material.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'settings_state.freezed.dart';

/// Settings state
@freezed
class SettingsState with _$SettingsState {
  const factory SettingsState({
    /// Theme mode (light/dark/system)
    @Default(ThemeMode.system) ThemeMode themeMode,

    /// Selected country code (e.g., 'US', 'UA')
    @Default('US') String country,

    /// Selected language code (e.g., 'en', 'ru', 'uk')
    @Default('en') String language,

    /// Selected currency code (e.g., 'USD', 'UAH')
    @Default('USD') String currency,

    /// Whether to enable notifications
    @Default(true) bool notificationsEnabled,

    /// Whether to enable sound
    @Default(true) bool soundEnabled,

    /// Whether settings are being saved
    @Default(false) bool isSaving,

    /// Error message
    String? error,
  }) = _SettingsState;

  const SettingsState._();

  /// Initial state with defaults
  factory SettingsState.initial() => const SettingsState();

  /// Check if dark mode is enabled
  bool isDarkMode(BuildContext context) {
    if (themeMode == ThemeMode.system) {
      return MediaQuery.of(context).platformBrightness == Brightness.dark;
    }
    return themeMode == ThemeMode.dark;
  }
}

/// Supported countries
class Country {
  final String code;
  final String name;
  final String flag;

  const Country({
    required this.code,
    required this.name,
    required this.flag,
  });
}

/// Supported languages
class Language {
  final String code;
  final String name;
  final String nativeName;

  const Language({
    required this.code,
    required this.name,
    required this.nativeName,
  });
}

/// Supported currencies
class Currency {
  final String code;
  final String name;
  final String symbol;

  const Currency({
    required this.code,
    required this.name,
    required this.symbol,
  });
}

/// Available countries
const availableCountries = [
  Country(code: 'US', name: 'United States', flag: '🇺🇸'),
  Country(code: 'UA', name: 'Ukraine', flag: '🇺🇦'),
  Country(code: 'GB', name: 'United Kingdom', flag: '🇬🇧'),
  Country(code: 'DE', name: 'Germany', flag: '🇩🇪'),
  Country(code: 'FR', name: 'France', flag: '🇫🇷'),
  Country(code: 'ES', name: 'Spain', flag: '🇪🇸'),
  Country(code: 'IT', name: 'Italy', flag: '🇮🇹'),
  Country(code: 'PL', name: 'Poland', flag: '🇵🇱'),
  Country(code: 'CA', name: 'Canada', flag: '🇨🇦'),
  Country(code: 'AU', name: 'Australia', flag: '🇦🇺'),
];

/// Available languages
const availableLanguages = [
  Language(code: 'en', name: 'English', nativeName: 'English'),
  Language(code: 'uk', name: 'Ukrainian', nativeName: 'Українська'),
  Language(code: 'ru', name: 'Russian', nativeName: 'Русский'),
  Language(code: 'de', name: 'German', nativeName: 'Deutsch'),
  Language(code: 'fr', name: 'French', nativeName: 'Français'),
  Language(code: 'es', name: 'Spanish', nativeName: 'Español'),
  Language(code: 'it', name: 'Italian', nativeName: 'Italiano'),
  Language(code: 'pl', name: 'Polish', nativeName: 'Polski'),
];

/// Available currencies
const availableCurrencies = [
  Currency(code: 'USD', name: 'US Dollar', symbol: '\$'),
  Currency(code: 'UAH', name: 'Ukrainian Hryvnia', symbol: '₴'),
  Currency(code: 'EUR', name: 'Euro', symbol: '€'),
  Currency(code: 'GBP', name: 'British Pound', symbol: '£'),
  Currency(code: 'CAD', name: 'Canadian Dollar', symbol: 'C\$'),
  Currency(code: 'AUD', name: 'Australian Dollar', symbol: 'A\$'),
];
