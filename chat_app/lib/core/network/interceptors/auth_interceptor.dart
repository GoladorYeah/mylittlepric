import 'package:dio/dio.dart';
import '../../config/app_config.dart';
import '../../storage/storage_service.dart';

/// Interceptor for adding authentication tokens to requests
class AuthInterceptor extends Interceptor {
  final StorageService _storageService;

  AuthInterceptor(this._storageService);

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    // Get access token from storage
    final accessToken = await _storageService.getString(AppConfig.accessTokenKey);

    // Add authorization header if token exists
    if (accessToken != null && accessToken.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $accessToken';
    }

    // Get session ID from storage
    final sessionId = await _storageService.getString(AppConfig.sessionIdKey);

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
      final refreshToken = await _storageService.getString(AppConfig.refreshTokenKey);

      if (refreshToken != null && refreshToken.isNotEmpty) {
        try {
          // TODO: Implement token refresh logic
          // For now, just clear tokens
          await _storageService.remove(AppConfig.accessTokenKey);
          await _storageService.remove(AppConfig.refreshTokenKey);
        } catch (e) {
          // Refresh failed, clear tokens
          await _storageService.remove(AppConfig.accessTokenKey);
          await _storageService.remove(AppConfig.refreshTokenKey);
        }
      }
    }

    handler.next(err);
  }
}
