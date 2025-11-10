import 'package:chat_app/shared/models/product.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'chat_response.freezed.dart';
part 'chat_response.g.dart';

@freezed
class ChatResponse with _$ChatResponse {
  const factory ChatResponse({
    @JsonKey(name: 'session_id') required String sessionId,
    required String message,
    @JsonKey(name: 'quick_replies') List<String>? quickReplies,
    List<Product>? products,
    @JsonKey(name: 'response_type') String? responseType,
    @JsonKey(name: 'search_type') String? searchType,
  }) = _ChatResponse;

  factory ChatResponse.fromJson(Map<String, dynamic> json) => _$ChatResponseFromJson(json);
}
