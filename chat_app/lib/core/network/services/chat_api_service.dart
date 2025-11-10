import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../shared/models/chat_response.dart';
import '../../../shared/models/product_details.dart';
import '../../config/app_config.dart';
import '../dio_client.dart';

/// API service for chat operations
class ChatApiService {
  final DioClient _client;

  ChatApiService(this._client);

  /// Send a chat message
  Future<ChatResponse> sendMessage({
    required String message,
    String? sessionId,
    String? country,
    String? language,
    String? currency,
  }) async {
    final response = await _client.post(
      AppConfig.chatEndpoint,
      data: {
        'message': message,
        if (sessionId != null) 'session_id': sessionId,
        if (country != null) 'country': country,
        if (language != null) 'language': language,
        if (currency != null) 'currency': currency,
      },
    );

    return ChatResponse.fromJson(response.data as Map<String, dynamic>);
  }

  /// Get product details by page token
  Future<ProductDetailsResponse> getProductDetails({
    required String pageToken,
    String? country,
    String? language,
  }) async {
    final response = await _client.post(
      AppConfig.productDetailsEndpoint,
      data: {
        'page_token': pageToken,
        if (country != null) 'country': country,
        if (language != null) 'language': language,
      },
    );

    return ProductDetailsResponse.fromJson(response.data as Map<String, dynamic>);
  }

  /// Send quick reply
  Future<ChatResponse> sendQuickReply({
    required String reply,
    required String sessionId,
  }) async {
    return sendMessage(
      message: reply,
      sessionId: sessionId,
    );
  }

  /// Continue search with parameters
  Future<ChatResponse> continueSearch({
    required String parameters,
    required String sessionId,
  }) async {
    return sendMessage(
      message: parameters,
      sessionId: sessionId,
    );
  }
}

/// Provider for ChatApiService
final chatApiServiceProvider = Provider<ChatApiService>((ref) {
  final client = ref.watch(dioClientProvider);
  return ChatApiService(client);
});
