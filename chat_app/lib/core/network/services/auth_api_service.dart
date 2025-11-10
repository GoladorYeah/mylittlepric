import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../shared/models/user.dart';
import '../../config/app_config.dart';
import '../dio_client.dart';

/// API service for authentication operations
class AuthApiService {
  final DioClient _client;

  AuthApiService(this._client);

  /// Login with OAuth provider
  Future<AuthResponse> login({
    required String provider,
    required String token,
  }) async {
    final response = await _client.post(
      AppConfig.loginEndpoint,
      data: {
        'provider': provider,
        'token': token,
      },
    );

    return AuthResponse.fromJson(response.data as Map<String, dynamic>);
  }

  /// Logout
  Future<void> logout() async {
    await _client.post(AppConfig.logoutEndpoint);
  }

  /// Refresh access token
  Future<TokenResponse> refreshToken(String refreshToken) async {
    final response = await _client.post(
      '/api/auth/refresh',
      data: {
        'refresh_token': refreshToken,
      },
    );

    return TokenResponse.fromJson(response.data as Map<String, dynamic>);
  }

  /// Get current user
  Future<User> getCurrentUser() async {
    final response = await _client.get('/api/auth/me');
    return User.fromJson(response.data as Map<String, dynamic>);
  }

  /// Update user preferences
  Future<User> updateUserPreferences({
    String? country,
    String? language,
    String? currency,
  }) async {
    final response = await _client.put(
      AppConfig.userPreferencesEndpoint,
      data: {
        if (country != null) 'country': country,
        if (language != null) 'language': language,
        if (currency != null) 'currency': currency,
      },
    );

    return User.fromJson(response.data as Map<String, dynamic>);
  }

  /// Delete user account
  Future<void> deleteAccount() async {
    await _client.delete('/api/auth/account');
  }
}

/// Authentication response
class AuthResponse {
  final String accessToken;
  final String refreshToken;
  final User user;

  AuthResponse({
    required this.accessToken,
    required this.refreshToken,
    required this.user,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
      user: User.fromJson(json['user'] as Map<String, dynamic>),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'access_token': accessToken,
      'refresh_token': refreshToken,
      'user': user.toJson(),
    };
  }
}

/// Token refresh response
class TokenResponse {
  final String accessToken;
  final String refreshToken;

  TokenResponse({
    required this.accessToken,
    required this.refreshToken,
  });

  factory TokenResponse.fromJson(Map<String, dynamic> json) {
    return TokenResponse(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'access_token': accessToken,
      'refresh_token': refreshToken,
    };
  }
}

/// Provider for AuthApiService
final authApiServiceProvider = Provider<AuthApiService>((ref) {
  final client = ref.watch(dioClientProvider);
  return AuthApiService(client);
});
