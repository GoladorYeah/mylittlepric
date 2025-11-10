import 'package:dio/dio.dart';
import '../../config/app_config.dart';

/// Interceptor for retrying failed requests
class RetryInterceptor extends Interceptor {
  final Dio _dio;

  RetryInterceptor(this._dio);

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    // Only retry on network errors and timeouts
    if (!_shouldRetry(err)) {
      return handler.next(err);
    }

    // Get retry count from request options
    final retryCount = err.requestOptions.extra['retry_count'] as int? ?? 0;

    // Check if max retries reached
    if (retryCount >= AppConfig.maxRetries) {
      return handler.next(err);
    }

    // Calculate delay with exponential backoff
    final delay = _calculateDelay(retryCount);

    // Wait before retrying
    await Future.delayed(delay);

    // Increment retry count
    err.requestOptions.extra['retry_count'] = retryCount + 1;

    // Retry the request
    try {
      final response = await _dio.fetch(err.requestOptions);
      return handler.resolve(response);
    } on DioException catch (e) {
      return handler.next(e);
    }
  }

  /// Check if request should be retried
  bool _shouldRetry(DioException err) {
    // Retry on connection errors and timeouts
    return err.type == DioExceptionType.connectionTimeout ||
        err.type == DioExceptionType.sendTimeout ||
        err.type == DioExceptionType.receiveTimeout ||
        err.type == DioExceptionType.connectionError ||
        // Retry on server errors (5xx)
        (err.response?.statusCode ?? 0) >= 500;
  }

  /// Calculate delay with exponential backoff
  Duration _calculateDelay(int retryCount) {
    // Exponential backoff: 1s, 2s, 4s
    final seconds = (retryCount + 1) * 1;
    return Duration(seconds: seconds);
  }
}
