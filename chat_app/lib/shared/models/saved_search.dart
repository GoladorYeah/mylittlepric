import 'package:chat_app/shared/models/chat_message.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'saved_search.freezed.dart';
part 'saved_search.g.dart';

@freezed
class SavedSearch with _$SavedSearch {
  const factory SavedSearch({
    required String id,
    required List<ChatMessage> messages,
    @JsonKey(name: 'session_id') required String sessionId,
    required String category,
    required int timestamp,
  }) = _SavedSearch;

  const SavedSearch._();

  factory SavedSearch.fromJson(Map<String, dynamic> json) => _$SavedSearchFromJson(json);

  /// Get the first user query from messages
  String get query {
    final userMessages = messages.where((m) => m.role == MessageRole.user);
    return userMessages.isNotEmpty ? userMessages.first.content : '';
  }

  /// Get formatted date
  String get formattedDate {
    final date = DateTime.fromMillisecondsSinceEpoch(timestamp);
    final now = DateTime.now();
    final difference = now.difference(date);

    if (difference.inDays == 0) {
      return 'Today ${date.hour}:${date.minute.toString().padLeft(2, '0')}';
    } else if (difference.inDays == 1) {
      return 'Yesterday ${date.hour}:${date.minute.toString().padLeft(2, '0')}';
    } else if (difference.inDays < 7) {
      return '${difference.inDays} days ago';
    } else {
      return '${date.day}/${date.month}/${date.year}';
    }
  }

  /// Get the number of products in the search
  int get productCount {
    return messages.where((m) => m.products?.isNotEmpty ?? false).fold(
      0,
      (sum, message) => sum + (message.products?.length ?? 0),
    );
  }
}
