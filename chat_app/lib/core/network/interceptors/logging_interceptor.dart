import 'package:dio/dio.dart';
import 'package:logger/logger.dart';

/// Interceptor for logging HTTP requests and responses
class LoggingInterceptor extends Interceptor {
  final Logger _logger = Logger(
    printer: PrettyPrinter(
      methodCount: 0,
      errorMethodCount: 5,
      lineLength: 80,
      colors: true,
      printEmojis: true,
      dateTimeFormat: DateTimeFormat.onlyTimeAndSinceStart,
    ),
  );

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    _logger.d(
      '┌─── REQUEST ───────────────────────────────────────\n'
      '│ ${options.method} ${options.uri}\n'
      '│ Headers: ${options.headers}\n'
      '│ Query: ${options.queryParameters}\n'
      '${options.data != null ? '│ Body: ${options.data}\n' : ''}'
      '└───────────────────────────────────────────────────',
    );
    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    _logger.i(
      '┌─── RESPONSE ──────────────────────────────────────\n'
      '│ ${response.statusCode} ${response.requestOptions.method} ${response.requestOptions.uri}\n'
      '│ Headers: ${response.headers}\n'
      '│ Body: ${response.data}\n'
      '└───────────────────────────────────────────────────',
    );
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    _logger.e(
      '┌─── ERROR ─────────────────────────────────────────\n'
      '│ ${err.requestOptions.method} ${err.requestOptions.uri}\n'
      '│ Type: ${err.type}\n'
      '│ Message: ${err.message}\n'
      '${err.response != null ? '│ Status: ${err.response?.statusCode}\n' : ''}'
      '${err.response != null ? '│ Response: ${err.response?.data}\n' : ''}'
      '└───────────────────────────────────────────────────',
    );
    handler.next(err);
  }
}
