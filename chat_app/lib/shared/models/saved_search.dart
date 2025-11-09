import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:hive/hive.dart';
import 'chat_message.dart';

part 'saved_search.freezed.dart';
part 'saved_search.g.dart';

@freezed
class SavedSearch with _$SavedSearch {
  const factory SavedSearch({
    required List<ChatMessage> messages,
    @JsonKey(name: 'session_id') required String sessionId,
    required String category,
    required int timestamp,
  }) = _SavedSearch;

  factory SavedSearch.fromJson(Map<String, dynamic> json) => _$SavedSearchFromJson(json);
}
