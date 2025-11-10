import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../shared/models/session.dart';
import '../../../shared/models/saved_search.dart';
import '../../config/app_config.dart';
import '../dio_client.dart';

/// API service for session operations
class SessionApiService {
  final DioClient _client;

  SessionApiService(this._client);

  /// Get or create a session
  Future<SessionResponse> getOrCreateSession({
    String? sessionId,
  }) async {
    final response = await _client.get(
      AppConfig.sessionEndpoint,
      queryParameters: {
        if (sessionId != null) 'session_id': sessionId,
      },
    );

    return SessionResponse.fromJson(response.data as Map<String, dynamic>);
  }

  /// Get session by ID
  Future<SessionResponse> getSession(String sessionId) async {
    final response = await _client.get(
      '${AppConfig.sessionEndpoint}/$sessionId',
    );

    return SessionResponse.fromJson(response.data as Map<String, dynamic>);
  }

  /// Delete session
  Future<void> deleteSession(String sessionId) async {
    await _client.delete(
      '${AppConfig.sessionEndpoint}/$sessionId',
    );
  }

  /// Clear session history
  Future<void> clearSession(String sessionId) async {
    await _client.post(
      '${AppConfig.sessionEndpoint}/$sessionId/clear',
    );
  }

  /// Get search history
  Future<List<SavedSearch>> getSearchHistory({
    int? limit,
    int? offset,
  }) async {
    final response = await _client.get(
      AppConfig.searchHistoryEndpoint,
      queryParameters: {
        if (limit != null) 'limit': limit,
        if (offset != null) 'offset': offset,
      },
    );

    final data = response.data as List<dynamic>;
    return data
        .map((json) => SavedSearch.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  /// Save search to history
  Future<SavedSearch> saveSearch({
    required String query,
    required String category,
    List<String>? productIds,
  }) async {
    final response = await _client.post(
      AppConfig.searchHistoryEndpoint,
      data: {
        'query': query,
        'category': category,
        if (productIds != null) 'product_ids': productIds,
      },
    );

    return SavedSearch.fromJson(response.data as Map<String, dynamic>);
  }

  /// Delete search from history
  Future<void> deleteSearch(String searchId) async {
    await _client.delete(
      '${AppConfig.searchHistoryEndpoint}/$searchId',
    );
  }
}

/// Provider for SessionApiService
final sessionApiServiceProvider = Provider<SessionApiService>((ref) {
  final client = ref.watch(dioClientProvider);
  return SessionApiService(client);
});
