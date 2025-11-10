import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/services/auth_api_service.dart';
import '../../../core/storage/storage_service.dart';
import '../../../core/config/app_config.dart';
import '../../../shared/utils/logger.dart';
import '../../../shared/models/user.dart';
import 'auth_state.dart';

/// Auth provider - manages authentication state
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final authApiService = ref.watch(authApiServiceProvider);
  return AuthNotifier(authApiService);
});

/// Auth notifier - handles authentication logic
class AuthNotifier extends StateNotifier<AuthState> {
  final AuthApiService _authApiService;

  AuthNotifier(this._authApiService) : super(AuthState.initial()) {
    // Try to restore session on initialization
    _restoreSession();
  }

  /// Restore session from storage
  Future<void> _restoreSession() async {
    try {
      final prefs = StorageService.prefs;
      final accessToken = prefs.getString(AppConfig.accessTokenKey);
      final refreshToken = prefs.getString(AppConfig.refreshTokenKey);

      if (accessToken != null && refreshToken != null) {
        state = AuthState.loading();

        // Try to get current user
        try {
          final user = await _authApiService.getCurrentUser();
          state = AuthState.authenticated(
            user: user,
            accessToken: accessToken,
            refreshToken: refreshToken,
          );
          AppLogger.info('✅ Session restored for user: ${user.email}');
        } catch (e) {
          // If getting user fails, try to refresh token
          await _tryRefreshToken();
        }
      }
    } catch (e, stackTrace) {
      AppLogger.error('Failed to restore session', e, stackTrace);
      state = AuthState.initial();
    }
  }

  /// Try to refresh access token
  Future<void> _tryRefreshToken() async {
    try {
      final prefs = StorageService.prefs;
      final refreshToken = prefs.getString(AppConfig.refreshTokenKey);

      if (refreshToken == null) {
        await logout();
        return;
      }

      final tokenResponse = await _authApiService.refreshToken(refreshToken);

      // Save new tokens
      await prefs.setString(AppConfig.accessTokenKey, tokenResponse.accessToken);
      await prefs.setString(AppConfig.refreshTokenKey, tokenResponse.refreshToken);

      // Get user with new token
      final user = await _authApiService.getCurrentUser();
      state = AuthState.authenticated(
        user: user,
        accessToken: tokenResponse.accessToken,
        refreshToken: tokenResponse.refreshToken,
      );
      AppLogger.info('✅ Token refreshed successfully');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to refresh token', e, stackTrace);
      await logout();
    }
  }

  /// Login with OAuth provider
  Future<void> login({
    required String provider,
    required String token,
  }) async {
    try {
      state = AuthState.loading();

      final response = await _authApiService.login(
        provider: provider,
        token: token,
      );

      // Save tokens to storage
      final prefs = StorageService.prefs;
      await prefs.setString(AppConfig.accessTokenKey, response.accessToken);
      await prefs.setString(AppConfig.refreshTokenKey, response.refreshToken);

      state = AuthState.authenticated(
        user: response.user,
        accessToken: response.accessToken,
        refreshToken: response.refreshToken,
      );

      AppLogger.info('✅ Login successful: ${response.user.email}');
    } catch (e, stackTrace) {
      AppLogger.error('Login failed', e, stackTrace);
      state = AuthState.error(e.toString());
    }
  }

  /// Logout
  Future<void> logout() async {
    try {
      // Call logout API
      await _authApiService.logout();
    } catch (e, stackTrace) {
      AppLogger.error('Logout API call failed', e, stackTrace);
    } finally {
      // Clear storage regardless of API call result
      final prefs = StorageService.prefs;
      await prefs.remove(AppConfig.accessTokenKey);
      await prefs.remove(AppConfig.refreshTokenKey);
      await StorageService.clearAll();

      state = AuthState.initial();
      AppLogger.info('✅ Logged out successfully');
    }
  }

  /// Update user preferences
  Future<void> updatePreferences({
    String? country,
    String? language,
    String? currency,
  }) async {
    try {
      final updatedUser = await _authApiService.updateUserPreferences(
        country: country,
        language: language,
        currency: currency,
      );

      state = state.copyWith(user: updatedUser);
      AppLogger.info('✅ User preferences updated');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to update preferences', e, stackTrace);
      state = state.copyWith(error: e.toString());
    }
  }

  /// Delete account
  Future<void> deleteAccount() async {
    try {
      state = AuthState.loading();
      await _authApiService.deleteAccount();
      await logout();
      AppLogger.info('✅ Account deleted successfully');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to delete account', e, stackTrace);
      state = state.copyWith(error: e.toString());
    }
  }

  /// Clear error
  void clearError() {
    state = state.copyWith(error: null);
  }
}

/// Helper providers for common auth checks
final isAuthenticatedProvider = Provider<bool>((ref) {
  return ref.watch(authProvider).isAuthenticated;
});

final currentUserProvider = Provider<User?>((ref) {
  return ref.watch(authProvider).user;
});

final authLoadingProvider = Provider<bool>((ref) {
  return ref.watch(authProvider).isLoading;
});
