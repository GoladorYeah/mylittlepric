import 'package:chat_app/core/config/app_config.dart';
import 'package:chat_app/shared/utils/logger.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:shared_preferences/shared_preferences.dart';

class StorageService {
  static late SharedPreferences _prefs;
  static bool _initialized = false;

  /// Initialize all storage services (Hive + SharedPreferences)
  static Future<void> init() async {
    if (_initialized) return;

    try {
      // Initialize Hive
      await Hive.initFlutter();

      // Open required boxes
      await Future.wait([
        Hive.openBox(AppConfig.messagesBox),
        Hive.openBox(AppConfig.sessionBox),
        Hive.openBox(AppConfig.preferencesBox),
        Hive.openBox(AppConfig.cacheBox),
      ]);

      // Initialize SharedPreferences
      _prefs = await SharedPreferences.getInstance();

      _initialized = true;
      AppLogger.info('✅ Storage services initialized');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to initialize storage', e, stackTrace);
      rethrow;
    }
  }

  /// Get SharedPreferences instance
  static SharedPreferences get prefs {
    if (!_initialized) {
      throw Exception('StorageService not initialized. Call StorageService.init() first.');
    }
    return _prefs;
  }

  /// Get a Hive box by name
  static Box getBox(String boxName) {
    if (!_initialized) {
      throw Exception('StorageService not initialized. Call StorageService.init() first.');
    }
    return Hive.box(boxName);
  }

  /// Clear all storage (useful for logout)
  static Future<void> clearAll() async {
    try {
      // Clear all Hive boxes
      await Future.wait([
        Hive.box(AppConfig.messagesBox).clear(),
        Hive.box(AppConfig.sessionBox).clear(),
        Hive.box(AppConfig.preferencesBox).clear(),
        Hive.box(AppConfig.cacheBox).clear(),
      ]);

      // Clear SharedPreferences
      await _prefs.clear();

      AppLogger.info('✅ All storage cleared');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to clear storage', e, stackTrace);
      rethrow;
    }
  }

  /// Clear expired cache entries
  static Future<void> clearExpiredCache() async {
    try {
      final cacheBox = Hive.box(AppConfig.cacheBox);
      final now = DateTime.now().millisecondsSinceEpoch;
      final keysToDelete = <String>[];

      for (final key in cacheBox.keys) {
        final entry = cacheBox.get(key);
        if (entry is Map && entry['expiresAt'] != null) {
          final expiresAt = entry['expiresAt'] as int;
          if (expiresAt < now) {
            keysToDelete.add(key.toString());
          }
        }
      }

      await Future.wait(
        keysToDelete.map((key) => cacheBox.delete(key)),
      );

      if (keysToDelete.isNotEmpty) {
        AppLogger.info('🗑️  Cleared ${keysToDelete.length} expired cache entries');
      }
    } catch (e, stackTrace) {
      AppLogger.error('Failed to clear expired cache', e, stackTrace);
    }
  }

  /// Dispose storage services
  static Future<void> dispose() async {
    await Hive.close();
    _initialized = false;
    AppLogger.info('Storage services disposed');
  }
}
