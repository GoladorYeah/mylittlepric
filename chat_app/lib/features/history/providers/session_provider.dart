import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:riverpod/riverpod.dart';
import '../../../core/network/services/session_api_service.dart';
import '../../../shared/utils/logger.dart';
import '../../../shared/models/saved_search.dart';
import 'session_state.dart';

/// Session provider - manages search history
final sessionProvider = StateNotifierProvider<SessionNotifier, SessionState>((ref) {
  final sessionApiService = ref.watch(sessionApiServiceProvider);
  return SessionNotifier(sessionApiService);
});

/// Session notifier - handles session history logic
class SessionNotifier extends StateNotifier<SessionState> {
  final SessionApiService _sessionApiService;

  SessionNotifier(this._sessionApiService) : super(SessionState.initial()) {
    // Load initial searches
    loadSearches();
  }

  /// Load searches from server
  Future<void> loadSearches({bool refresh = false}) async {
    // If refreshing, reset offset
    if (refresh) {
      state = SessionState.initial();
    }

    // Don't load if already loading or no more items
    if (state.isLoading || (!refresh && !state.hasMore)) {
      return;
    }

    try {
      state = state.copyWith(isLoading: true, error: null);

      final newSearches = await _sessionApiService.getSearchHistory(
        limit: state.limit,
        offset: state.offset,
      );

      // Update state
      final updatedSearches = refresh
          ? newSearches
          : [...state.searches, ...newSearches];

      state = state.copyWith(
        searches: updatedSearches,
        offset: state.offset + newSearches.length,
        hasMore: newSearches.length >= state.limit,
        isLoading: false,
      );

      AppLogger.info('✅ Loaded ${newSearches.length} searches (total: ${updatedSearches.length})');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to load searches', e, stackTrace);
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Load more searches (pagination)
  Future<void> loadMoreSearches() async {
    await loadSearches();
  }

  /// Refresh searches
  Future<void> refreshSearches() async {
    await loadSearches(refresh: true);
  }

  /// Save a new search
  Future<void> saveSearch({
    required String query,
    required String category,
    List<String>? productIds,
  }) async {
    try {
      final savedSearch = await _sessionApiService.saveSearch(
        query: query,
        category: category,
        productIds: productIds,
      );

      // Add to the beginning of the list
      state = state.copyWith(
        searches: [savedSearch, ...state.searches],
      );

      AppLogger.info('✅ Search saved: $query');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to save search', e, stackTrace);
      state = state.copyWith(error: e.toString());
    }
  }

  /// Delete a search
  Future<void> deleteSearch(String searchId) async {
    try {
      state = state.copyWith(isDeleting: true);

      await _sessionApiService.deleteSearch(searchId);

      // Remove from list
      state = state.copyWith(
        searches: state.searches.where((s) => s.id != searchId).toList(),
        isDeleting: false,
      );

      AppLogger.info('✅ Search deleted: $searchId');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to delete search', e, stackTrace);
      state = state.copyWith(
        isDeleting: false,
        error: e.toString(),
      );
    }
  }

  /// Clear all searches
  Future<void> clearAllSearches() async {
    try {
      // Delete all searches one by one
      for (final search in state.searches) {
        await _sessionApiService.deleteSearch(search.id);
      }

      state = SessionState.initial();
      AppLogger.info('✅ All searches cleared');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to clear searches', e, stackTrace);
      state = state.copyWith(error: e.toString());
    }
  }

  /// Clear error
  void clearError() {
    state = state.copyWith(error: null);
  }

  /// Search in history by query
  List<SavedSearch> searchInHistory(String query) {
    if (query.isEmpty) {
      return state.searches;
    }

    final lowerQuery = query.toLowerCase();
    return state.searches.where((search) {
      return search.query.toLowerCase().contains(lowerQuery) ||
          search.category.toLowerCase().contains(lowerQuery);
    }).toList();
  }

  /// Get searches by category
  List<SavedSearch> getSearchesByCategory(String category) {
    return state.searches.where((search) {
      return search.category.toLowerCase() == category.toLowerCase();
    }).toList();
  }

  /// Get unique categories from searches
  List<String> getUniqueCategories() {
    final categories = state.searches.map((s) => s.category).toSet();
    return categories.toList()..sort();
  }
}

/// Helper providers
final searchHistoryProvider = Provider<List<SavedSearch>>((ref) {
  return ref.watch(sessionProvider).searches;
});

final searchHistoryLoadingProvider = Provider<bool>((ref) {
  return ref.watch(sessionProvider).isLoading;
});

final searchHistoryHasMoreProvider = Provider<bool>((ref) {
  return ref.watch(sessionProvider).hasMore;
});

final uniqueCategoriesProvider = Provider<List<String>>((ref) {
  final notifier = ref.read(sessionProvider.notifier);
  return notifier.getUniqueCategories();
});
