import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../shared/models/saved_search.dart';

part 'session_state.freezed.dart';

/// Session history state
@freezed
class SessionState with _$SessionState {
  const factory SessionState({
    /// List of saved searches
    @Default([]) List<SavedSearch> searches,

    /// Whether searches are being loaded
    @Default(false) bool isLoading,

    /// Whether more searches can be loaded
    @Default(true) bool hasMore,

    /// Current page offset for pagination
    @Default(0) int offset,

    /// Number of items per page
    @Default(20) int limit,

    /// Error message
    String? error,

    /// Whether a search is being deleted
    @Default(false) bool isDeleting,
  }) = _SessionState;

  const SessionState._();

  /// Initial state
  factory SessionState.initial() => const SessionState();

  /// Check if there are any searches
  bool get hasSearches => searches.isNotEmpty;

  /// Get total number of searches
  int get totalSearches => searches.length;
}
