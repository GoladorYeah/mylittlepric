import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../config/app_config.dart';
import '../storage/storage_service.dart';
import 'interceptors/auth_interceptor.dart';
import 'interceptors/logging_interceptor.dart';
import 'interceptors/retry_interceptor.dart';
import 'api_exception.dart';

/// HTTP client based on Dio with interceptors and error handling
class DioClient {
  late final Dio _dio;

  DioClient() {
    _dio = Dio(_baseOptions);
    _setupInterceptors();
  }

  BaseOptions get _baseOptions => BaseOptions(
        baseUrl: AppConfig.apiBaseUrl,
        connectTimeout: AppConfig.connectionTimeout,
        receiveTimeout: AppConfig.receiveTimeout,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      );

  void _setupInterceptors() {
    _dio.interceptors.addAll([
      AuthInterceptor(),
      RetryInterceptor(_dio),
      if (AppConfig.enableLogging) LoggingInterceptor(),
    ]);
  }

  /// GET request
  Future<Response<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      return await _dio.get<T>(
        path,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// POST request
  Future<Response<T>> post<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      return await _dio.post<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// PUT request
  Future<Response<T>> put<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      return await _dio.put<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// DELETE request
  Future<Response<T>> delete<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      return await _dio.delete<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// PATCH request
  Future<Response<T>> patch<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      return await _dio.patch<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// Handle Dio errors and convert to ApiException
  ApiException _handleError(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return ApiException(
          message: 'Connection timeout',
          statusCode: 408,
          type: ApiExceptionType.timeout,
        );

      case DioExceptionType.badResponse:
        final statusCode = error.response?.statusCode ?? 500;
        final message = _extractErrorMessage(error.response?.data);

        if (statusCode == 401) {
          return ApiException(
            message: message ?? 'Unauthorized',
            statusCode: statusCode,
            type: ApiExceptionType.unauthorized,
          );
        } else if (statusCode == 404) {
          return ApiException(
            message: message ?? 'Resource not found',
            statusCode: statusCode,
            type: ApiExceptionType.notFound,
          );
        } else if (statusCode >= 500) {
          return ApiException(
            message: message ?? 'Server error',
            statusCode: statusCode,
            type: ApiExceptionType.server,
          );
        } else {
          return ApiException(
            message: message ?? 'Request failed',
            statusCode: statusCode,
            type: ApiExceptionType.badRequest,
          );
        }

      case DioExceptionType.cancel:
        return ApiException(
          message: 'Request cancelled',
          type: ApiExceptionType.cancel,
        );

      case DioExceptionType.connectionError:
        return ApiException(
          message: 'No internet connection',
          type: ApiExceptionType.network,
        );

      case DioExceptionType.badCertificate:
        return ApiException(
          message: 'SSL certificate error',
          type: ApiExceptionType.network,
        );

      case DioExceptionType.unknown:
      default:
        return ApiException(
          message: error.message ?? 'Unknown error occurred',
          type: ApiExceptionType.unknown,
        );
    }
  }

  /// Extract error message from response data
  String? _extractErrorMessage(dynamic data) {
    if (data == null) return null;

    if (data is Map<String, dynamic>) {
      // Try common error message keys
      return data['message'] ??
          data['error'] ??
          data['detail'] ??
          data['msg'];
    }

    if (data is String) {
      return data;
    }

    return null;
  }

  /// Update authorization header
  void updateAuthToken(String token) {
    _dio.options.headers['Authorization'] = 'Bearer $token';
  }

  /// Remove authorization header
  void clearAuthToken() {
    _dio.options.headers.remove('Authorization');
  }

  /// Get underlying Dio instance (use with caution)
  Dio get dio => _dio;
}

/// Provider for DioClient
final dioClientProvider = Provider<DioClient>((ref) {
  return DioClient();
});
