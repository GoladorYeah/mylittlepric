import 'package:dio/dio.dart';
import '../../config/app_config.dart';
import '../../storage/storage_service.dart';

/// Interceptor for adding authentication tokens to requests
class AuthInterceptor extends Interceptor {
  AuthInterceptor();

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    // Get access token from storage
    final prefs = StorageService.prefs;
    final accessToken = prefs.getString(AppConfig.accessTokenKey);

    // Add authorization header if token exists
    if (accessToken != null && accessToken.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $accessToken';
    }

    // Get session ID from storage
    final sessionId = prefs.getString(AppConfig.sessionIdKey);

    // Add session ID header if exists
    if (sessionId != null && sessionId.isNotEmpty) {
      options.headers['X-Session-ID'] = sessionId;
    }

    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    // Handle 401 Unauthorized - token expired
    if (err.response?.statusCode == 401) {
      // Try to refresh token
      final prefs = StorageService.prefs;
      final refreshToken = prefs.getString(AppConfig.refreshTokenKey);

      if (refreshToken != null && refreshToken.isNotEmpty) {
        try {
          // TODO: Implement token refresh logic
          // For now, just clear tokens
          await prefs.remove(AppConfig.accessTokenKey);
          await prefs.remove(AppConfig.refreshTokenKey);
        } catch (e) {
          // Refresh failed, clear tokens
          await prefs.remove(AppConfig.accessTokenKey);
          await prefs.remove(AppConfig.refreshTokenKey);
        }
      }
    }

    handler.next(err);
  }
}
