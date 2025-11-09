import 'package:freezed_annotation/freezed_annotation.dart';
import 'product.dart';

part 'chat_message.freezed.dart';
part 'chat_message.g.dart';

enum MessageRole {
  @JsonValue('user')
  user,
  @JsonValue('assistant')
  assistant,
}

@freezed
class ChatMessage with _$ChatMessage {
  const factory ChatMessage({
    required String id,
    required MessageRole role,
    required String content,
    required int timestamp,
    @JsonKey(name: 'quick_replies') List<String>? quickReplies,
    List<Product>? products,
    @JsonKey(name: 'search_type') String? searchType,
    @Default(false) bool isLocal, // true if sent from this device
  }) = _ChatMessage;

  factory ChatMessage.fromJson(Map<String, dynamic> json) => _$ChatMessageFromJson(json);
}

// Session message format (from backend)
@freezed
class SessionMessage with _$SessionMessage {
  const factory SessionMessage({
    required String role,
    required String content,
    String? timestamp,
    @JsonKey(name: 'quick_replies') List<String>? quickReplies,
    List<Product>? products,
    @JsonKey(name: 'search_type') String? searchType,
  }) = _SessionMessage;

  factory SessionMessage.fromJson(Map<String, dynamic> json) => _$SessionMessageFromJson(json);
}
