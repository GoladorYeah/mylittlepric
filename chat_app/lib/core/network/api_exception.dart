/// Custom exception for API errors
class ApiException implements Exception {
  final String message;
  final int? statusCode;
  final ApiExceptionType type;
  final dynamic data;

  const ApiException({
    required this.message,
    this.statusCode,
    this.type = ApiExceptionType.unknown,
    this.data,
  });

  @override
  String toString() {
    return 'ApiException: $message (statusCode: $statusCode, type: $type)';
  }

  /// Check if error is due to network issues
  bool get isNetworkError =>
      type == ApiExceptionType.network ||
      type == ApiExceptionType.timeout;

  /// Check if error requires authentication
  bool get isAuthError => type == ApiExceptionType.unauthorized;

  /// Check if error is server-side
  bool get isServerError => type == ApiExceptionType.server;
}

/// Types of API exceptions
enum ApiExceptionType {
  /// Network connectivity issues
  network,

  /// Request timeout
  timeout,

  /// Bad request (4xx except 401, 404)
  badRequest,

  /// Unauthorized (401)
  unauthorized,

  /// Not found (404)
  notFound,

  /// Server error (5xx)
  server,

  /// Request cancelled
  cancel,

  /// Unknown error
  unknown,
}
