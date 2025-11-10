import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/storage/storage_service.dart';
import '../../../core/config/app_config.dart';
import '../../../shared/utils/logger.dart';
import '../../auth/providers/auth_provider.dart';
import 'settings_state.dart';

/// Settings provider - manages app settings
final settingsProvider = StateNotifierProvider<SettingsNotifier, SettingsState>((ref) {
  final authNotifier = ref.read(authProvider.notifier);
  return SettingsNotifier(authNotifier);
});

/// Settings notifier - handles settings logic
class SettingsNotifier extends StateNotifier<SettingsState> {
  final AuthNotifier _authNotifier;

  SettingsNotifier(this._authNotifier) : super(SettingsState.initial()) {
    _loadSettings();
  }

  /// Load settings from storage
  Future<void> _loadSettings() async {
    try {
      final prefs = StorageService.prefs;

      // Load theme mode
      final themeModeString = prefs.getString(_themeModeKey);
      final themeMode = themeModeString != null
          ? ThemeMode.values.firstWhere(
              (mode) => mode.toString() == themeModeString,
              orElse: () => ThemeMode.system,
            )
          : ThemeMode.system;

      // Load other settings
      final country = prefs.getString(_countryKey) ?? 'US';
      final language = prefs.getString(_languageKey) ?? 'en';
      final currency = prefs.getString(_currencyKey) ?? 'USD';
      final notificationsEnabled = prefs.getBool(_notificationsKey) ?? true;
      final soundEnabled = prefs.getBool(_soundKey) ?? true;

      state = SettingsState(
        themeMode: themeMode,
        country: country,
        language: language,
        currency: currency,
        notificationsEnabled: notificationsEnabled,
        soundEnabled: soundEnabled,
      );

      AppLogger.info('✅ Settings loaded: theme=$themeMode, country=$country, language=$language');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to load settings', e, stackTrace);
    }
  }

  /// Save settings to storage
  Future<void> _saveSettings() async {
    try {
      final prefs = StorageService.prefs;

      await prefs.setString(_themeModeKey, state.themeMode.toString());
      await prefs.setString(_countryKey, state.country);
      await prefs.setString(_languageKey, state.language);
      await prefs.setString(_currencyKey, state.currency);
      await prefs.setBool(_notificationsKey, state.notificationsEnabled);
      await prefs.setBool(_soundKey, state.soundEnabled);

      AppLogger.info('✅ Settings saved');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to save settings', e, stackTrace);
      state = state.copyWith(error: e.toString());
    }
  }

  /// Set theme mode
  Future<void> setThemeMode(ThemeMode mode) async {
    state = state.copyWith(themeMode: mode);
    await _saveSettings();
  }

  /// Toggle between light and dark mode
  Future<void> toggleTheme() async {
    final newMode = state.themeMode == ThemeMode.dark
        ? ThemeMode.light
        : ThemeMode.dark;
    await setThemeMode(newMode);
  }

  /// Set country
  Future<void> setCountry(String countryCode) async {
    state = state.copyWith(country: countryCode, isSaving: true);
    await _saveSettings();

    // Update user preferences on backend if authenticated
    try {
      await _authNotifier.updatePreferences(country: countryCode);
    } catch (e, stackTrace) {
      AppLogger.error('Failed to update country on backend', e, stackTrace);
    } finally {
      state = state.copyWith(isSaving: false);
    }
  }

  /// Set language
  Future<void> setLanguage(String languageCode) async {
    state = state.copyWith(language: languageCode, isSaving: true);
    await _saveSettings();

    // Update user preferences on backend if authenticated
    try {
      await _authNotifier.updatePreferences(language: languageCode);
    } catch (e, stackTrace) {
      AppLogger.error('Failed to update language on backend', e, stackTrace);
    } finally {
      state = state.copyWith(isSaving: false);
    }
  }

  /// Set currency
  Future<void> setCurrency(String currencyCode) async {
    state = state.copyWith(currency: currencyCode, isSaving: true);
    await _saveSettings();

    // Update user preferences on backend if authenticated
    try {
      await _authNotifier.updatePreferences(currency: currencyCode);
    } catch (e, stackTrace) {
      AppLogger.error('Failed to update currency on backend', e, stackTrace);
    } finally {
      state = state.copyWith(isSaving: false);
    }
  }

  /// Set all preferences at once
  Future<void> setPreferences({
    String? country,
    String? language,
    String? currency,
  }) async {
    state = state.copyWith(
      country: country ?? state.country,
      language: language ?? state.language,
      currency: currency ?? state.currency,
      isSaving: true,
    );

    await _saveSettings();

    // Update user preferences on backend if authenticated
    try {
      await _authNotifier.updatePreferences(
        country: country,
        language: language,
        currency: currency,
      );
    } catch (e, stackTrace) {
      AppLogger.error('Failed to update preferences on backend', e, stackTrace);
    } finally {
      state = state.copyWith(isSaving: false);
    }
  }

  /// Toggle notifications
  Future<void> toggleNotifications() async {
    state = state.copyWith(notificationsEnabled: !state.notificationsEnabled);
    await _saveSettings();
  }

  /// Toggle sound
  Future<void> toggleSound() async {
    state = state.copyWith(soundEnabled: !state.soundEnabled);
    await _saveSettings();
  }

  /// Clear error
  void clearError() {
    state = state.copyWith(error: null);
  }

  /// Reset settings to defaults
  Future<void> resetToDefaults() async {
    state = SettingsState.initial();
    await _saveSettings();
    AppLogger.info('✅ Settings reset to defaults');
  }

  // Storage keys
  static const _themeModeKey = '${AppConfig.preferencesBox}_theme_mode';
  static const _countryKey = '${AppConfig.preferencesBox}_country';
  static const _languageKey = '${AppConfig.preferencesBox}_language';
  static const _currencyKey = '${AppConfig.preferencesBox}_currency';
  static const _notificationsKey = '${AppConfig.preferencesBox}_notifications';
  static const _soundKey = '${AppConfig.preferencesBox}_sound';
}

/// Helper providers
final themeModeProvider = Provider<ThemeMode>((ref) {
  return ref.watch(settingsProvider).themeMode;
});

final countryProvider = Provider<String>((ref) {
  return ref.watch(settingsProvider).country;
});

final languageProvider = Provider<String>((ref) {
  return ref.watch(settingsProvider).language;
});

final currencyProvider = Provider<String>((ref) {
  return ref.watch(settingsProvider).currency;
});

final notificationsEnabledProvider = Provider<bool>((ref) {
  return ref.watch(settingsProvider).notificationsEnabled;
});

final soundEnabledProvider = Provider<bool>((ref) {
  return ref.watch(settingsProvider).soundEnabled;
});
