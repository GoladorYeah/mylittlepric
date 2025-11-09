import 'package:freezed_annotation/freezed_annotation.dart';
import 'chat_message.dart';

part 'session.freezed.dart';
part 'session.g.dart';

@freezed
class SearchState with _$SearchState {
  const factory SearchState({
    String? category,
    String? status,
    @JsonKey(name: 'last_product') LastProduct? lastProduct,
  }) = _SearchState;

  factory SearchState.fromJson(Map<String, dynamic> json) => _$SearchStateFromJson(json);
}

@freezed
class LastProduct with _$LastProduct {
  const factory LastProduct({
    required String name,
    required String price,
  }) = _LastProduct;

  factory LastProduct.fromJson(Map<String, dynamic> json) => _$LastProductFromJson(json);
}

@freezed
class SessionResponse with _$SessionResponse {
  const factory SessionResponse({
    @JsonKey(name: 'session_id') required String sessionId,
    required List<SessionMessage> messages,
    @JsonKey(name: 'search_state') SearchState? searchState,
  }) = _SessionResponse;

  factory SessionResponse.fromJson(Map<String, dynamic> json) => _$SessionResponseFromJson(json);
}
